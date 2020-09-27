package app

import (
	"time"

	"github.com/pkg/errors"
)

// Principle who is performing the op
type Principle struct {
	Type string
	ID   string
}

// StartNewEntry start writing a new entry
func (a *App) StartNewEntry(p *Principle, text string, creatorID string) (int, error) {
	// TODO: Check policy to make sure principle can do this
	return a.Writer.CreateEntry(text, creatorID)
}

// ResetEntries drop all entries. Usually used for testing
func (a *App) ResetEntries(p *Principle) error {
	// TODO: Check policy to make sure principal can do this
	return a.Writer.DropEntries()
}

// Entry DTO
type Entry struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ReadEntry get an entry for reading
func (a *App) ReadEntry(p *Principle, id int) (*Entry, error) {
	// TODO: Check policy to make sure principal can do this
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

	return entry, nil
}
