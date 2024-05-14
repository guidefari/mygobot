package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Title string
	URL   string
}

func ParseRSS() {
	feedList := [2]string{"https://goosebumps.fm/rss.xml", "http://musicforprogramming.net/rss.php"}
	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{Timeout: time.Second * 10}

	feed_items := make([]FeedItem, 1)

	for true {
		for k := 0; k < len(feedList); k++ {
			feed, err := feedParser.ParseURL(feedList[k])

			if err != nil {
				fmt.Println(err)
				// local_log.Printf("[WARN] FeedItem couldn't create since link or title is empty!. URL: %s", feedList[k])
				return
			}

			local_log.Printf("[INFO] RSS Parsing started for %s", feedList[k])
			items := feed.Items

			for i := 0; i < len(items); i++ {
				if items[i].Title != "" && items[i].Link != "" && !strings.Contains(readFile("feed_item.list"), items[i].Link) {
					feedItem := FeedItem{Title: items[i].Title, URL: items[i].Link}
					feed_items = append(feed_items, feedItem)
					local_log.Printf("[INFO] FeedItem is created. Title: %s Link: %s", items[i].Title, items[i].Link)

					// Write link to the file
					file, err := openOrCreateFile("feed_item.list")
					if err != nil {
						local_log.Println(err)
					}
					defer file.Close()
					if _, err := file.WriteString(items[i].Link + "\n"); err != nil {
						local_log.Fatal(err)
					}

					msg := "New thing posted 🚨: **" + items[i].Title + "**\n" + items[i].Link
					disco_session.ChannelMessageSend(preferred_channel_id, msg)
				} else {
					local_log.Printf("[WARN] FeedItem couldn't create since link or title is empty!. URL: %s", feedList[k])
				}
			}

		}

		// Regenrate the items list
		feed_items = make([]FeedItem, 1)

		// Wait 8 hours
		time.Sleep(28800 * time.Second)
	}
}
