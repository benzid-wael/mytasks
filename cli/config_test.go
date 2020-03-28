package cli

import (
	"github.com/benzid-wael/mytasks/tasks"
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
