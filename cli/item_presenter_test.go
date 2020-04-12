package cli

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetColorStatus(t *testing.T) {
	type args struct {
		min int
		max int
		i   int
	}
	tests := []struct {
		name string
		args args
		want COLOR
	}{
		{
			name: "completed",
			args: args{min: 0, max: 100, i: 100},
			want: GREEN,
		}, {
			name: "in-progress",
			args: args{min: 0, max: 100, i: 50},
			want: YELLOW,
		}, {
			name: "min == max",
			args: args{min: 0, max: 0, i: 50},
			want: YELLOW,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetColorStatus(tt.args.min, tt.args.max, tt.args.i); got != tt.want {
				t.Errorf("GetColorStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemPresenter_PrintSummary(t *testing.T) {
	// Given
	presenter := itemPresenter{renderer: NewRenderer()}
	summary := Summary{
		TasksCount:      3,
		DoneCount:       2,
		InProgressCount: 0,
		PendingCount:    1,
		NoteCount:       3,
	}
	// When
	actual := CaptureOutput(func() { // nolint
		presenter.PrintSummary(summary, "") // nolint
	})
	// Then
	//assert.Equal(t, actual, "66.67% of all tasks complete.\n2  done  0  in-progress  1  pending  3  notes")
	assert.True(t, strings.Contains(actual, "66.67%"))
}
