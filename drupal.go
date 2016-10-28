package main

import (
	"time"

	"github.com/SlyMarbo/rss"
)

// Field struct type
type Field struct {
	Value string `json:"value"`
}

type FieldWithAlias struct {
	Field

	Alias string `json:"alias"`
}

// RefField struct type
type RefField struct {
	ID   string `json:"target_id"`
	Type string `json:"target_type"`
	UUID string `json:"target_uuid"`
}

// RefWithURLField struct type
type RefWithURLField struct {
	RefField

	URL string `json:"url"`
}

// FileField struct type
type FileField struct {
	RefWithURLField

	Display     string `json:"display"`
	Description string `json:"description"`
}

// TextField struct type
type TextField struct {
	Field

	Format string `json:"format"`
}

// TextFieldWithSummary struct type
type TextFieldWithSummary struct {
	TextField

	Summary string `json:"summary"`
}

// Node struct type
type Node struct {
	Type            []RefField             `json:"type"`
	Title           []Field                `json:"title"`
	Body            []TextFieldWithSummary `json:"body"`
	UID             []RefWithURLField      `json:"uid"`
	Status          []Field                `json:"status"`
	Created         []Field                `json:"created"`
	Path            []FieldWithAlias       `json:"path"`
	BackgroundImage []FileField            `json:"field_background_image"`
	Company         []Field                `json:"field_company"`
	Sidebar         []TextField            `json:"field_sidebar"`
	DatePublished   []Field                `json:"field_date_published"`
	Language        []Field                `json:"field_language"`
}

// LegacyNode struct type contains legacy fields
type LegacyNode struct {
	Node

	Content          []TextFieldWithSummary `json:"field_press_release_content"`
	SidebarOld       []TextField            `json:"field_press_release_sidebar"`
	DatePublishedOld []Field                `json:"field_press_rel_date_published"`
	LanguageOld      []Field                `json:"field_press_rel_language"`
}

// NodeService handle setting of values and saving item
type NodeService interface {
	Save(item *rss.Item) error
}

// DrupalService struct type.
type DrupalService struct {
	registry    chan *rss.Item
	nodeService NodeService
	fetched     chan bool
	done        chan bool
	err         chan error
}

func (s *DrupalService) Done() chan bool {
	return s.done
}

func (s *DrupalService) Error() chan error {
	return s.err
}

// Fetch retrieves data from Drupal RSS feed.
func (s *DrupalService) Fetch(fetchFunc func() (*rss.Feed, error)) {
	feed, err := fetchFunc()

	if err != nil {
		s.err <- err
		return
	}

	for _, item := range feed.Items {
		select {
		case s.registry <- item:
		}
	}
	s.fetched <- true
}

// Save saves RSS item to new Drupal website.
func (s DrupalService) Save() {
	fetched := false
	for {
		select {
		case item := <-s.registry:
			if err := s.nodeService.Save(item); err != nil {
				s.err <- err
			}
		case fetched = <-s.fetched:
		case <-time.After(5 * time.Second):
			if fetched {
				s.done <- true
			}
		}
	}
}
