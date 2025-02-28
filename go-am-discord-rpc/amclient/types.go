package amclient

import "time"

type playerState struct {
	State       string `json:"player_state"`
	Title       string `json:"track_name"`
	Artist      string `json:"artist_name"`
	Album       string `json:"album_name"`
	Playhead    string `json:"playhead_time"`
	TrackLength string `json:"track_length"`
	Playtime    time.Time
	Url         string
	isPlaying   bool
}

type musicEvent struct {
	songChanged     bool
	stateChanged    bool
	playheadChanged bool
	noTrackPlaying  bool
}
