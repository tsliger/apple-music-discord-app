package main

import "github.com/altfoxie/drpc"

var client *drpc.Client
var DISCORD_APP_ID string = "1332158263432708146"

func setActivity(info CurrentTrack) {
	err := client.SetActivity(drpc.Activity{
		State:   info.Artist,
		Details: info.Title,
		Timestamps: &drpc.Timestamps{
			Start: info.Playtime,
		},
		Assets: &drpc.Assets{
			LargeImage: info.Url,
			LargeText:  info.Album,
			SmallText:  "",
			SmallImage: "https://static-00.iconduck.com/assets.00/apple-music-icon-1024x1024-zncv5jwr.png",
		},
	})

	if err != nil { /* handle error */
	}
}

func initializeDiscord() {
	var err error
	client, err = drpc.New(DISCORD_APP_ID)

	if err != nil {
	}
}
