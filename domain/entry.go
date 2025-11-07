package domain

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	URL          string    `json:"url"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastViewedAt time.Time `json:"last_viewed_at"`
}

func NewEntry(title, username, password, url, notes string) *Entry {
	now := time.Now()
	return &Entry{
		ID:           uuid.New().String(),
		Title:        title,
		Username:     username,
		Password:     password,
		URL:          url,
		Notes:        notes,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastViewedAt: now,
	}
}

func (e *Entry) Update(title, username, password, url, notes string) {
	e.Title = title
	e.Username = username
	e.Password = password
	e.URL = url
	e.Notes = notes
	e.UpdatedAt = time.Now()
}

func (e *Entry) MarkAsViewed() {
	e.LastViewedAt = time.Now()
}
