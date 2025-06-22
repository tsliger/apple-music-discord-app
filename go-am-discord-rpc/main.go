package main

import (
	"am-discord-rpc/amclient"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var ctx context.Context
var cancel context.CancelFunc

func pollingProcess(ctx context.Context) {
	amclient.Poll(ctx) // This runs until context is cancelled
	fmt.Println("Polling process exited.")
}

func main() {
	amclient.CreateScraper()
	amclient.NewClient()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Setup OS signal handling (Ctrl+C, SIGTERM)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nReceived OS interrupt signal")
		cancel()
	}()

	go pollingProcess(ctx)

	var pString string
	_, err := fmt.Scan(&pString)
	if err != nil {
		fmt.Println("Exiting due to port parsing issue.")
		cancel()
		return
	}

	port := ":" + pString

	r := gin.Default()

	r.GET("/kill", func(c *gin.Context) {
		cancel() // signal cancellation
		c.JSON(http.StatusOK, gin.H{"message": "Killed process, shutting down..."})
	})

	// Run Gin server in main goroutine
	if err := r.Run(port); err != nil && err != http.ErrServerClosed {
		fmt.Println("Gin server error:", err)
		cancel()
	}

	// Wait a moment to allow goroutines to exit cleanly
	time.Sleep(2 * time.Second)
	fmt.Println("Exiting main")
}
