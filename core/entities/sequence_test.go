package entities

import (
	"gotest.tools/assert"
	"testing"
)

func TestSequence_Current_ReturnsCurrentIndex(t *testing.T) {
	// Given
	var testee Sequence
	// When
	actual := testee.Current()
	// Then
	assert.Equal(t, 0, actual)
}

func TestSequence_Next_ReturnsIncrementCurrentIndex(t *testing.T) {
	// Given
	current := 2
	testee := Sequence(current)
	// When
	testee.Next()
	// Then
	actual := testee.Current()
	assert.Equal(t, current+1, actual)
}
