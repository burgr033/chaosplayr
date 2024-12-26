package file

import (
	"io"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
)

// downloadFile gets the feed xml from media.ccc.de
func downloadFile(url string) string {
	resp, err := http.Get(url)
	if err != nil {
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		bodyString := string(bodyBytes)
		return bodyString
	}
	return ""
}

// get string from url and feed into Feed type
func GetFeed(url string) []*gofeed.Item {
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(downloadFile(url))
	if err != nil {
		log.Fatalf("error obtaining feed file %v\n", err)
	}
	return feed.Items
}
