package amclient

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
)

var cache *bigcache.BigCache

func createCache() {
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(30*time.Minute))
}

func setUrlCache(artist string, album string, url string) {
	cache_key := gen_key(artist, album)

	err := cache.Set(cache_key, []byte(url))

	if err != nil {
		fmt.Println("Could not cache url")
	}
}

func getUrlFromCache(artist string, album string) (string, error) {
	cache_key := gen_key(artist, album)

	entry, err := cache.Get(cache_key)

	if err != nil {
		return "", err
	}

	return string(entry), nil
}

func gen_key(artist string, album string) string {
	return artist + album
}
