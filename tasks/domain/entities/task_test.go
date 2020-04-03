package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTask_NewTask_CreatesInstanceInTodoSate(t *testing.T) {
	// When
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
	// Then
	assert.Equal(t, ToDo, testee.Status)
}

func TestTask_TriggerEvent_ChangeStatusForValidTransistions(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
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

func TestTask_TriggerEvent_ToMarkItInProgress_UpdatesStartedAtDate(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
	now := time.Now()
	// When
	err := testee.TriggerEvent("start")
	// Then
	assert.Nil(t, err)
	assert.True(t, now.Before(testee.StartedAt))
}

func TestTask_TriggerEvent_ToMarkItInProgress_CreatesNewLogEntry(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
	now := time.Now()
	// When
	err := testee.TriggerEvent("start")
	// Then
	assert.Nil(t, err)
	assert.Len(t, testee.Logs, 1)
	log := testee.Logs[0]
	assert.True(t, log.StartedAt.After(now))
}

func TestTask_TriggerEvent_ToMarkItStopped_UpdatesPausedAt(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
	testee.TriggerEvent("start") // nolint
	now := time.Now()
	// When
	testee.TriggerEvent("stop") // nolint
	// Then
	assert.Len(t, testee.Logs, 1)
	log := testee.Logs[0]
	assert.True(t, log.PausedAt.After(now))
}

func TestTask_TriggerEvent_ToMarkItCompleted_UpdatesCompletedAt(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
	testee.TriggerEvent("start") // nolint
	now := time.Now()
	// When
	testee.TriggerEvent("complete") // nolint
	// Then
	assert.True(t, testee.CompletedAt.After(now))
}

func TestTask_TriggerEvent_ToMarkItCompleted_UpdatesPausedAt(t *testing.T) {
	// Given
	testee := NewTask("Learn Golang", "", ToDo, "@coding", "@goilang")
	testee.TriggerEvent("start") // nolint
	now := time.Now()
	// When
	testee.TriggerEvent("complete") // nolint
	// Then
	assert.Len(t, testee.Logs, 1)
	log := testee.Logs[0]
	assert.True(t, log.PausedAt.After(now))
}
