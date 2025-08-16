package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handler_validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		handler  handler
		hasError bool
	}{
		{
			name:     "failed: not enough args",
			handler:  handler{args: []string{}},
			hasError: true,
		},
		{
			name:     "failed: invalid args",
			handler:  handler{args: []string{"invalid"}},
			hasError: true,
		},
		{
			name:     "pass: generate",
			handler:  handler{args: []string{"generate", "title", "url"}},
			hasError: false,
		},
		{
			name:     "pass: generate help",
			handler:  handler{args: []string{"generate", "help"}},
			hasError: false,
		},
		{
			name:     "failed: generate, many args",
			handler:  handler{args: []string{"generate", "title", "url", "unnecessary"}},
			hasError: true,
		},
		{
			name:     "failed: generate, fewer args",
			handler:  handler{args: []string{"generate"}},
			hasError: true,
		},
		{
			name:     "failed: generate, fewer args",
			handler:  handler{args: []string{"generate", "title"}},
			hasError: true,
		},
		{
			name:     "pass: list",
			handler:  handler{args: []string{"list"}},
			hasError: false,
		},
		{
			name:     "pass: list help",
			handler:  handler{args: []string{"list", "help"}},
			hasError: false,
		},
		{
			name:     "failed: list, many args",
			handler:  handler{args: []string{"list", "unnecessary"}},
			hasError: true,
		},
		{
			name:     "pass: get",
			handler:  handler{args: []string{"get", "1"}},
			hasError: false,
		},
		{
			name:     "pass: get help",
			handler:  handler{args: []string{"get", "help"}},
			hasError: false,
		},
		{
			name:     "failed: get, many args",
			handler:  handler{args: []string{"get", "1", "2"}},
			hasError: true,
		},
		{
			name:     "failed: get, fewer args",
			handler:  handler{args: []string{"get"}},
			hasError: true,
		},
		{
			name:     "failed: get, non-numeric",
			handler:  handler{args: []string{"get", "abc"}},
			hasError: true,
		},
		{
			name:     "failed: get, id = 0",
			handler:  handler{args: []string{"get", "0"}},
			hasError: true,
		},
		{
			name:     "pass: help",
			handler:  handler{args: []string{"help"}},
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.handler.validate()
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_handler_mapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler handler
		want    subCommand
	}{
		{
			name:    "generate",
			handler: handler{args: []string{"generate", "title", "url"}},
			want: subCommand{
				action: Generate,
				generateData: generateData{
					title: "title",
					url:   "url",
				},
			},
		},
		{
			name:    "generate with help",
			handler: handler{args: []string{"generate", "help"}},
			want: subCommand{
				action: Generate,
				help:   true,
			},
		},
		{
			name:    "list",
			handler: handler{args: []string{"list"}},
			want: subCommand{
				action: List,
			},
		},
		{
			name:    "list with help",
			handler: handler{args: []string{"list", "help"}},
			want: subCommand{
				action: List,
				help:   true,
			},
		},
		{
			name:    "get",
			handler: handler{args: []string{"get", "1"}},
			want: subCommand{
				action:  Get,
				getData: 1,
			},
		},
		{
			name:    "get with help",
			handler: handler{args: []string{"get", "help"}},
			want: subCommand{
				action: Get,
				help:   true,
			},
		},
		{
			name:    "help",
			handler: handler{args: []string{"help"}},
			want: subCommand{
				action: Help,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			test.handler.mapper()
			assert.Equal(t, test.want, test.handler.subCommand)
		})
	}
}

func Test_isHelp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "return true",
			arg:  "help",
			want: true,
		},
		{
			name: "return false",
			arg:  "test",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := isHelp(test.arg)
			assert.Equal(t, test.want, got)
		})
	}
}
