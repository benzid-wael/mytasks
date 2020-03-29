package cli

import (
	"encoding/json"
	"fmt"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
)

type ItemGroup struct {
	Name  string
	Items entities.ItemCollection
}

type ItemSummarizer func(item entities.Manageable) string
type GetKey func(item entities.Manageable) string

func GroupBy(items entities.ItemCollection, keyFunc GetKey) []ItemGroup {
	itemsByKey := make(map[string]entities.ItemCollection, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		key := keyFunc(item)
		if key == "" {
			continue
		}
		itemsByKey[key] = append(itemsByKey[key], item)
	}

	index := 0
	groupedItems := make([]ItemGroup, len(itemsByKey))
	for key, group := range itemsByKey {
		groupedItems[index] = ItemGroup{
			Name:  key,
			Items: group,
		}
		index++
	}
	return groupedItems
}

func GroupByCreationDate(item entities.Manageable) string {
	return item.GetCreationDateTime().Format("Monday, 02 Jan 2006")
}

func GroupByStatus(item entities.Manageable) string {
	return item.GetStatus()
}

func GroupByTag(item entities.Manageable) string {
	if len(item.GetTags()) > 0 {
		return item.GetTags()[0]
	}
	return "My Board"
}

func FlatByTags(items entities.ItemCollection) entities.ItemCollection {
	newItems := make(entities.ItemCollection, len(items)*5)
	index := 0
	for _, item := range items {
		tags := item.GetTags()
		if len(tags) > 0 {
			for _, tag := range tags {
				var payload map[string]interface{}
				data, _ := json.Marshal(item)
				json.Unmarshal(data, &payload)
				payload["tags"] = []interface{}{interface{}(tag)}
				newItems[index] = entities.CreateItem(payload)
				index++
			}
		} else {
			newItems[index] = item
			index++
		}
	}
	return newItems
}

func (g *ItemGroup) Print(renderer Renderer, summarizer ItemSummarizer) {
	summary := Summarize(g.Items)

	// Print  Header
	color := GetColorStatus(0, summary.TasksCount, summary.DoneCount)
	taskStatus := fmt.Sprintf("[%v/%v]", summary.DoneCount, summary.TasksCount)
	fmt.Printf("\n  %v %v\n", g.Name, renderer.Colorify(taskStatus, color))

	// Print Body
	for _, item := range g.Items {
		fmt.Printf("    %v\n", summarizer(item))
	}
}
