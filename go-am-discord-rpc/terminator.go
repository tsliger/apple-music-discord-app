package main

import (
	"am-discord-rpc/amclient"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func createTerminator() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("Received signal: %s\n", sig)
		fmt.Println("Running cleanup tasks...")

		// Clean functions
		amclient.CloseClient()

		fmt.Println("Cleanup completed. Exiting now.")
		os.Exit(0)
	}()
}
