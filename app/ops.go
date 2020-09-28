package app

import (
	"time"

	"github.com/pkg/errors"
)

// Principle - entity that is performing the op
type Principle struct {
	Type string // TODO: Some kind of enum for type?
	ID   string
}

// Entry - text entry
type Entry struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// StartNewEntry start writing a new entry
func (a *App) StartNewEntry(p *Principle, text string, creatorID string) (int, error) {
	// TODO: Check policy to make sure principle can do this
	return a.Writer.CreateEntry(text, creatorID)
}

// ResetEntries drop all entries. Usually used for testing
func (a *App) ResetEntries(p *Principle) error {
	if !canResetEntries(p) {
		return &UnauthorizedError{s: "Principle cannot drop entries"}
	}
	return a.Writer.DropEntries()
}

// ReadEntry get an entry for reading
func (a *App) ReadEntry(p *Principle, id int) (*Entry, error) {
	output, err := a.Reader.GetEntry(id)
	if err != nil {
		return nil, errors.Wrap(err, "GetEntry failed")
	}

	entry := &Entry{
		ID:        output.ID,
		Text:      output.Text,
		CreatorID: output.CreatorID,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt.Time,
	}

	if !canReadEntry(p, entry) {
		return nil, &UnauthorizedError{s: "Principle cannot read entry"}
	}

	return entry, nil
}
