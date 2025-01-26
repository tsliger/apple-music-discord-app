package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var currentTrack CurrentTrack
var playingTrack string

func poll() {
	var previousState string

	for {
		currentState, err := getPlayerState()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}

		currentState = strings.TrimSpace(currentState)

		// Check if the state has changed
		if currentState != previousState {
			fmt.Printf("State changed: %s\n", currentState)
			previousState = currentState
		}

		if currentState != "paused" {
			currTrack, err := getCurrentTrack()

			if err == nil {
				setActivity(currTrack)
			}
		}

		// Poll every second
		time.Sleep(1 * time.Second)
	}
}

func getCurrentTrack() (CurrentTrack, error) {
	script := `
	tell application "Music"
	    set currentTrack to current track
	    set trackName to name of currentTrack
	    set artistName to artist of currentTrack
	    set albumName to album of currentTrack
	    set playheadTime to player position
	    set endTime to duration of currentTrack
	end tell

	set playheadTimeFormatted to (playheadTime as string)
	set endTimeFormatted to (endTime as string)

	set jsonResult to "{ \"track_name\": \"" & trackName & "\", " & ¬
	    "\"artist_name\": \"" & artistName & "\", " & ¬
	    "\"album_name\": \"" & albumName & "\", " & ¬
		"\"playhead_time\": \"" & playheadTimeFormatted & "\", " & ¬
	    "\"end_time\": \"" & endTimeFormatted & "\" }"

	return jsonResult
    `

	cmd := exec.Command("osascript", "-e", script)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return CurrentTrack{}, err
	}

	output := strings.TrimSpace(out.String())
	err = json.Unmarshal([]byte(output), &currentTrack)

	if err != nil {
		fmt.Println("failed to unmarshal json")
	}

	currentTrack.Playtime, _ = getPlayheadTime(currentTrack.Playhead)

	currTrack := currentTrack.Title + " by " + currentTrack.Artist

	if currTrack != playingTrack {
		playingTrack = currTrack
		fmt.Println(playingTrack)

		// Grab album art from cache or scrape
		cached_url, err := GetUrlFromCache(currentTrack.Artist, currentTrack.Album)

		if err == nil {
			currentTrack.Url = cached_url
		} else {
			scraped_url, err := scrapeAlbumArt(currentTrack.Artist, currentTrack.Album)

			// use default url
			if err != nil {

			} else {
				// Set cache to scraped url
				SetUrlCache(currentTrack.Artist, currentTrack.Album, scraped_url)
				currentTrack.Url = scraped_url
			}
		}
	}

	return currentTrack, nil
}

func getPlayerState() (string, error) {
	// AppleScript to get the player state
	script := `
        tell application "Music"
            return player state as string
        end tell
    `
	cmd := exec.Command("osascript", "-e", script)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	// Return the trimmed player state
	return out.String(), nil
}

func getPlayheadTime(time_float string) (time.Time, error) {
	parsedFloat, err := strconv.ParseFloat(time_float, 32)

	if err != nil {
		return time.Time{}, err
	}

	time_stamp := time.Now().Add(-time.Duration(parsedFloat) * time.Second)
	return time_stamp, nil
}
