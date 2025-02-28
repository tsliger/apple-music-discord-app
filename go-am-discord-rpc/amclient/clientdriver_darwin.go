package amclient

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var cancelPrevious context.CancelFunc
var rate time.Duration = songPollingRate * time.Millisecond
var DataCtx context.Context
var DataCancel context.CancelFunc

func Poll() {
	playingState, event, err := eventDetector()

	if err != nil {
		// Clear out on error?
	}

	// Call event handler
	eventHandler(event, playingState)

	time.Sleep(rate)
}

func eventHandler(event musicEvent, state playerState) {
	// Set activity based off event changes
	if event.noTrackPlaying {
		state.Url = "https://static-00.iconduck.com/assets.00/apple-music-icon-1024x1024-zncv5jwr.png"
	}

	//|| currentDiscordState == state
	if event.stateChanged || event.playheadChanged || event.songChanged {
		setDiscordActivity(state)
	}
}

func getAlbumArtUrl(state playerState) (string, error) {
	cachedUrl, err := getUrlFromCache(state.Artist, state.Album)

	// Cache hit
	if err == nil {
		return cachedUrl, nil
	}

	// Scrape and insert into cache
	state.Url, err = scrapeAlbumArt(state.Artist, state.Album)

	if err != nil {
		// Set to default art
		state.Url = "https://static-00.iconduck.com/assets.00/apple-music-icon-1024x1024-zncv5jwr.png"
	}

	setUrlCache(state.Artist, state.Album, state.Url)

	return state.Url, nil
}

var previousState playerState

func newEvent() musicEvent {
	return musicEvent{false, false, false, false}
}

func eventDetector() (playerState, musicEvent, error) {
	if DataCancel != nil {
		DataCancel()
	}

	DataCtx, DataCancel = context.WithCancel(context.Background())

	currentState, err := getPlayerState(DataCtx)
	defer DataCancel()

	if err != nil {
		return playerState{}, musicEvent{}, err
	}

	event := newEvent()

	// Detect song change
	didChange := previousState.Title != currentState.Title && previousState.Artist != currentState.Artist && previousState.Album != currentState.Album
	if didChange {
		event.songChanged = true
	}

	// Detect song state change
	stateChanged := previousState.State != currentState.State
	if stateChanged {
		event.stateChanged = true
	}

	// Detect change in track player location
	playheadMoved := previousState.Playtime.Sub(currentState.Playtime).Abs().Milliseconds() > 1000
	if playheadMoved {
		event.playheadChanged = true
	}

	// Detect if no track is playing
	if strings.TrimSpace(strings.ToLower(currentState.Title)) == "no track playing" {
		event.noTrackPlaying = true
	} else {
		event.noTrackPlaying = false
	}

	previousState = currentState

	return currentState, event, nil
}

func getPlayerState(ctx context.Context) (playerState, error) {
	cliCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	script := `
	set jsonResult to ""
	try
	    tell application "Music"
	        if player state is not stopped then
	            set playerState to player state
	            set playheadTime to player position
	            set trackName to name of current track
	            set artistName to artist of current track
	            set albumName to album of current track
				set trackLength to duration of current track
	        else
	            set playerState to "stopped"
	            set playheadTime to 0
	            set trackName to "No track playing"
	            set artistName to "No artist"
	            set albumName to "No album"
				set trackLength to 0
	        end if
	    end tell

	    set playheadTimeFormatted to (playheadTime as string)

	    set jsonResult to "{ \"track_name\": \"" & trackName & "\", " & ¬
	        "\"artist_name\": \"" & artistName & "\", " & ¬
	        "\"album_name\": \"" & albumName & "\", " & ¬
	        "\"player_state\": \"" & playerState & "\", " & ¬
			"\"track_length\": \"" & trackLength & "\", " & ¬
	        "\"playhead_time\": \"" & playheadTimeFormatted & "\" }"
	on error errMsg
	    set jsonResult to "{ \"error\": \"" & errMsg & "\" }"
	end try

	return jsonResult
    `

	// cmd := exec.Command("osascript", "-e", script)
	cmd := exec.CommandContext(cliCtx, "osascript", "-e", script)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return playerState{}, err
	}

	var parsedTrack playerState
	output := strings.TrimSpace(out.String())
	err = json.Unmarshal([]byte(output), &parsedTrack)

	if err != nil {
		return playerState{}, err
	}

	// Convert playhead time
	parsedTrack.Playtime, err = getPlayheadTime(parsedTrack.Playhead)

	if err != nil {
		return playerState{}, err
	}

	if strings.TrimSpace(parsedTrack.State) == "playing" {
		parsedTrack.isPlaying = true
	} else {
		parsedTrack.isPlaying = false
	}

	// Grab URL
	albumUrl, err := getAlbumArtUrl(parsedTrack)

	if err == nil {
		parsedTrack.Url = albumUrl
	}

	return parsedTrack, nil
}

func getPlayheadTime(time_float string) (time.Time, error) {
	parsedFloat, err := strconv.ParseFloat(time_float, 32)

	if err != nil {
		return time.Time{}, err
	}

	time_stamp := time.Now().Add(-time.Duration(parsedFloat) * time.Second)
	return time_stamp, nil
}
