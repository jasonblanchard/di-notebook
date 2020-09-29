package app

import (
	"github.com/pkg/errors"
)

// StartNewEntry start writing a new entry
func (a *App) StartNewEntry(p *Principle, text string, creatorID string) (int, error) {
	// TODO: Check policy to make sure principle can do this
	return a.StoreWriter.CreateEntry(text, creatorID)
}

// ResetEntries drop all entries. Usually used for testing
func (a *App) ResetEntries(p *Principle) error {
	if !canResetEntries(p) {
		return &UnauthorizedError{s: "Principle cannot drop entries"}
	}
	return a.StoreWriter.DropEntries()
}

// ReadEntry get an entry for reading
func (a *App) ReadEntry(p *Principle, id int) (*Entry, error) {
	output, err := a.StoreReader.GetEntry(id)
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
