package main

import (
	"flag"

	"github.com/SlyMarbo/rss"
)

var feedURL = flag.String("url", "https://site.com/rss", "Drupal RSS feed URL")

// PressRelease struct type
type PressRelease struct {
	LegacyNode
}

// SetValues handles field values of PressRelease
func (pr *PressRelease) SetValues(item *rss.Item) error {
	pr.Title = append(
		pr.Title,
		Field{
			Value: item.Title,
		},
	)

	body := TextFieldWithSummary{}
	body.Format = "basic_html"
	body.Summary = ""

	pr.Body = append(pr.Body, body)

	return nil
}

func (pr *PressRelease) Save() error {
	return nil
}

func main() {
	flag.Parse()

	svc := DrupalService{
		feedURL:     *feedURL,
		registry:    make(chan *rss.Item),
		nodeService: &PressRelease{},
	}

	go svc.Fetch()

	if err := svc.Save(); err != nil {
		panic(err)
	}
}
