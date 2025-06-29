package amclient

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

// func Poll(ctx context.Context) {
// 	ticker := time.NewTicker(songPollingRate * time.Millisecond)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Poll stopped")
// 			return
// 		case <-ticker.C:
// 			fmt.Println("test")

// 			test := playerState{
// 				Artist:      "test",
// 				Album:       "test",
// 				isPlaying:   true,
// 				State:       "play",
// 				TrackLength: "100",
// 				Playhead:    "20",
// 				Playtime:    time.Now().UTC(),
// 				Title:       "test",
// 			}

// 			setDiscordActivity(test)
// 		}
// 	}
// }

func onNewOutput(line string) {
	fmt.Println("New output detected:", line)
	var parsedTrack playerState
	output := strings.TrimSpace(line)
	err := json.Unmarshal([]byte(output), &parsedTrack)

	if err != nil {
		return
		// return playerState{}, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	if albumUrl, err := getAlbumArtUrl(parsedTrack); err == nil {
		parsedTrack.Url = albumUrl
	}

	if parsedTrack.State == "Playing" {
		parsedTrack.isPlaying = true
		parsedTrack.Playtime, err = getPlayheadTime(parsedTrack.Playhead)
	}

	fmt.Println(parsedTrack)
	setDiscordActivity(parsedTrack)
}

func getPlayheadTime(time_float string) (time.Time, error) {
	parsedFloat, err := strconv.ParseFloat(time_float, 32)

	if err != nil {
		return time.Time{}, err
	}

	time_stamp := time.Now().Add(-time.Duration(parsedFloat) * time.Second)
	return time_stamp, nil
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
		set.Add(state.Artist + state.Album)
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

		return albumArtUrl, nil
	}
}

func Poll(ctx context.Context) {
	cmd := exec.Command("./windows-apple-music-info.exe")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Kill process on context cancel
	go func() {
		<-ctx.Done()
		fmt.Println("Context cancelled: killing subprocess")
		_ = cmd.Process.Kill()
	}()

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		onNewOutput(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading stdout:", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Process exited with error:", err)
	}
}
