package app

import "github.com/jasonblanchard/di-notebook/store"

// Principle who is performing the op
type Principle struct {
	Type string
	ID   string
}

// StartNewEntry start writing a new entry
func (a *App) StartNewEntry(p *Principle, text string, creatorID string) (int, error) {
	// TODO: Check policy to make sure principle can do the write
	return a.Writer.CreateEntry(text, creatorID)
}

// ResetEntries drop all entries. Usually used for testing
func (a *App) ResetEntries(p *Principle) error {
	// TODO: Check policy to make sure principal can do this
	return a.Writer.DropEntries()
}

// ReadEntry get an entry for reading
func (a *App) ReadEntry(p *Principle, id int) (*store.GetEntryOutput, error) { // TODO: Change output type
	// TODO: Check policy to make sure principal can do this
	return a.Reader.GetEntry(id)
}
