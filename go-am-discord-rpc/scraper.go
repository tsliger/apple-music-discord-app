package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
)

// var browser playwright.Browser
var ctx context.Context
var cancel context.CancelFunc

func initializeScraper() {
	ctx, cancel = chromedp.NewContext(context.Background())
	// defer cancel()
}

func scrapeAlbumArt(artist string, album string) (string, error) {
	// Check if multiple artists are being listed by &
	artists := strings.Split(artist, "&")
	album = strings.ReplaceAll(album, "#", "")
	album = strings.ReplaceAll(album, "&", "")

	if len(artists) > 0 {
		artist = artists[0]
	}

	searchString := fmt.Sprintf("%s %s", artist, album)
	searchURL := "https://music.apple.com/us/search?term=" + searchString

	var urls string
	var ok bool
	err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(".artwork-component"),
		// chromedp.Nodes(".artwork-component > picture > source", &nodes, chromedp.ByQueryAll),
		chromedp.AttributeValue(".artwork-component > picture > source", "srcset", &urls, &ok, chromedp.ByQueryAll),
	)
	if err != nil {
		fmt.Println("Error navigating to the website:", err)
	}

	re := regexp.MustCompile(`https?://[^,\s]+\.webp`)
	var url string

	// Find the first match
	matches := re.FindStringSubmatch(urls)
	if len(matches) > 0 {
		url = strings.TrimSpace(matches[0])
	} else {
		fmt.Println("No .webp URLs found")
	}

	fmt.Println(url)
	return url, nil
}
