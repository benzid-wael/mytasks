package entities

import (
	"gotest.tools/assert"
	"testing"
)

func Test_getTags_extract_tags_from_map(t *testing.T) {
	// Given
	expected := []string{"coding", "golang"}
	data := map[string]interface{}{"tags": expected}
	// When
	tags := getTags(data)
	// Then
	assert.DeepEqual(t, expected, tags)
}

func Test_GetId_extract_id_from_map(t *testing.T) {
	// Given
	expected := 22
	data := map[string]interface{}{"id": float64(expected)}
	// When
	id := GetId(data)
	// Then
	assert.Equal(t, expected, id)
}

func Test_NoteFactory_Create_Returns_pointer_to_note(t *testing.T) {
	// Given
	data := map[string]interface{}{
		"id":          float64(2),
		"title":       "My awesome item",
		"description": "",
		"tags":        []string{"coding", "golang"},
	}
	sequence := NewSequence(1)
	factory := new(NoteFactory)
	// When
	actual := factory.Create(sequence, data)
	// Then
	assert.Assert(t, actual != nil)
}

func Test_CreateItem_Creates_Note_from_json_data(t *testing.T) {
	// Given
	data := map[string]interface{}{
		"id":          float64(2),
		"type":        "note",
		"title":       "My awesome item",
		"description": "",
		"tags":        []string{"coding", "golang"},
	}
	// When
	actual := CreateItem(data)
	// Then
	assert.Assert(t, actual != nil)
}
