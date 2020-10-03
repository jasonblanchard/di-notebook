package app

import (
	"github.com/pkg/errors"
)

// StartNewEntryInput input for StartNewEntry
type StartNewEntryInput struct {
	Principle *Principle
	Text      string
	CreatorID string
}

// StartNewEntry start writing a new entry
func (a *App) StartNewEntry(i *StartNewEntryInput) (int, error) {
	// TODO: Check policy to make sure principle can do this
	return a.StoreWriter.CreateEntry(i.Text, i.CreatorID)
}

// ResetEntries drop all entries. Usually used for testing
func (a *App) ResetEntries(p *Principle) error {
	if !canResetEntries(p) {
		return &UnauthorizedError{s: "Principle cannot drop entries"}
	}
	return a.StoreWriter.DropEntries()
}

// ReadEntryInput Input for ReadEntry
type ReadEntryInput struct {
	Principle *Principle
	ID        int
}

// ReadEntry get an entry for reading
func (a *App) ReadEntry(i *ReadEntryInput) (*Entry, error) {
	output, err := a.StoreReader.GetEntry(i.ID)
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

	if !canReadEntry(i.Principle, entry) {
		return nil, &UnauthorizedError{s: "Principle cannot read entry"}
	}

	return entry, nil
}
