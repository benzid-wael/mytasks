package entities

import "time"

type TaskStatus string

const (
	ToDo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Stopped    TaskStatus = "stopped"
	Completed  TaskStatus = "completed"
	NotDoing   TaskStatus = "not-doing"
	Cancelled  TaskStatus = "cancelled"
)

type Item struct {
	Id          int64     `json: "id"`
	Title       string    `json: "title"`
	Description string    `json: "description"`
	CreatedAt   time.Time `json: "created_at"`
	IsStarred   bool      `json: "is_starred"`
	Tags        []string  `json: "tags"`
}

type Note Item

type Task struct {
	Item
	status   TaskStatus `json: "status"`
	priority int8       `json:"priority"`
}
