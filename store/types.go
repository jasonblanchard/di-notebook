package store

import (
	"database/sql"
	"time"
)

// CreateEntryInput input for CreateEntry
type CreateEntryInput struct {
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Writer interface for write ops to the store
type Writer interface {
	CreateEntry(*CreateEntryInput) (int, error)
	DeleteEntry(int) error
	DropEntries() error
}

// GetEntryOutput output for GetEntry
type GetEntryOutput struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// Reader interface
type Reader interface {
	GetEntry(id int) (*GetEntryOutput, error)
}
