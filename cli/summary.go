package cli

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
)

type Summary struct {
	TasksCount      int
	DoneCount       int
	InProgressCount int
	PendingCount    int
	NoteCount       int
}

func Summarize(items entities.ItemCollection) Summary {
	summary := Summary{}
	for _, item := range items {
		kind := item.GetType()
		switch kind {
		case "note":
			summary.NoteCount++
		case "task":
			summary.TasksCount++
			switch item.GetStatus() {
			case string(entities.InProgress):
				summary.InProgressCount++
			case string(entities.Completed):
				summary.DoneCount++
			case string(entities.ToDo):
				summary.PendingCount++
			case string(entities.Stopped):
				summary.PendingCount++
			}
		}
	}
	return summary
}

func (s *Summary) GetDonePercentage() float32 {
	var percentage float32
	if s.TasksCount != 0 && s.DoneCount > 0 && s.DoneCount < s.TasksCount {
		percentage = float32(s.DoneCount) * 100 / float32(s.TasksCount)
	} else if s.TasksCount == s.DoneCount {
		percentage = 100
	}
	return percentage
}
