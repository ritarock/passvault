package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Update(t *testing.T) {
	t.Parallel()
	entry := NewEntry("old title", "old username", "old password", "old url", "old notes")
	tests := []struct {
		name     string
		title    string
		username string
		password string
		url      string
		notes    string
	}{
		{
			name:     "succeed",
			title:    "new title",
			username: "new username",
			password: "new password",
			url:      "new url",
			notes:    "new notes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry.Update(test.title, test.username, test.password, test.url, test.notes)
			assert.Equal(t, entry.Title, test.title)
			assert.Equal(t, entry.Username, test.username)
			assert.Equal(t, entry.Password, test.password)
			assert.Equal(t, entry.URL, test.url)
			assert.Equal(t, entry.Notes, test.notes)
		})
	}
}
func TestEntry_MarkAsViewed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		entry Entry
	}{
		{
			name:  "succeed",
			entry: *NewEntry("title", "username", "password", "url", "notes"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldLastViewdAt := test.entry.LastViewedAt
			test.entry.MarkAsViewed()
			newLastViewedAt := test.entry.LastViewedAt
			assert.NotEqual(t, newLastViewedAt, oldLastViewdAt)
		})
	}
}
