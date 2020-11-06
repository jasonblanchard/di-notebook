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

	entry := storeGetEntryOutputToEntry(getEntryOutput)

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

type callback func(*Entry)

// ChangeEntry Change an existing entry
func (a *App) ChangeEntry(i *ChangeEntryInput, callbacks ...callback) (*Entry, error) {
	getEntryOutput, err := a.StoreReader.GetEntry(i.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetEntry failed")
	}

	entry := storeGetEntryOutputToEntry(getEntryOutput)

	if !canChangeEntry(i.Principal, entry) {
		return nil, errors.Wrap(&UnauthorizedError{s: "Principal cannot change entry"}, "Unauthorized")
	}

	updateOutput, err := a.StoreWriter.UpdateEntry(&store.UpdateEntryInput{
		ID:   i.ID,
		Text: i.Text,
	})

	updatedEntry := storeUpdateEntryOutputToEntry(updateOutput)

	for _, f := range callbacks {
		f(updatedEntry)
	}

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

	entry := storeGetEntryOutputToEntry(getEntryOutput)

	if !canDiscardEntry(i.Principal, entry) {
		return errors.Wrap(&UnauthorizedError{s: "Principal cannot read entry"}, "Unauthorized")
	}

	err = a.StoreWriter.DeleteEntry(i.ID)
	if err != nil {
		return errors.Wrap(err, "Delete entry failed")
	}

	return nil
}

// ListEntriesInput input for ListEntries
type ListEntriesInput struct {
	Principal *Principal
	CreatorID string
	First     int
	After     int
}

// ListEntriesPagination pagination info for entries list
type ListEntriesPagination struct {
	TotalCount  int
	HasNextPage bool
	StartCursor int
	EndCursor   int
}

// ListEntriesOutput output for ListEntries
type ListEntriesOutput struct {
	Entries    []Entry
	Pagination ListEntriesPagination
}

// ListEntries lists entries
func (a *App) ListEntries(i *ListEntriesInput) (*ListEntriesOutput, error) {
	if !canListEntries(i.Principal, i.CreatorID) {
		return nil, errors.Wrap(&UnauthorizedError{s: "Principal cannot list entries"}, "Unauthorized")
	}

	input := &store.ListEntriesInput{
		CreatorID: i.CreatorID,
		First:     i.First,
		After:     i.After,
	}

	listEntriesOutput, err := a.StoreReader.ListEntries(input)

	if err != nil {
		return nil, errors.Wrap(err, "Error listing entries")
	}

	entries := listEntryOutputToEntries(listEntriesOutput)

	pagination, err := a.StoreReader.GetEntriesPaginationInfo(&store.GetPaginationInfoInput{
		CreatorID:   i.CreatorID,
		First:       i.First,
		StartCursor: entries[0].ID,
		EndCursor:   entries[len(entries)-1].ID,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error getting pagination info")
	}

	output := &ListEntriesOutput{
		Entries: entries,
		Pagination: ListEntriesPagination{
			TotalCount:  pagination.TotalCount,
			HasNextPage: pagination.HasNextPage,
			StartCursor: pagination.StartCursor,
			EndCursor:   pagination.EndCursor,
		},
	}

	return output, nil
}
