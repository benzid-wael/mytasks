package cli

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupBy_CallsGetKeyFuncForEachItem(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
		entities.NewNote("my note 2", ""),
	}
	// When
	groupedByTitle := GroupBy(items, func(item entities.Manageable) string {
		return item.GetTitle()
	})
	groupedByType := GroupBy(items, func(item entities.Manageable) string {
		return item.GetType()
	})
	// Then
	assert.Len(t, groupedByTitle, 2)
	assert.Len(t, groupedByType, 1)
}

func TestGroupByCreationDate(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
		entities.NewNote("my note 2", ""),
	}
	// When
	actual := GroupBy(items, GroupByCreationDate)
	// Then
	assert.Len(t, actual, 1)
}

func TestGroupByStatus_RemovesNotes(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
		entities.NewNote("my note 2", ""),
		entities.NewTask("my task 1", "", entities.ToDo),
	}
	// When
	actual := GroupBy(items, GroupByStatus)
	// Then
	assert.Len(t, actual, 1)
	assert.Len(t, actual[0].Items, 1)
}

func TestGroupByTag(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
		entities.NewNote("my note 2", ""),
		entities.NewTask("my task 1", "", entities.ToDo),
	}
	// When
	actual := GroupBy(items, GroupByTag)
	// Then
	assert.Len(t, actual, 1)
	assert.Equal(t, actual[0].Name, MainBoardName)
}

func TestFlatByTags_ReturnsArrayOfItemsWithOneTagEach_WhenItemHasTags(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", "", "tag1", "tag2"),
	}
	// When
	actual := FlatByTags(items)
	// Then
	assert.Len(t, actual, 2)
	assert.Len(t, actual[0].GetTags(), 1)
	assert.Len(t, actual[1].GetTags(), 1)
}

func TestFlatByTags_ReturnsArrayWithOneItem_WhenItemDoesNotHaveTags(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
	}
	// When
	actual := FlatByTags(items)
	// Then
	assert.Len(t, actual, 1)
	assert.Len(t, actual[0].GetTags(), 0)
}
