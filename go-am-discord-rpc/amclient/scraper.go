package amclient

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	ctx           context.Context
	allocCtx      context.Context
	allocCancel   context.CancelFunc
	scraperCancel context.CancelFunc
)

func CreateScraper() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)

	allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, scraperCancel = chromedp.NewContext(allocCtx)
}

func scrapeAlbumArt(artist string, album string) (string, error) {
	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()

	timeoutCtx, cancel := context.WithTimeout(tabCtx, 10*time.Second)
	defer cancel()

	album = strings.ReplaceAll(album, "#", "")
	album = strings.ReplaceAll(album, "&", "")
	artists := strings.Split(artist, "&")
	if len(artists) > 0 {
		artist = artists[0]
	}

	searchString := fmt.Sprintf("%s %s", artist, album)
	searchURL := "https://music.apple.com/us/search?term=" + searchString

	var urls string
	var ok bool
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(".artwork-component", chromedp.ByQuery),
		chromedp.AttributeValue(".artwork-component > picture > source", "srcset", &urls, &ok, chromedp.ByQuery),
	)
	if err != nil {
		return "", fmt.Errorf("error navigating: %w", err)
	}

	re := regexp.MustCompile(`https?://[^,\s]+\.webp`)
	matches := re.FindStringSubmatch(urls)

	if len(matches) == 0 {
		return "", fmt.Errorf("no .webp URLs found")
	}

	return strings.TrimSpace(matches[0]), nil
}

func cleanScraper() {
	CleanScraper()
}

func CleanScraper() {
	if scraperCancel != nil {
		scraperCancel() // ends browser tab
	}
	if allocCancel != nil {
		allocCancel() // shuts down Chrome
	}
}
