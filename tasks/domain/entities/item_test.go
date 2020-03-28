package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_NewItem_CreatesInstanceWithNegativeId(t *testing.T) {
	// When
	item := newItem("My item", "", "")
	// Then
	assert.Equal(t, -1, item.Id)
}

func TestItem_NewItem_CreatesUnstarredItem(t *testing.T) {
	// When
	item := newItem("My item", "", "")
	// Then
	assert.False(t, item.IsStarred)
}

func TestItem_NewNote_CreatesUnstarredNote(t *testing.T) {
	// When
	note := NewNote("My note", "", "")
	// Then
	assert.False(t, note.IsStarred)
}

func TestItem_Star_MarkItemAsStarred(t *testing.T) {
	// Given
	item := newItem("My item", "", "")
	// When
	item.Star()
	// Then
	assert.True(t, item.IsStarred)
}

func TestItem_Unstar_MarkItemAsUnStarred(t *testing.T) {
	// Given
	testee := newItem("My item", "", "")
	testee.Star()
	// When
	testee.Unstar()
	// Then
	assert.False(t, testee.IsStarred)
}
