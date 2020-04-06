package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItemCollection_Filter_ReturnsEmptyCollection_WhenFalsePredicateIsPassed(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
	}
	// When
	actual := testee.Filter(func(item Manageable) bool { return false })
	// Then
	assert.Len(t, actual, 0)
}

func TestItemCollection_Filter_ReturnsSameCollection_WhenTruthPredicateIsPassed(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
	}
	// When
	actual := testee.Filter(func(item Manageable) bool { return true })
	// Then
	assert.Equal(t, testee, actual)
}

func TestItemCollection_Exclude_ReturnsEmptyCollection_WhenTruthPredicateIsPassed(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
	}
	// When
	actual := testee.Exclude(func(item Manageable) bool { return true })
	// Then
	assert.Len(t, actual, 0)
}

func TestItemCollection_Exclude_ReturnsSameCollection_WhenFalsePredicateIsPassed(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
	}
	// When
	actual := testee.Exclude(func(item Manageable) bool { return false })
	// Then
	assert.Equal(t, testee, actual)
}

func TestItemCollection_FilterByType(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
		NewTask("my task", "", ToDo, "coding", "python"),
	}
	// When
	actual := testee.FilterByType("note")
	// Then
	assert.Len(t, actual, 2)
}

func TestItemCollection_FilterPending(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
		NewTask("my task", "", ToDo, "coding", "python"),
	}
	// When
	actual := testee.FilterPending()
	// Then
	assert.Len(t, actual, 1)
}

func TestItemCollection_FilterByTag(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
		NewTask("my task", "", ToDo, "coding", "python"),
	}
	// Whem / Then
	assert.Len(t, testee.FilterByTag("@special"), 0)
	assert.Len(t, testee.FilterByTag("python"), 2)
	assert.Len(t, testee.FilterByTag("coding"), 3)
}

func TestItemCollection_FilterByStatus(t *testing.T) {
	// Given
	testee := ItemCollection{
		NewNote("my note 1", "", "coding", "golang"),
		NewNote("my note 2", "", "coding", "python"),
		NewTask("my task 1", "", ToDo, "coding", "python"),
		NewTask("my task 2", "", InProgress, "coding", "python"),
	}
	// Whem / Then
	assert.Len(t, testee.FilterByStatus("cancelled"), 0)
	assert.Len(t, testee.FilterByStatus("todo"), 1)
	assert.Len(t, testee.FilterByStatus("in-progress"), 1)
}
