package main

import (
	"am-discord-rpc/amclient"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateTerminator() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("Received signal: %s\n", sig)

		fmt.Println("Running cleanup tasks...")

		time.Sleep(1 * time.Second)

		// Clean functions
		amclient.CleanScraper()
		amclient.CloseDiscordClient()
		EndPolling()

		fmt.Println("Cleanup completed. Exiting now.")
		os.Exit(0)
	}()
}
