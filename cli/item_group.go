package cli

import (
	"encoding/json"
	"fmt"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"log"
)

type ItemGroup struct {
	Name  string
	Items entities.ItemCollection
}

type ItemGroupCollection []ItemGroup
type ItemSummarizer func(item entities.Manageable) string
type GetKey func(item entities.Manageable) string
type Predicate func(group ItemGroup) bool

const MainBoardName = "My Board"

func GroupBy(items entities.ItemCollection, keyFunc GetKey) ItemGroupCollection {
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
	return MainBoardName
}

func FlatByTags(items entities.ItemCollection) entities.ItemCollection {
	newItems := make(entities.ItemCollection, 0, len(items)*2)
	for _, item := range items {
		tags := item.GetTags()
		if len(tags) > 1 {
			for _, tag := range tags {
				var payload map[string]interface{}
				data, _ := json.Marshal(item)
				err := json.Unmarshal(data, &payload)
				if err != nil {
					log.Fatal("Cannot unmarshal item with ID: ", item.GetId())
				}
				payload["tags"] = []interface{}{interface{}(tag)}
				newItems = append(newItems, entities.CreateItem(payload))
			}
		} else {
			newItems = append(newItems, item)
		}
	}
	return newItems
}

func Filter(summary ItemGroupCollection, predicate Predicate) ItemGroupCollection {
	var result []ItemGroup
	for _, group := range summary {
		if predicate(group) {
			result = append(result, group)
		}
	}
	return result
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

// Sort Protocol
func (c ItemGroupCollection) Len() int           { return len(c) }
func (c ItemGroupCollection) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ItemGroupCollection) Less(i, j int) bool { return c[i].Name < c[j].Name }
