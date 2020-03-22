package entities

import (
	"time"
)

type Item struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsStarred   bool      `json:"is_starred"`
	Tags        []string  `json:"tags"`
}

type Note struct {
	Item
}

func newItem(id *Sequence, title string, description string, tags ...string) *Item {
	id.Next()
	return &Item{
		Id:          id.Current(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		IsStarred:   false,
		Tags:        tags,
	}
}

func NewNote(id *Sequence, title string, description string, tags ...string) *Note {
	return &Note{Item: *newItem(id, title, description, tags...)}
}

func (item *Item) Star() {
	item.IsStarred = true
}

func (item *Item) Unstar() {
	item.IsStarred = false
}
