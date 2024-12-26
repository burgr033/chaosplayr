package file

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// CacheFile represents the cached file location
const CacheFile = "podcast-hq.xml"

// MaxReadBytes defines how many bytes to read for metadata
const MaxReadBytes = 1024

// downloadFile gets the feed xml from media.ccc.de
func downloadFile(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error fetching URL: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body: %v\n", err)
		}
		return string(bodyBytes)
	}
	return ""
}

// readLastBuildDate extracts the <lastBuildDate> tag from the feed
func readLastBuildDate(url string) (time.Time, error) {
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Limit the read to MaxReadBytes to avoid downloading the full feed
	limitedReader := io.LimitReader(resp.Body, MaxReadBytes)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return time.Time{}, fmt.Errorf("error reading response body: %v", err)
	}

	bodyString := string(bodyBytes)

	// Find <lastBuildDate> tag in the feed
	startTag := "<lastBuildDate>"
	endTag := "</lastBuildDate>"
	startIdx := strings.Index(bodyString, startTag)
	endIdx := strings.Index(bodyString, endTag)

	if startIdx == -1 || endIdx == -1 {
		return time.Time{}, fmt.Errorf("<lastBuildDate> tag not found")
	}

	// Extract and parse the date
	dateStr := bodyString[startIdx+len(startTag) : endIdx]
	lastBuildDate, err := time.Parse(time.RFC1123, strings.TrimSpace(dateStr))
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing <lastBuildDate>: %v", err)
	}

	return lastBuildDate, nil
}

// getCachedLastBuildDate reads the last cached build date
func getCachedLastBuildDate() (time.Time, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting cache path: %v", err)
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, nil // No cache yet
		}
		return time.Time{}, fmt.Errorf("error reading cache: %v", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(data))
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing cached feed: %v", err)
	}

	return *feed.UpdatedParsed, nil
}

// GetFeed fetches and parses the feed, using caching to avoid unnecessary downloads
func GetFeed(url string) []*gofeed.Item {
	lastBuildDate, err := readLastBuildDate(url)
	if err != nil {
		log.Fatalf("error reading lastBuildDate: %v\n", err)
	}

	cachedLastBuildDate, err := getCachedLastBuildDate()
	if err != nil {
		log.Fatalf("error reading cached lastBuildDate: %v\n", err)
	}

	// Check if the feed has been updated
	if !lastBuildDate.After(cachedLastBuildDate) {
		cachePath, err := getCachePath()
		if err != nil {
			log.Fatalf("error getting cache path: %v\n", err)
		}

		data, err := os.ReadFile(cachePath)
		if err != nil {
			log.Fatalf("error reading cached feed: %v\n", err)
		}

		fp := gofeed.NewParser()
		feed, err := fp.ParseString(string(data))
		if err != nil {
			log.Fatalf("error parsing cached feed: %v\n", err)
		}

		return feed.Items
	}

	// Fetch the updated feed
	feedContent := downloadFile(url)

	// Update the cache
	err = updateCache(feedContent)
	if err != nil {
		log.Fatalf("error updating cache: %v\n", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(feedContent)
	if err != nil {
		log.Fatalf("error parsing feed: %v\n", err)
	}

	return feed.Items
}
