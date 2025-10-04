package amclient

import (
	"fmt"
	"strconv"
	"time"

	discordrpc "github.com/tsliger/go-discord-rpc"
)

var client *discordrpc.Client
var currentPlayerState playerState

func setDiscordActivity(info playerState) error {
	currentPlayerState = info

	if info.isPlaying {
		floatVal, err := strconv.ParseFloat(info.TrackLength, 64)
		if err != nil {
			fmt.Println("Error converting string to float:", err)
		}

		playHeadFloat, err := strconv.ParseFloat(info.Playhead, 64)
		if err != nil {
			fmt.Println("Error converting string to float:", err)
		}

		endTime := time.Now().Unix() + int64(floatVal+0.5) - int64(playHeadFloat)

		data := &discordrpc.ActivityData{
			Name:    info.Artist,
			State:   info.Artist,
			Type:    discordrpc.LISTENTING_TYPE,
			Details: info.Title,
			Assets: discordrpc.ActivityAssets{
				LargeImage: info.Url,
				LargeText:  info.Album,
			},
			Timestamps: discordrpc.ActivityTimestamp{
				Start: int64(info.Playtime.Unix()),
				End:   int64(endTime),
			},
		}

		err = client.SetActivity(data)

		if err != nil {
			fmt.Printf("Error setting: %v", err)
			return err
		}
	} else {
		data := &discordrpc.ActivityData{
			Type:    discordrpc.LISTENTING_TYPE,
			Details: "Paused",
		}

		client.SetActivity(data)
	}

	return nil
}

func initializeDiscord() error {
	var err error
	client, err = discordrpc.NewClient(discord_APP_ID)

	return err
}

func closeDiscordClient() {
	client.CloseClient()
}
