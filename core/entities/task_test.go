package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTask_NewTask_CreatesInstanceInTodoSate(t *testing.T) {
	// Given
	var seq Sequence
	// When
	testee := NewTask(&seq, "Learn Golang", "@coding", "@goilang")
	// Then
	assert.Equal(t, testee.Status, ToDo)
}
