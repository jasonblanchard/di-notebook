package app

import (
	"github.com/jasonblanchard/di-notebook/store"
	"github.com/pkg/errors"
)

// StartNewEntryInput input for StartNewEntry
type StartNewEntryInput struct {
	Principal *Principal
	Text      string
	CreatorID string
}

// StartNewEntry start writing a new entry
func (a *App) StartNewEntry(i *StartNewEntryInput) (int, error) {
	// TODO: Check policy to make sure principle can do this
	createEntryInput := &store.CreateEntryInput{
		Text:      i.Text,
		CreatorID: i.CreatorID,
	}
	return a.StoreWriter.CreateEntry(createEntryInput)
}

// ReadEntryInput Input for ReadEntry
type ReadEntryInput struct {
	Principal *Principal
	ID        int
}

// ReadEntry get an entry for reading
func (a *App) ReadEntry(i *ReadEntryInput) (*Entry, error) {
	getEntryOutput, err := a.StoreReader.GetEntry(i.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetEntry failed")
	}

	entry := StoreGetEntryOutputToEntry(getEntryOutput)

	if !canReadEntry(i.Principal, entry) {
		return nil, errors.Wrap(&UnauthorizedError{s: "Principal cannot read entry"}, "Unauthorized")
	}

	return entry, nil
}

// ChangeEntryInput input for ChangeEntry
type ChangeEntryInput struct {
	Principal *Principal
	ID        int
	Text      string
}

// ChangeEntry Change an existing entry
func (a *App) ChangeEntry(i *ChangeEntryInput) (*Entry, error) {
	getEntryOutput, err := a.StoreReader.GetEntry(i.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetEntry failed")
	}

	entry := StoreGetEntryOutputToEntry(getEntryOutput)

	if !canChangeEntry(i.Principal, entry) {
		return nil, errors.Wrap(&UnauthorizedError{s: "Principal cannot change entry"}, "Unauthorized")
	}

	updateOutput, err := a.StoreWriter.UpdateEntry(&store.UpdateEntryInput{
		ID:   i.ID,
		Text: i.Text,
	})

	updatedEntry := StoreUpdateEntryOutputToEntry(updateOutput)

	return updatedEntry, nil
}

// DiscardEntryInput Input for ReadEntry
type DiscardEntryInput struct {
	Principal *Principal
	ID        int
}

// DiscardEntry marks entry as deleted
func (a *App) DiscardEntry(i *DiscardEntryInput) error {
	getEntryOutput, err := a.StoreReader.GetEntry(i.ID)
	if err != nil {
		return errors.Wrap(err, "Error getting entry")
	}

	entry := StoreGetEntryOutputToEntry(getEntryOutput)

	if !canDiscardEntry(i.Principal, entry) {
		return errors.Wrap(&UnauthorizedError{s: "Principal cannot read entry"}, "Unauthorized")
	}

	err = a.StoreWriter.DeleteEntry(i.ID)
	if err != nil {
		return errors.Wrap(err, "Delete entry failed")
	}

	return nil
}
