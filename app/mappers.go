package app

import (
	"github.com/jasonblanchard/di-notebook/store"
)

// StoreGetEntryOutputToEntry mapper
func StoreGetEntryOutputToEntry(o *store.GetEntryOutput) *Entry {
	return &Entry{
		ID:        o.ID,
		Text:      o.Text,
		CreatorID: o.CreatorID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt.Time,
	}
}

// StoreUpdateEntryOutputToEntry mapper
func StoreUpdateEntryOutputToEntry(o *store.UpdateEntryOutput) *Entry {
	return &Entry{
		ID:        o.ID,
		Text:      o.Text,
		CreatorID: o.CreatorID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

// listEntryOutputToEntries mapper
func listEntryOutputToEntries(o *store.ListEntriesOutputCollection) []Entry {
	entries := make([]Entry, len(*o))

	for i, s := range *o {
		entries[i] = Entry{
			ID:        s.ID,
			Text:      s.Text,
			CreatorID: s.CreatorID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt.Time,
		}
	}

	return entries
}
