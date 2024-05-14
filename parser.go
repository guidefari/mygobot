package main

import (
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Title string
	URL   string
}

func ParseRss() {
	feedList := [2]string{"https://goosebumps.fm/rss.xml", "http://musicforprogramming.net/rss.php"}
	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{Timeout: time.Second * 10}

	feed_items := make([]FeedItem, 1)

	for true {
		for k := 0; k < len(feedList); k++ {

		}

		// Regenrate the items list
		feed_items = make([]FeedItem, 1)

		// Wait 8 hours
		time.Sleep(28800 * time.Second)
	}
}
