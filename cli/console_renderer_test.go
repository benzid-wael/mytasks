package cli

import (
	"gotest.tools/assert"
	"log"
	"testing"
)

func TestConsoleRenderer_Colorify_WrapsTextWithCorrectAnsiSequences(t *testing.T) {
	tests := []struct {
		name  string
		color COLOR
		want  string
	}{
		{
			name:  "Convert to Black - 3/4 bits",
			color: BLACK,
			want:  "\033[0;30mHello\033[0m",
		},
		{
			name:  "Convert to Red - 3/4 bits",
			color: RED,
			want:  "\033[0;31mHello\033[0m",
		},
		{
			name:  "Convert to Grey - 24 bits",
			color: GREY,
			want:  "\033[0;38;2;128;128;128mHello\033[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Logger := &ConsoleRenderer{}
			if got := Logger.Colorify("Hello", tt.color); got != tt.want {
				t.Errorf("Colorify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsoleRenderer_Log_PrintsMessageIntoStdout(t *testing.T) {
	// Given
	testee := ConsoleRenderer{}
	// When
	actual := CaptureOutput(func() {
		err := testee.Log(LoggerOptions{Badge: "X", Color: RED}, "Hello")
		if err != nil {
			log.Fatal("Cannot capture log")
		}
	})
	// Then
	assert.Equal(t, "\033[0;31mX\033[0m  Hello\n", actual)
}
