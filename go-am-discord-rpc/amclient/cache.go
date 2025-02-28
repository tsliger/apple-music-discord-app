package amclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/allegro/bigcache/v3"
)

var cache *bigcache.BigCache
var loadedCache bool = false

func createCache() {
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(30*time.Minute))

	combinedPath, err := getFolderPath()

	if err == nil {
		loadCacheFile(combinedPath)
	}
}

func loadCacheFile(filename string) error {
	if cache == nil {
		return errors.New("Cache not initialized.")
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var cacheData map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cacheData); err != nil {
		return err
	}

	for key, value := range cacheData {
		if byteString, ok := value.(string); ok {
			decodedBytes, err := base64.StdEncoding.DecodeString(byteString)

			fmt.Println(string(decodedBytes))

			if err != nil {
				cache.Set(key, decodedBytes)
			}
		}
	}

	loadedCache = true
	return nil
}

func saveCacheFile(filename string) error {
	if loadedCache == false {
		return errors.New("Cache never loaded, won't overwrite previous cache.")
	}

	cacheData := make(map[string]interface{})

	iterator := cache.Iterator()
	for iterator.SetNext() {
		current, err := iterator.Value()

		if err == nil {
			cacheData[current.Key()] = current.Value()
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cacheData); err != nil {
		return fmt.Errorf("failed to encode cache data: %w", err)
	}

	return nil
}

func setUrlCache(artist string, album string, url string) error {
	cache_key := gen_key(artist, album)

	err := cache.Set(cache_key, []byte(url))

	if err != nil {
		fmt.Println("Could not cache url")
	}

	return err
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

func getFolderPath() (string, error) {
	homePath, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	combinedPath := filepath.Join(homePath, "/Library/Application Support/applemusicrpc/", "artwork_cache.json")

	return combinedPath, nil
}

func cleanCache() error {
	combinedPath, err := getFolderPath()

	if err != nil {
		return err
	}

	saveCacheFile(combinedPath)

	return nil
}
