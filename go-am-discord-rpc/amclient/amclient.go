package amclient

import (
	"fmt"
	"sync"
)

func NewClient() {
	fmt.Println("Initializing...")

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
