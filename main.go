package main

import (
	"flag"

	"github.com/SlyMarbo/rss"
)

var feedURL = flag.String("url", "https://www.nagra.com/press-releases-xml-complete/Kudelski%20Group,Nagra%20Kudelski,SmarDTV,Kudelski%20Security", "Drupal RSS feed URL")

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
