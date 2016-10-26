package main

import (
	"flag"

	"github.com/SlyMarbo/rss"
)

var feedURL = flag.String("url", "https://site.com/rss", "Drupal RSS feed URL")

func main() {
	flag.Parse()

	svc := DrupalService{
		feedURL:  *feedURL,
		registry: make(chan *rss.Item),
	}

	go svc.Fetch()

	if err := svc.Save(); err != nil {
		panic(err)
	}
}
