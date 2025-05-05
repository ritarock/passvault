package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAction_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		action action
		want   string
	}{
		{
			name:   "Get Generate",
			action: Generate,
			want:   "GENERATE",
		},
		{
			name:   "Get List",
			action: List,
			want:   "LIST",
		},
		{
			name:   "Get Get",
			action: Get,
			want:   "GET",
		},
		{
			name:   "Get Help",
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
