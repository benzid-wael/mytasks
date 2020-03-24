package entities

type ItemFactory interface {
	Create(sequence *Sequence, data map[string]interface{}) Manageable
}

type NoteFactory struct{}

func getTags(data map[string]interface{}) []string {
	var tags []string = make([]string, len(data["tags"].([]string)))
	for idx, val := range data["tags"].([]string) {
		tags[idx] = val
	}

	return tags
}

func GetId(data map[string]interface{}) int {
	return int(data["id"].(float64))
}

func (factory *NoteFactory) Create(sequence *Sequence, data map[string]interface{}) Manageable {
	return NewNote(sequence, data["title"].(string), data["description"].(string), getTags(data)...)
}

type TaskFactory struct{}

func (factory *TaskFactory) Create(sequence *Sequence, data map[string]interface{}) Manageable {
	return NewTask(sequence, data["title"].(string), data["description"].(string), getTags(data)...)
}

func CreateItem(item map[string]interface{}) Manageable {
	kind := item["type"].(string)
	id := GetId(item)
	sequence := Sequence(id - 1)
	var factory ItemFactory = nil

	switch kind {
	case "note":
		factory = &NoteFactory{}
	case "task":
		factory = &TaskFactory{}
	}

	if factory != nil {
		return factory.Create(&sequence, item)
	}
	return nil
}
