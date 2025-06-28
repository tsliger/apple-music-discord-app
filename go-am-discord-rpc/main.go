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

func pollingProcess(ctx context.Context) {
	amclient.Poll(ctx)
	fmt.Println("Polling process exited.")
}

func main() {
	amclient.CreateScraper()
	amclient.NewClient()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) < 2 {
		fmt.Println("Exiting: no port argument provided.")
		return
	}
	port := ":" + os.Args[1]

	go pollingProcess(ctx)

	// Gin setup
	r := gin.Default()

	r.GET("/kill", func(c *gin.Context) {
		cancel()
		c.JSON(http.StatusOK, gin.H{"message": "Killed process, shutting down..."})
	})

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	// Signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		fmt.Println("\nReceived OS interrupt signal")
		amclient.CleanScraper()
		cancel()
	}()

	// Wait for context to be canceled
	<-ctx.Done()

	// Graceful shutdown
	fmt.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	fmt.Println("Exited cleanly.")
}
