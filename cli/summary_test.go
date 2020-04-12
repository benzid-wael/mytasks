package cli

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSummarize(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
		entities.NewNote("my note 2", ""),
		entities.NewTask("my task 1", "", entities.ToDo),
		entities.NewTask("my task 2", "", entities.InProgress),
		entities.NewTask("my task 3", "", entities.ToDo),
		entities.NewTask("my task 4", "", entities.Completed),
	}
	// When
	actual := Summarize(items)
	// Then
	assert.Equal(t, actual.TasksCount, 4)
	assert.Equal(t, actual.PendingCount, 2)
	assert.Equal(t, actual.DoneCount, 1)
	assert.Equal(t, actual.NoteCount, 2)
}

func TestSummary_GetDonePercentage(t *testing.T) {
	// Given
	items := entities.ItemCollection{
		entities.NewNote("my note 1", ""),
		entities.NewNote("my note 2", ""),
		entities.NewTask("my task 1", "", entities.ToDo),
		entities.NewTask("my task 2", "", entities.InProgress),
		entities.NewTask("my task 3", "", entities.ToDo),
		entities.NewTask("my task 4", "", entities.Completed),
	}
	testee := Summarize(items)
	// When
	actual := testee.GetDonePercentage()
	// Then
	assert.Equal(t, actual, float32(25))
}
