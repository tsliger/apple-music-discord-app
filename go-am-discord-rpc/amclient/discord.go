package amclient

import (
	"fmt"
	"strconv"
	"time"

	discordrpc "github.com/tsliger/go-discord-rpc"
)

var client *discordrpc.Client

func setDiscordActivity(info playerState) error {
	var err error

	floatVal, err := strconv.ParseFloat(info.TrackLength, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
	}

	playHeadFloat, err := strconv.ParseFloat(info.Playhead, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
	}

	endTime := time.Now().Unix() + int64(floatVal) - int64(playHeadFloat)

	var data discordrpc.ActivityData
	if info.isPlaying {
		data = discordrpc.ActivityData{
			State:   info.Artist,
			Type:    discordrpc.LISTENTING_TYPE,
			Details: info.Title,
			Assets: discordrpc.ActivityAssets{
				LargeImage: info.Url,
				LargeText:  info.Album,
				// SmallImage: smallImg,
			},
			Timestamps: discordrpc.ActivityTimestamp{
				Start: int(info.Playtime.Unix()),
				End:   int(endTime),
			},
		}
	} else {
		data = discordrpc.ActivityData{}
	}

	client.SendActivity(data)

	return err
}

func initializeDiscord() error {
	var err error
	client, err = discordrpc.NewClient(discord_APP_ID)

	return err
}

func closeDiscordClient() {
	client.CloseClient()
}
