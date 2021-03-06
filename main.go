package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/SlyMarbo/rss"
	"github.com/pkg/errors"
)

var feedURL = flag.String("feedUrl", "https://site.com/rss", "Drupal RSS feed URL")
var apiURL = flag.String("apiUrl", "http://127.0.0.1:8088/entity/node?_format=json", "Drupal API URL")
var credentials = flag.String("authbasic", "admin:admin", "Basic Auth credentials in username:password format")

// PressRelease struct type
type PressRelease struct {
	Node

	BackgroundImage []FileField `json:"field_background_image"`
	Company         []Field     `json:"field_company"`
	Sidebar         []TextField `json:"field_sidebar"`
	DatePublished   []Field     `json:"field_date_published"`
	Language        []Field     `json:"field_language"`
}

// SetValues handles field values of PressRelease
func (pr *PressRelease) SetValues(item *rss.Item) error {
	fmt.Println("Preparing:", item.Title)
	// Set title
	pr.Title = []Field{{
		Value: item.Title,
	}}

	// Set body
	body := TextFieldWithSummary{}
	body.Format = "basic_html"
	body.Summary = ""
	body.Value = item.Summary

	pr.Body = []TextFieldWithSummary{body}

	// Set content (legacy)
	content := TextFieldWithSummary{}
	content.Format = "basic_html"
	content.Summary = ""
	content.Value = "Intentionally empty"

	// Set type
	pr.Type = []RefField{{
		ID:   "press_release",
		Type: "node_type",
		UUID: "04b06708-ade5-49ec-a71f-5ae9d834a35f",
	}}

	// Set user ID
	uid := RefWithURLField{}
	uid.ID = "1"
	uid.Type = "user"
	uid.UUID = "a4c2373c-1b55-4c44-baf0-e40555c8ba7a"
	uid.URL = "/user/1"

	pr.UID = []RefWithURLField{uid}

	// Set status
	pr.Status = []Field{{Value: "1"}}

	// Set publish date and time
	datePublished := fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d",
		item.Date.Year(),
		item.Date.Month(),
		item.Date.Day(),
		item.Date.Hour(),
		item.Date.Minute(),
		item.Date.Second(),
	)
	pr.DatePublished = []Field{{Value: datePublished}}

	// Set date created
	pr.Created = []Field{{Value: fmt.Sprintf("%d", item.Date.Unix())}}

	// Set company
	pr.Company = []Field{{Value: item.Category}}

	// Set path
	path := FieldWithAlias{}

	u, err := url.Parse(*feedURL)
	if err != nil {
		return err
	}

	path.Alias = strings.Replace(item.Link, fmt.Sprintf("%s://%s", u.Scheme, u.Host), "", 1)
	pr.Path = []FieldWithAlias{path}

	return nil
}

// Save handle save of press releases to the API endpoint.
func (pr PressRelease) Save(item *rss.Item) error {
	pr.SetValues(item)

	b, err := json.Marshal(pr)
	if err != nil {
		return err
	}

	// Some parts is generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("POST", *apiURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	basicAuth := strings.Split(*credentials, ":")
	req.SetBasicAuth(basicAuth[0], basicAuth[1])

	fmt.Println("Saving:", pr.Title[0].Value)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}
	defer resp.Body.Close()
	fmt.Println(pr.Title[0].Value, "Saved!")
	return nil
}

func main() {
	flag.Parse()

	svc := New()

	go svc.Fetch(func() (*rss.Feed, error) {
		return rss.Fetch(*feedURL)
	}, 0)

	// FIXME: Check semaphore.
	for _ = range make([]int, runtime.NumCPU()) {
		go svc.Save()
	}

	for {
		select {
		case err := <-svc.Error():
			// TODO: Don't panic, use logging.
			panic(err)
		case done := <-svc.Done():
			if done {
				os.Exit(0)
			}
		}
	}

}
