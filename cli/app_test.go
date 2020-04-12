package cli

import (
	"errors"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
	"testing"
)

func TestGetCliApp_ReturnsNonNilObject(t *testing.T) {
	// When
	cli := GetCliApp()
	// Then
	assert.NotNil(t, cli)
}

func TestGetItemUseCase_InitializesGivenDirectory(t *testing.T) {
	// Given
	dir := tasks.CreateTempDirectory("GetItemUseCase")
	// When
	actual := GetItemUseCase("", AppConfig{DataDirectory: dir})
	// Then
	assert.Len(t, actual.GetItems(), 0)
}

func TestBulkFunc_PrintErrorMessageWhenCallbackReturnsError(t *testing.T) {
	// Given
	renderer := NewRenderer()
	dir := tasks.CreateTempDirectory("BulkFunc")
	testee := GetItemUseCase("", AppConfig{DataDirectory: dir})
	testee.CreateNote("my note 1")     // nolint
	testee.CreateNote("my note 2")     // nolint
	testee.CreateTask("My first task") // nolint
	// When
	actual := CaptureOutput(func() {
		BulkFunc([]int{1, 2}, "starred", renderer, nil, func(id int, c *cli.Context) error { // nolint
			return errors.New("My Error " + strconv.Itoa(id))
		})
	})
	// Then
	assert.True(t, strings.Contains(actual, "My Error 1"))
	assert.True(t, strings.Contains(actual, "My Error 2"))
}

func TestBulkFunc_(t *testing.T) {
	// Given
	renderer := NewRenderer()
	dir := tasks.CreateTempDirectory("BulkFunc")
	testee := GetItemUseCase("", AppConfig{DataDirectory: dir})
	testee.CreateNote("my note 1")     // nolint
	testee.CreateNote("my note 2")     // nolint
	testee.CreateTask("My first task") // nolint
	// When
	actual := CaptureOutput(func() {
		BulkFunc([]int{1, 2}, "starred", renderer, nil, func(id int, c *cli.Context) error { // nolint
			if id == 2 {
				return errors.New("My Error 2")
			}
			return nil
		})
	})
	// Then
	assert.True(t, strings.Contains(actual, "My Error 2"))
	assert.True(t, strings.Contains(actual, "starred"))
}
