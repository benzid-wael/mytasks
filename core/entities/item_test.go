package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_NewItem_CreatesInstanceWithNextSequenceValue(t *testing.T) {
	// Given
	var seq Sequence
	expected := seq.Current() + 1
	// When
	item := newItem(&seq, "My item", "")
	// Then
	assert.Equal(t, expected, item.Id)
}

func TestItem_NewItem_CreatesUnstarredItem(t *testing.T) {
	// Given
	var seq Sequence
	// When
	item := newItem(&seq, "My item", "")
	// Then
	assert.False(t, item.IsStarred)
}

func TestItem_NewNote_CreatesUnstarredNote(t *testing.T) {
	// Given
	var seq Sequence
	// When
	note := NewNote(&seq, "My note", "")
	// Then
	assert.False(t, note.IsStarred)
}

func TestItem_Star_MarkItemAsStarred(t *testing.T) {
	// Given
	var seq Sequence
	item := newItem(&seq, "My item", "")
	// When
	item.Star()
	// Then
	assert.True(t, item.IsStarred)
}

func TestItem_Unstar_MarkItemAsUnStarred(t *testing.T) {
	// Given
	var seq Sequence
	testee := newItem(&seq, "My item", "")
	testee.Star()
	// When
	testee.Unstar()
	// Then
	assert.False(t, testee.IsStarred)
}
