package amclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

func Poll() {
	ticker := time.NewTicker(songPollingRate * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			playingState, event, err := eventDetector()

			if err != nil {
				continue
			}

			eventHandler(event, playingState)
		}
	}
}

var previousState playerState
var start time.Time = time.Now().Add(-DISCORD_RATE_DELAY * time.Millisecond)

func eventHandler(event musicEvent, state playerState) {
	// Set activity based off event changes
	if event.noTrackPlaying {
		state.Url = DEFAULT_ALBUM_URI
	}

	if event.albumArtUpdated || event.stateChanged || event.songChanged || event.playheadChanged {
		if time.Since(start) >= DISCORD_RATE_DELAY*time.Millisecond {

			if state != previousState {
				setDiscordActivity(state)
				previousState = state
			}
			start = time.Now() // Update start only after the delay condition is met
		}
	}
}

func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

var set mapset.Set[string] = mapset.NewSet[string]()

func getAlbumArtUrl(state playerState) (string, error) {
	cachedUrl, err := getUrlFromCache(state.Artist, state.Album)
	if err == nil {
		// Cache hit, return the cached URL
		return cachedUrl, nil
	}

	contains := set.Contains(state.Artist + state.Album)
	if contains {
		return DEFAULT_ALBUM_URI, nil
	} else {
		// q = append(q, state.Artist+state.Album)
		set.Add(state.Artist + state.Album)
		go func() {
			albumArtUrl, err := scrapeAlbumArt(state.Artist, state.Album)
			if err != nil {
				// If scraping fails, fall back to a default album art URL
				albumArtUrl = DEFAULT_ALBUM_URI
			}

			if err := setUrlCache(state.Artist, state.Album, albumArtUrl); err != nil {
				// Log or handle the error if caching fails
				// This should ideally be non-blocking if it's not crucial
				fmt.Printf("Failed to cache album art URL for %s - %s: %v\n", state.Artist, state.Album, err)
			}
			set.Remove(state.Artist + state.Album)
		}()
	}

	return DEFAULT_ALBUM_URI, nil
}

func newEvent() musicEvent {
	return musicEvent{false, false, false, false, false}
}

func eventDetector() (playerState, musicEvent, error) {
	currentState, err := getPlayerState()

	if err != nil {
		return playerState{}, musicEvent{}, err
	}

	event := newEvent()

	if hasSongChanged(currentPlayerState, currentState) {
		event.songChanged = true
	}

	if currentPlayerState.State != currentState.State {
		event.stateChanged = true
	}

	if playheadMoved(currentPlayerState.Playtime, currentState.Playtime) {
		event.playheadChanged = playheadMoved(currentPlayerState.Playtime, currentState.Playtime)
	}

	if updatedAlbumArt(currentPlayerState, currentState) {
		event.albumArtUpdated = true
	}

	event.noTrackPlaying = isNoTrackPlaying(currentState.Title)

	return currentState, event, nil
}

func updatedAlbumArt(prevState, currState playerState) bool {
	return prevState.Url != currState.Url
}

func hasSongChanged(prevState, currState playerState) bool {
	return prevState.Title != currState.Title ||
		prevState.Artist != currState.Artist ||
		prevState.Album != currState.Album
}

func playheadMoved(prevPlaytime time.Time, currPlaytime time.Time) bool {
	return prevPlaytime.Sub(currPlaytime).Abs().Seconds() > 2
}

func isNoTrackPlaying(title string) bool {
	return strings.TrimSpace(strings.ToLower(title)) == "no track playing"
}

func getPlayerState() (playerState, error) {
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

	cmd := exec.Command("osascript", "-e", script)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return playerState{}, fmt.Errorf("failed to execute AppleScript: %v", err)
	}

	var parsedTrack playerState
	output := strings.TrimSpace(out.String())
	err = json.Unmarshal([]byte(output), &parsedTrack)
	if err != nil {
		return playerState{}, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	// Convert playhead time
	parsedTrack.Playtime, err = getPlayheadTime(parsedTrack.Playhead)
	if err != nil {
		return playerState{}, fmt.Errorf("failed to get playhead time: %v", err)
	}

	parsedTrack.isPlaying = strings.TrimSpace(parsedTrack.State) == "playing"

	// Grab URL
	if albumUrl, err := getAlbumArtUrl(parsedTrack); err == nil {
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
