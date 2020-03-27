package cli

import (
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
	"github.com/stretchr/testify/assert"
	path2 "path"
	"testing"
)

func TestGetAppConfig_ReturnsDefaultConfig_WhenConfigFileDoesNotExist(t *testing.T) {
	// Given / When
	defaultDir := tasks.CreateTempDirectory("GetAppConfig")
	actual := GetAppConfig("/tmp/__unknown/mytasks.json", defaultDir)
	// Then
	assert.Equal(t, defaultDir, actual.DataDirectory)
}

func TestGetAppConfig_CreatesConfigFile_WhenItDoesNotExist(t *testing.T) {
	// Given
	defaultDir := tasks.CreateTempDirectory("GetAppConfig")
	expected := path2.Join(defaultDir, "mytasks.json")
	// When
	actual := GetAppConfig(expected, defaultDir)
	// Then
	assert.Equal(t, actual.DataDirectory, defaultDir)
	assert.FileExists(t, expected)
}

func TestGetAppState_ReturnsDefaultAppState_WhenStateFileDoesNotExist(t *testing.T) {
	// Given / When
	defaultDir := tasks.CreateTempDirectory("GetAppState")
	actual := GetAppState(path2.Join(defaultDir, "app.state"))
	// Then
	assert.Equal(t, 0, actual.ItemSequence.Current())
}

func TestGetAppState_DoesNotCreatesStateFile_WhenItDoesNotExist(t *testing.T) {
	// Given
	defaultDir := tasks.CreateTempDirectory("GetAppState")
	expected := path2.Join(defaultDir, "app.state")
	// When
	GetAppState(expected)
	// Then
	assert.NoFileExists(t, expected)
}

func Test_storeAppState_CreatesStateFile(t *testing.T) {
	// Given
	defaultDir := tasks.CreateTempDirectory("storeAppState")
	appStatePath := path2.Join(defaultDir, "app.state")
	appState := AppState{ItemSequence: *value_objects.NewSequence(5)}
	// When
	storeAppState(appStatePath, &appState)
	// Then
	assert.FileExists(t, appStatePath)
}

func TestGetCliApp_ReturnsNonNilObject(t *testing.T) {
	// Given
	defaultDir := tasks.CreateTempDirectory("storeAppState")
	appConfig := AppConfig{
		DataDirectory:         defaultDir,
		DisplayCompletedTasks: false,
		DefaultDisplayMode:    "",
	}
	appState := AppState{ItemSequence: *value_objects.NewSequence(4)}
	// When
	cli := GetCliApp(appConfig, &appState)
	// Then
	assert.NotNil(t, cli)
}
