package amclient

import "github.com/altfoxie/drpc"

var client *drpc.Client

const discord_APP_ID = "1332158263432708146"
const playImage = "https://i.ibb.co/5gW2VJLX/play.png"
const pauseImage = "https://i.ibb.co/RTC7L0zK/pause.png"

func setDiscordActivity(info playerState) error {
	// when album name is 1 character, there is an issue with drpc causing the activity to not be set
	info.Album += "      "

	var smallImg string
	if info.isPlaying {
		smallImg = playImage
	} else {
		smallImg = pauseImage
	}

	var err error
	if info.isPlaying {
		err = client.SetActivity(drpc.Activity{
			State:   info.Artist,
			Details: info.Title,
			Timestamps: &drpc.Timestamps{
				Start: info.Playtime,
			},
			Assets: &drpc.Assets{
				LargeImage: info.Url,
				LargeText:  info.Album,
				SmallText:  "",
				SmallImage: smallImg,
			},
		})
	} else {
		err = client.SetActivity(drpc.Activity{
			State:   info.Artist,
			Details: info.Title,
			Assets: &drpc.Assets{
				LargeImage: info.Url,
				LargeText:  info.Album,
				SmallText:  "",
				SmallImage: smallImg,
			},
		})
	}

	return err
}

func initializeDiscord() error {
	var err error
	client, err = drpc.New(discord_APP_ID)

	return err
}

func closeDiscordClient() {
	client.Close()
}
