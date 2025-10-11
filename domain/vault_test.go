package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVault_CreateEntry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(*Vault)
		entry  Entry
		hasErr bool
	}{
		{
			name:  "succeed: create new entry",
			setup: func(v *Vault) {},
			entry: Entry{
				ID: "test-id-1",
			},
			hasErr: false,
		},
		{
			name: "failed: duplicate ID with existing entry",
			setup: func(v *Vault) {
				v.Entries["test-id-1"] = &Entry{
					ID: "test-id-1",
				}
			},
			entry: Entry{
				ID: "test-id-1",
			},
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vault := NewVault()
			test.setup(vault)
			err := vault.CreateEntry(test.entry)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVault_GetEntry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(*Vault)
		id     string
		want   *Entry
		hasErr bool
	}{
		{
			name: "succeed: get existing entry",
			setup: func(v *Vault) {
				v.Entries["test-id-1"] = &Entry{
					ID: "test-id-1",
				}
			},
			id:     "test-id-1",
			hasErr: false,
		},
		{
			name:   "failed: get non-existent entry",
			setup:  func(v *Vault) {},
			id:     "test-id-1",
			hasErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vault := NewVault()
			test.setup(vault)

			entry, err := vault.GetEntry(test.id)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, entry)
			}
		})
	}
}

func TestVault_UpdateEntry(t *testing.T) {
	t.Parallel()
	newTitle := "new title"
	tests := []struct {
		name   string
		setup  func(*Vault)
		entry  Entry
		hasErr bool
	}{
		{
			name: "succeed: update existing entry",
			setup: func(v *Vault) {
				v.Entries["test-id-1"] = &Entry{
					ID:    "test-id-1",
					Title: "old title",
				}
			},
			entry: Entry{
				ID:    "test-id-1",
				Title: newTitle,
			},
			hasErr: false,
		},
		{
			name:   "failed: non-existent entry",
			setup:  func(v *Vault) {},
			entry:  Entry{ID: "test-id-1"},
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vault := NewVault()
			test.setup(vault)
			oldUpdatedAt := vault.UpdatedAt
			err := vault.UpdateEntry(test.entry)

			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, oldUpdatedAt, vault.UpdatedAt)
				if entry, exists := vault.Entries[test.entry.ID]; exists {
					assert.Equal(t, entry.Title, newTitle)
				}
			}
		})
	}
}

func TestVault_DeleteEntry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(*Vault)
		id     string
		hasErr bool
	}{
		{
			name: "succeed: delete existing entry",
			setup: func(v *Vault) {
				v.Entries["test-id-1"] = &Entry{
					ID: "test-id-1",
				}
			},
			id:     "test-id-1",
			hasErr: false,
		},
		{
			name:   "failed: non-existent entry",
			setup:  func(v *Vault) {},
			id:     "test-id-1",
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vault := NewVault()
			test.setup(vault)
			oldUpdatedAt := vault.UpdatedAt
			err := vault.DeleteEntry(test.id)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, oldUpdatedAt, vault.UpdatedAt)
			}
		})
	}
}

func TestVault_ListEntries(t *testing.T) {
	t.Parallel()
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	tests := []struct {
		name  string
		setup func(*Vault)
		want  []string
	}{
		{
			name:  "empty vault",
			setup: func(v *Vault) {},
			want:  []string{},
		},
		{
			name: "single entry",
			setup: func(v *Vault) {
				v.Entries["entry-1"] = &Entry{
					ID:           "entry-1",
					LastViewedAt: now,
				}
			},
			want: []string{"entry-1"},
		},
		{
			name: "multiple entries sorted by LastViewedAt",
			setup: func(v *Vault) {
				v.Entries["entry-1"] = &Entry{
					ID:           "entry-1",
					LastViewedAt: twoHoursAgo,
				}
				v.Entries["entry-2"] = &Entry{
					ID:           "entry-2",
					LastViewedAt: oneHourAgo,
				}
				v.Entries["entry-3"] = &Entry{
					ID:           "entry-3",
					LastViewedAt: now,
				}
			},
			want: []string{"entry-3", "entry-2", "entry-2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vault := NewVault()
			test.setup(vault)
			entries := vault.ListEntries()
			assert.Len(t, entries, len(test.want))
		})
	}
}
