package main

import (
	"fmt"

	"github.com/SlyMarbo/rss"
)

// DrupalService struct type.
type DrupalService struct {
	feedURL string

	registry chan *rss.Item
}

// Fetch retrieves data from Drupal RSS feed.
func (s *DrupalService) Fetch() error {
	feed, err := rss.Fetch(s.feedURL)

	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		select {
		case s.registry <- item:
		}
	}

	return nil
}

// Save saves RSS item to new Drupal website.
func (s *DrupalService) Save() error {
	for {
		select {
		case item := <-s.registry:
			// TODO: Save to Drupal webservice.
			fmt.Println(item.Title)
			fmt.Println(item.Summary)
			fmt.Println(item.Content)
		}
	}
	return nil
}
