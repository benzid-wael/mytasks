package entities

import (
	"gotest.tools/assert"
	"testing"
)

func Test_getTags_extract_tags_from_map(t *testing.T) {
	// Given
	raw := []interface{}{"coding", "golang"}
	data := map[string]interface{}{"tags": raw}
	// When
	tags := getTags(data)
	// Then
	assert.DeepEqual(t, []string{"coding", "golang"}, tags)
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
	item := GenerateDummyRawItem()
	factory := new(NoteFactory)
	// When
	actual := factory.Create(item)
	// Then
	assert.Assert(t, actual != nil)
}

func Test_CreateItem_Creates_Note_from_json_data(t *testing.T) {
	// Given
	item := GenerateDummyRawItem()
	// When
	actual := CreateItem(item)
	// Then
	assert.Assert(t, actual != nil)
}
