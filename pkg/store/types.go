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

// UpdateEntryInput input for UpdateEntry
type UpdateEntryInput struct {
	ID   int
	Text string
}

// UpdateEntryOutput output for UpdateEntry
type UpdateEntryOutput struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DeleteEntryOutput output for UpdateEntry
type DeleteEntryOutput struct {
	ID         int
	Text       string
	CreatorID  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeleteTime time.Time
}

// Writer interface for write ops to the store
// TODO: Add reader methods to writer and initialize
type Writer interface {
	CreateEntry(*CreateEntryInput) (int, error)
	UpdateEntry(*UpdateEntryInput) (*UpdateEntryOutput, error)
	DeleteEntry(int) (*DeleteEntryOutput, error)
	DropEntries() error
	UndeleteEntry(int) error
}

// GetEntryOutput output for GetEntry
type GetEntryOutput struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// ListEntriesInput input for ListEntries
type ListEntriesInput struct {
	CreatorID string
	First     int
	After     int
}

// ListEntriesOutput singular output for ListEntries
type ListEntriesOutput struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// ListEntriesOutputCollection composit output
type ListEntriesOutputCollection []ListEntriesOutput

// GetPaginationInfoInput singular output for GetPaginationInfo
type GetPaginationInfoInput struct {
	CreatorID   string
	First       int
	StartCursor int
	EndCursor   int
}

// GetEntriesPaginationInfoOutput output for GetEntriesPaginationInfo
type GetEntriesPaginationInfoOutput struct {
	TotalCount  int
	HasNextPage bool
	StartCursor int
	EndCursor   int
}

// Reader interface
type Reader interface {
	GetEntry(id int) (*GetEntryOutput, error)
	GetDeletedEntry(id int) (*GetEntryOutput, error)
	ListEntries(*ListEntriesInput) (*ListEntriesOutputCollection, error)
	GetEntriesPaginationInfo(*GetPaginationInfoInput) (*GetEntriesPaginationInfoOutput, error)
}
