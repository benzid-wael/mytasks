package entities

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/benzid-wael/mytasks/core"
	"github.com/bluele/factory-go/factory"
	"strconv"
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

func GenerateDummyRawItem() map[string]interface{} {
	note := DummyItemFactory.MustCreate().(*Item)
	item, _ := core.ToMap(note)
	return *item
}

func GenerateDummyRawItems(size int) map[string]map[string]interface{} {
	var items map[string]map[string]interface{} = make(map[string]map[string]interface{})
	for i := 0; i < size; i++ {
		note := DummyItemFactory.MustCreate().(*Item)
		item, _ := core.ToMap(note)
		items[strconv.Itoa(i)] = *item
	}
	return items
}
