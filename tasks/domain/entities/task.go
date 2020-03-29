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
	Cancelled  TaskStatus = "cancelled"
)

type Task struct {
	Item
	Status   TaskStatus `json:"status"`
	Priority uint       `json:"priority"`
	fsm      *fsm.FSM
}

func (task *Task) GetStatus() string {
	return string(task.Status)
}

func NewTask(title string, description string, status TaskStatus, tags ...string) *Task {
	task := &Task{
		Item:     *newItem(title, "task", description, tags...),
		Status:   status,
		Priority: 0,
	}

	task.fsm = fsm.NewFSM(
		string(status),
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
	err := task.fsm.Event(event, args...)
	if err != nil {
		task.Status = TaskStatus(task.fsm.Current())
	}
	return err
}
