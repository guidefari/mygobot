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
}
