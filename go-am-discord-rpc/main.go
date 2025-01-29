package main

import (
	"am-discord-rpc/amclient"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var ctx context.Context
var cancel context.CancelFunc

func pollingProcess() {
	for {
		select {
		case <-ctx.Done():
			amclient.CloseDiscordClient()
			fmt.Println("Process was cancelled.")
			return
		default:
			// Simulate some ongoing work
			fmt.Println("Process is running...")
			amclient.Poll()
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// Initialize client
	amclient.NewClient()

	r := gin.Default()

	ctx, cancel = context.WithCancel(context.Background())

	go pollingProcess()

	r.GET("/kill", func(c *gin.Context) {
		if cancel != nil {
			cancel()
			os.Exit(0)
			c.JSON(http.StatusOK, gin.H{
				"message": "Killed process.",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "No process started.",
			})
		}
	})

	r.Run()

	defer cancel()
}
