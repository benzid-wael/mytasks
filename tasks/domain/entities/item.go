package entities

import (
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
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

func newItem(id *value_objects.Sequence, title string, kind string, description string, tags ...string) *Item {
	id.Next()
	return &Item{
		Id:          id.Current(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		IsStarred:   false,
		Tags:        tags,
		Type:        kind,
	}
}

func NewNote(id *value_objects.Sequence, title string, description string, tags ...string) *Note {
	return &Note{Item: *newItem(id, title, "note", description, tags...)}
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
