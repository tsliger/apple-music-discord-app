package amclient

import (
	"sync"
)

func NewClient() {
	var wg sync.WaitGroup
	wg.Add(2)

	// go func() {
	// 	sigCh := make(chan os.Signal, 1)
	// 	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	// 	<-sigCh
	// 	fmt.Println("\nReceived OS interrupt signal")
	// 	// CloseClient()
	// 	os.Exit(0)
	// }()

	go func() {
		defer wg.Done()
		createCache()
	}()

	go func() {
		defer wg.Done()
		initializeDiscord()
	}()

	wg.Wait()
}

func CloseClient() {
	closeDiscordClient()
	cleanScraper()
	cleanCache()
}
