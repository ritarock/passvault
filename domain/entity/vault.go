package entity

import (
	"errors"
	"time"
)

const CurrentVaultVersion = "1.0"

var (
	ErrEntryNotFound = errors.New("entry not found")
	ErrEntryExists   = errors.New("entry already exists")
)

type Vault struct {
	Version   string            `json:"version"`
	Entries   map[string]*Entry `json:"entries"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func NewVault() *Vault {
	return &Vault{
		Version:   CurrentVaultVersion,
		Entries:   make(map[string]*Entry),
		UpdatedAt: time.Now(),
	}
}

func (v *Vault) CreateEntry(entry Entry) error {
	if _, exists := v.Entries[entry.ID]; exists {
		return ErrEntryExists
	}
	v.Entries[entry.ID] = &entry
	v.UpdatedAt = time.Now()
	return nil
}

func (v *Vault) GetEntry(id string) (*Entry, error) {
	entry, exists := v.Entries[id]
	if !exists {
		return nil, ErrEntryNotFound
	}
	return entry, nil
}

func (v *Vault) UpdateEntry(entry Entry) error {
	if _, exists := v.Entries[entry.ID]; !exists {
		return ErrEntryNotFound
	}
	v.Entries[entry.ID] = &entry
	v.UpdatedAt = time.Now()
	return nil
}

func (v *Vault) DeleteEntry(id string) error {
	if _, exists := v.Entries[id]; !exists {
		return ErrEntryNotFound
	}
	delete(v.Entries, id)
	v.UpdatedAt = time.Now()
	return nil
}

func (v *Vault) ListEntries() []*Entry {
	entries := make([]*Entry, 0, len(v.Entries))
	for _, entry := range v.Entries {
		entries = append(entries, entry)
	}
	return entries
}
