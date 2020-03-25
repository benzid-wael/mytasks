package entities

import (
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
	"github.com/looplab/fsm"
)

type TaskStatus string

const (
	ToDo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Stopped    TaskStatus = "stopped"
	Completed  TaskStatus = "completed"
	Cancelled  TaskStatus = "cancelled"
)

type Task struct {
	Item
	Status   TaskStatus `json:"status"`
	Priority int8       `json:"priority"`
	fsm      *fsm.FSM
}

func NewTask(id *value_objects.Sequence, title string, description string, tags ...string) *Task {
	task := &Task{
		Item:   *newItem(id, title, "task", description, tags...),
		Status: ToDo,
	}

	task.fsm = fsm.NewFSM(
		string(ToDo),
		fsm.Events{
			{Name: "start", Src: []string{string(ToDo), string(Stopped), string(Cancelled)}, Dst: string(InProgress)},
			{Name: "stop", Src: []string{string(InProgress)}, Dst: string(ToDo)},
			{Name: "complete", Src: []string{string(InProgress)}, Dst: string(Completed)},
			{Name: "cancel", Src: []string{string(ToDo), string(Stopped), string(InProgress)}, Dst: string(Cancelled)},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { task.Status = TaskStatus(task.fsm.Current()) },
		},
	)
	return task
}

func (task *Task) TriggerEvent(event string, args ...interface{}) error {
	return task.fsm.Event(event, args...)
}
