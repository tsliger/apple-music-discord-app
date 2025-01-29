package main

import "time"

type PlayerState struct {
	State     string `json:"player_state"`
	Title     string `json:"track_name"`
	Artist    string `json:"artist_name"`
	Album     string `json:"album_name"`
	Playhead  string `json:"playhead_time"`
	Playtime  time.Time
	Url       string
	isPlaying bool
}

type MusicEvent struct {
	songChanged     bool
	stateChanged    bool
	playheadChanged bool
	noTrackPlaying  bool
}
