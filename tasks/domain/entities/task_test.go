package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTask_NewTask_CreatesInstanceInTodoSate(t *testing.T) {
	// When
	testee := NewTask("Learn Golang", "", "@coding", "@goilang")
	// Then
	assert.Equal(t, ToDo, testee.Status)
}

func TestTask_TriggerEvent_ChangeStatusForValidTransistions(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", "@coding", "@goilang")
	// When
	err := testee.TriggerEvent("start")
	// Then
	assert.Equal(t, InProgress, testee.Status)
	assert.Nil(t, err)
}

func TestTask_TriggerEvent_ReturnsErrorForInvalidTransitions(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", "@coding", "@goilang")
	// When
	actual := testee.TriggerEvent("stop")
	// Then
	assert.Error(t, actual)
}
