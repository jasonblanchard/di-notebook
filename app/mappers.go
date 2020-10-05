package app

import "github.com/jasonblanchard/di-notebook/store"

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
