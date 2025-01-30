package amclient

import (
	"sync"
)

func NewClient() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		CreateCache()
	}()

	go func() {
		defer wg.Done()
		initializeScraper()
	}()

	go func() {
		defer wg.Done()
		initializeDiscord()
	}()

	wg.Wait()
}
