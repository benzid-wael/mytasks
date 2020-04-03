package entities

import (
	"github.com/looplab/fsm"
	"time"
)

type TaskStatus string
type Priority uint

const (
	ToDo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Stopped    TaskStatus = "stopped"
	Completed  TaskStatus = "completed"
	Cancelled  TaskStatus = "cancelled"
)

const (
	UNKNOWN  Priority = 0
	TRIVIAL  Priority = 2
	LOW      Priority = 4
	MEDIUM   Priority = 8
	HIGH     Priority = 16
	CRITICAL Priority = 32
)

type TaskLog struct {
	StartedAt time.Time `json:"started_at"`
	PausedAt  time.Time `json:"paused_at"`
}

type Task struct {
	Item
	Status              TaskStatus `json:"status"`
	Priority            Priority   `json:"priority"`
	DueDate             time.Time  `json:"due_date"`
	CompletedPercentage uint       `json:"completed_percentage"`
	DurationInMinutes   uint       `json:"duration_in_minutes"`
	StartedAt           time.Time  `json:"started_at"`
	CompletedAt         time.Time  `json:"completed_at"`
	CancelledAt         time.Time  `json:"cancelled_at"`
	Logs                []TaskLog  `json:"logs"`
	fsm                 *fsm.FSM
}

func (task *Task) GetStatus() string {
	return string(task.Status)
}

func (task *Task) setLastPausedAt(date time.Time) {
	index := len(task.Logs) - 1
	if index < 0 {
		return
	}
	lastLog := task.Logs[index]
	if lastLog.PausedAt.IsZero() {
		lastLog.PausedAt = date
		task.Logs[index] = lastLog
	}
}

func NewTask(title string, description string, status TaskStatus, tags ...string) *Task {
	task := &Task{
		Item:        *newItem(title, "task", description, tags...),
		Status:      status,
		Priority:    UNKNOWN,
		DueDate:     time.Time{},
		StartedAt:   time.Time{},
		CompletedAt: time.Time{},
		CancelledAt: time.Time{},
	}

	task.fsm = fsm.NewFSM(
		string(status),
		fsm.Events{
			{Name: "start", Src: []string{string(ToDo), string(Stopped), string(Cancelled)}, Dst: string(InProgress)},
			{Name: "stop", Src: []string{string(InProgress)}, Dst: string(Stopped)},
			{Name: "complete", Src: []string{string(InProgress), string(Stopped)}, Dst: string(Completed)},
			{Name: "cancel", Src: []string{string(ToDo), string(Stopped), string(InProgress)}, Dst: string(Cancelled)},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { task.Status = TaskStatus(task.fsm.Current()) },
			"enter_in-progress": func(e *fsm.Event) {
				now := time.Now()
				if task.StartedAt.IsZero() {
					task.StartedAt = now
				}
				task.Logs = append(task.Logs, TaskLog{StartedAt: now})
			},
			"enter_stopped": func(e *fsm.Event) {
				task.setLastPausedAt(time.Now())
			},
			"enter_completed": func(e *fsm.Event) {
				now := time.Now()
				task.setLastPausedAt(now)
				task.CompletedAt = now
			},
			"enter_cancelled": func(e *fsm.Event) {
				now := time.Now()
				task.setLastPausedAt(now)
				task.CancelledAt = now
			},
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

func (task *Task) GetDueDate() *time.Time {
	if task.DueDate.IsZero() {
		return nil
	}
	return &task.DueDate
}
