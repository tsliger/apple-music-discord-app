package main

import "time"

type CurrentTrack struct {
	Title    string `json:"track_name"`
	Artist   string `json:"artist_name"`
	Album    string `json:"album_name"`
	Playhead string `json:"playhead_time"`
	Length   string `json:"end_time"`
	Url      string
	Playtime time.Time
}
