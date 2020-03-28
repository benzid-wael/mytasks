package entities

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
	return NewNote(data["title"].(string), data["description"].(string), getTags(data)...)
}

type TaskFactory struct{}

func (factory *TaskFactory) Create(data map[string]interface{}) Manageable {
	return NewTask(data["title"].(string), data["description"].(string), getTags(data)...)
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
