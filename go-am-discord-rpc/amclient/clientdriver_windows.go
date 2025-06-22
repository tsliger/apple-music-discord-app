package amclient

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
	// fmt.Println("New output detected:", line)
	var parsedTrack playerState
	output := strings.TrimSpace(line)
	err := json.Unmarshal([]byte(output), &parsedTrack)

	if err != nil {
		return
		// return playerState{}, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	if albumUrl, err := getAlbumArtUrl(parsedTrack); err == nil {
		parsedTrack.Url = albumUrl
		fmt.Println("Art scraped: " + albumUrl)
	}

	if parsedTrack.State == "Playing" {
		parsedTrack.isPlaying = true
		parsedTrack.TrackLength = "500"
		parsedTrack.Playhead = "0"
		parsedTrack.Playtime = time.Now()
		// parsedTrack.Url = "https://alamocitygolftrail.com/wp-content/uploads/2022/11/canstockphoto22402523-arcos-creator.com_-1024x1024-1.jpg"
	}

	fmt.Println(parsedTrack)
	setDiscordActivity(parsedTrack)
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

	return DEFAULT_ALBUM_URI, nil
}

func Poll(ctx context.Context) {
	cmd := exec.Command("amclient\\win_client\\windows-apple-music-info.exe")

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
