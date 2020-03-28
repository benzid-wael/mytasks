package entities

import (
	"time"
)

type Manageable interface {
	GetTitle() string
	GetTags() []string
	GetType() string
	Star()
	Unstar()
}

type Item struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsStarred   bool      `json:"is_starred"`
	Tags        []string  `json:"tags"`
	Type        string    `json:"type"`
}

type Note struct {
	Item
}

func newItem(title string, kind string, description string, tags ...string) *Item {
	return &Item{
		Id:          -1,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		IsStarred:   false,
		Tags:        tags,
		Type:        kind,
	}
}

func NewNote(title string, description string, tags ...string) *Note {
	return &Note{Item: *newItem(title, "note", description, tags...)}
}

func (item *Item) GetTitle() string {
	return item.Title
}

func (item *Item) GetTags() []string {
	return item.Tags
}

func (item *Item) GetType() string {
	return item.Type
}

func (item *Item) Star() {
	item.IsStarred = true
}

func (item *Item) Unstar() {
	item.IsStarred = false
}
