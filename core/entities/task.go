package entities

import (
	"github.com/looplab/fsm"
)

type TaskStatus string

const (
	ToDo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Stopped    TaskStatus = "stopped"
	Completed  TaskStatus = "completed"
	NotDoing   TaskStatus = "not-doing"
	Cancelled  TaskStatus = "cancelled"
)

type Task struct {
	Item
	Status   TaskStatus `json:"status"`
	Priority int8       `json:"priority"`
	fsm      *fsm.FSM
}

func NewTask(id *Sequence, title string, description string, tags ...string) *Task {
	return &Task{
		Item:     *newItem(id, title, description, tags...),
		Status:   ToDo,
		Priority: 0,
		fsm:      nil,
	}
}
