package entities

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/benzid-wael/mytasks/core"
	"github.com/bluele/factory-go/factory"
	"gotest.tools/assert"
	"testing"
)

var DummyItemFactory = factory.NewFactory(
	&Item{},
).SeqInt("Id", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Type", func(args factory.Args) (interface{}, error) {
	return "note", nil
}).Attr("Title", func(args factory.Args) (interface{}, error) {
	return "note", nil
}).Attr("IsStarred", func(args factory.Args) (interface{}, error) {
	return randomdata.Boolean(), nil
}).Attr("Tags", func(args factory.Args) (interface{}, error) {
	len := randomdata.Number(0, 5)
	var tags = make([]string, len)
	for i := 0; i < len; i++ {
		tags[i] = randomdata.SillyName()
	}
	return tags, nil
})

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
	note := DummyItemFactory.MustCreate().(*Item)
	data, _ := core.ToMap(note)
	sequence := NewSequence(1)
	factory := new(NoteFactory)
	// When
	actual := factory.Create(sequence, *data)
	// Then
	assert.Assert(t, actual != nil)
}

func Test_CreateItem_Creates_Note_from_json_data(t *testing.T) {
	// Given
	note := DummyItemFactory.MustCreate().(*Item)
	data, _ := core.ToMap(note)
	// When
	actual := CreateItem(*data)
	// Then
	assert.Assert(t, actual != nil)
}
