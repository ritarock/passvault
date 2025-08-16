package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_action_string(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		action action
		want   string
	}{
		{
			name:   "get Generate",
			action: Generate,
			want:   "GENERATE",
		},
		{
			name:   "get List",
			action: List,
			want:   "LIST",
		},
		{
			name:   "get Get",
			action: Get,
			want:   "GET",
		},
		{
			name:   "get Help",
			action: Help,
			want:   "HELP",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.action.string()
			assert.Equal(t, test.want, got)
		})
	}
}
