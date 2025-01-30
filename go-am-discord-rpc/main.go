package main

import (
	"am-discord-rpc/amclient"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var ctx context.Context
var cancel context.CancelFunc

func pollingProcess() {
	for {
		select {
		case <-ctx.Done():
			amclient.CloseClient()
			fmt.Println("Process was cancelled.")
			return
		default:
			amclient.Poll()
		}
	}
}

func main() {
	// Init on termination signal detection
	createTerminator()

	// Initialize client
	amclient.NewClient()
	ctx, cancel = context.WithCancel(context.Background())

	var pString string
	_, err := fmt.Scan(&pString)

	if err != nil {
		fmt.Println("Exiting due to port parsing issue.")
		AppExit()
	}

	port := ":" + pString

	r := gin.Default()

	go pollingProcess()

	r.GET("/kill", func(c *gin.Context) {
		if cancel != nil {
			AppExit()
			c.JSON(http.StatusOK, gin.H{
				"message": "Killed process.",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "No process started.",
			})
		}
	})

	err = r.Run(port)

	if err != nil {
		AppExit()
	}

	defer cancel()
}

func AppExit() {
	amclient.CloseClient()
	cancel()
	os.Exit(0)
}
