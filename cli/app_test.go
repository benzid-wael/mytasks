package cli

import (
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCliApp_ReturnsNonNilObject(t *testing.T) {
	// Given
	defaultDir := tasks.CreateTempDirectory("storeAppState")
	appConfig := AppConfig{
		DataDirectory:         defaultDir,
		DisplayCompletedTasks: false,
		DefaultDisplayMode:    "",
	}
	// When
	cli := GetCliApp(appConfig)
	// Then
	assert.NotNil(t, cli)
}
