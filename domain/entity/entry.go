package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	URL          string    `json:"url"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastViewedAt time.Time `json:"last_viewed_at"`
}

func NewEntry(name, password, url, notes string) *Entry {
	now := time.Now()
	return &Entry{
		ID:           uuid.New().String(),
		Name:         name,
		Password:     password,
		URL:          url,
		Notes:        notes,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastViewedAt: time.Time{},
	}
}

func (e *Entry) Update(name, password, url, notes string) {
	e.Name = name
	e.Password = password
	e.URL = url
	e.Notes = notes
	e.UpdatedAt = time.Now()
}

func (e *Entry) MarkAsViewed() {
	e.LastViewedAt = time.Now()
}
