package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCliApp_ReturnsNonNilObject(t *testing.T) {
	// When
	cli := GetCliApp()
	// Then
	assert.NotNil(t, cli)
}
