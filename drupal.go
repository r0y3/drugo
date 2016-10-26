package main

import "github.com/SlyMarbo/rss"

// Field struct type
type Field struct {
	Value string `json:"value"`
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
	NID             []Field                `json:"nid"`
	UUID            []Field                `json:"uuid"`
	VID             []Field                `json:"vid"`
	Langcode        []Field                `json:"langcode"`
	Type            []RefField             `json:"type"`
	Title           []Field                `json:"title"`
	Body            []TextFieldWithSummary `json:"body"`
	UID             []RefWithURLField      `json:"uid"`
	Status          []Field                `json:"status"`
	Created         []Field                `json:"created"`
	Path            []Field                `json:"path"`
	BackgroundImage []FileField            `json:"field_background_image"`
	Company         []Field                `json:"field_company"`
	Layout          []Field                `json:"field_layout"`
	PDFDE           []FileField            `json:"field_pdf_de"`
	PDFEN           []FileField            `json:"field_pdf_en"`
	PDFFR           []FileField            `json:"field_pdf_fr"`
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
	SetValues(item *rss.Item) error
	Save() error
}

// DrupalService struct type.
type DrupalService struct {
	feedURL     string
	registry    chan *rss.Item
	nodeService NodeService
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
			if err := s.nodeService.SetValues(item); err != nil {
				return err
			}
			if err := s.nodeService.Save(); err != nil {
				return err
			}
		}
	}
	return nil
}
