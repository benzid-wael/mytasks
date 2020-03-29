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

func GetId(data map[string]interface{}) int {
	return int(data["id"].(float64))
}

func (factory *NoteFactory) Create(data map[string]interface{}) Manageable {
	note := NewNote(data["title"].(string), data["description"].(string), getTags(data)...)
	if id := data["id"]; id != nil {
		note.Id = GetId(data)
	}
	if created_at := data["created_at"].(string); created_at != "" {
		note.CreatedAt, _ = time.Parse(time.RFC3339, created_at)
	}
	if starred := data["is_starred"].(bool); starred == true {
		note.IsStarred = starred
	}
	return note
}

type TaskFactory struct{}

func (factory *TaskFactory) Create(data map[string]interface{}) Manageable {
	task := NewTask(data["title"].(string), data["description"].(string), getTags(data)...)
	if id := data["id"]; id != nil {
		task.Id = GetId(data)
	}
	if created_at := data["created_at"].(string); created_at != "" {
		task.CreatedAt, _ = time.Parse(time.RFC3339, created_at)
	}
	if starred := data["is_starred"].(bool); starred == true {
		task.IsStarred = starred
	}
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
