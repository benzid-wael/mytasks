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
	assert.Equal(t, actual, 0)
}

func TestSequence_Next_ReturnsIncrementCurrentIndex(t *testing.T) {
	// Given
	current := 2
	var testee Sequence = Sequence(current)
	// When
	testee.Next()
	// Then
	actual := testee.Current()
	assert.Equal(t, actual, current+1)
}
