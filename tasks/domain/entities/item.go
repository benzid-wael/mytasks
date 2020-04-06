package entities

import (
	"time"
)

type Manageable interface {
	GetId() int
	GetTitle() string
	GetTags() []string
	GetType() string
	GetStatus() string
	GetCreationDateTime() time.Time
	HasStar() bool
	Star()
	Unstar()
	GetDueDate() *time.Time
	HasTag(tag string) bool
	HasAnyTag(tags ...string) bool
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

func (item *Item) GetId() int {
	return item.Id
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

func (item *Item) GetStatus() string {
	return ""
}

func (item *Item) GetCreationDateTime() time.Time {
	return item.CreatedAt
}

func (item *Item) HasStar() bool {
	return item.IsStarred
}

func (item *Item) Star() {
	item.IsStarred = true
}

func (item *Item) Unstar() {
	item.IsStarred = false
}

func (item *Item) GetDueDate() *time.Time {
	return nil
}

func (item *Item) HasTag(tag string) bool {
	tags := item.GetTags()
	for _, val := range tags {
		if val == tag {
			return true
		}
	}
	return false
}

func (item *Item) HasAnyTag(tags ...string) bool {
	for _, tag := range tags {
		if item.HasTag(tag) {
			return true
		}
	}
	return false
}
