package amclient

import (
	"sync"
)

func NewClient() {
	var wg sync.WaitGroup
	wg.Add(2)

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
