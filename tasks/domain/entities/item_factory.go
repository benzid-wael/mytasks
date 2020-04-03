package entities

import (
	"time"
)

type ItemFactory interface {
	Create(data map[string]interface{}) Manageable
}

type NoteFactory struct{}

func getTags(data map[string]interface{}) []string {
	var tags []string = make([]string, len(data["tags"].([]interface{})))
	for idx, val := range data["tags"].([]interface{}) {
		tags[idx] = interface{}(val).(string)
	}

	return tags
}

func getLogs(logs interface{}) []TaskLog {
	if logs == nil {
		return []TaskLog{}
	}
	result := make([]TaskLog, 0, len(logs.([]interface{})))
	for idx, log := range logs.([]interface{}) {
		startedAt, _ := time.Parse(time.RFC3339, log.(map[string]interface{})["started_at"].(string))
		pausedAt, _ := time.Parse(time.RFC3339, log.(map[string]interface{})["paused_at"].(string))
		result[idx] = TaskLog{StartedAt: startedAt, PausedAt: pausedAt}
	}
	return result
}

func GetId(data map[string]interface{}) int {
	return int(data["id"].(float64))
}

func (factory *NoteFactory) Create(data map[string]interface{}) Manageable {
	note := NewNote(data["title"].(string), data["description"].(string), getTags(data)...)
	if id := data["id"]; id != nil {
		note.Id = GetId(data)
	}
	note.CreatedAt = getDate(data["created_at"])
	if starred := data["is_starred"].(bool); starred {
		note.IsStarred = starred
	}
	return note
}

type TaskFactory struct{}

func getDate(value interface{}) time.Time {
	if value == nil || value.(string) == "" {
		return time.Time{}
	}
	date, _ := time.Parse(time.RFC3339, value.(string))
	return date
}

func getFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}
	return value.(float64)
}

func (factory *TaskFactory) Create(data map[string]interface{}) Manageable {
	status := TaskStatus(data["status"].(string))
	task := NewTask(data["title"].(string), data["description"].(string), status, getTags(data)...)
	if id := data["id"]; id != nil {
		task.Id = GetId(data)
	}
	if starred := data["is_starred"].(bool); starred {
		task.IsStarred = starred
	}
	if priority := uint(data["priority"].(float64)); priority > 0 {
		task.Priority = Priority(priority)
	}

	task.CompletedPercentage = uint(getFloat64(data["completed_percentage"]))
	task.DurationInMinutes = uint(getFloat64(data["duration_in_minutes"]))

	task.CreatedAt = getDate(data["created_at"])
	task.StartedAt = getDate(data["started_at"])
	task.CompletedAt = getDate(data["completed_at"])
	task.CancelledAt = getDate(data["cancelled_at"])
	task.DueDate = getDate(data["due_date"])

	task.Logs = getLogs(data["logs"])

	return task
}

func CreateItem(item map[string]interface{}) Manageable {
	kind := item["type"].(string)
	var factory ItemFactory = nil

	switch kind {
	case "note":
		factory = &NoteFactory{}
	case "task":
		factory = &TaskFactory{}
	}

	if factory != nil {
		return factory.Create(item)
	}
	return nil
}
