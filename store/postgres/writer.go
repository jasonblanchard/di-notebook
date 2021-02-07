package postgres

import (
	"database/sql"
	"time"

	"github.com/jasonblanchard/di-notebook/store"
	"github.com/pkg/errors"
)

// Writer postgres write connection
type Writer struct {
	Db *sql.DB
}

// CreateEntry Creates an entry in the DB
func (w *Writer) CreateEntry(i *store.CreateEntryInput) (int, error) {
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}

	var updatedAt sql.NullTime

	if !i.UpdatedAt.IsZero() {
		updatedAt.Time = i.UpdatedAt
		updatedAt.Valid = true
	}

	row := w.Db.QueryRow(`
INSERT INTO entries (text, creator_id, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id
`, i.Text, i.CreatorID, i.CreatedAt, updatedAt)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return id, errors.Wrap(err, "Insert failed")
	}

	return id, nil
}

// UpdateEntry update entry instance
func (w *Writer) UpdateEntry(i *store.UpdateEntryInput) (*store.UpdateEntryOutput, error) {
	now := time.Now()

	row := w.Db.QueryRow(`
UPDATE entries
SET text = $1, updated_at = $2
WHERE id = $3
AND delete_time is null
RETURNING id, text, creator_id, created_at, updated_at
`, i.Text, now, i.ID)

	output := &store.UpdateEntryOutput{}
	err := row.Scan(&output.ID, &output.Text, &output.CreatorID, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "Error updateing entry")
	}

	return output, nil
}

// DeleteEntry delete entry instance
func (w *Writer) DeleteEntry(id int) (*store.DeleteEntryOutput, error) {
	now := time.Now()

	row := w.Db.QueryRow(`
UPDATE entries
SET delete_time = $1
WHERE id = $2
RETURNING id, text, creator_id, created_at, updated_at, delete_time
`, now, id)

	output := &store.DeleteEntryOutput{}
	err := row.Scan(&output.ID, &output.Text, &output.CreatorID, &output.CreatedAt, &output.UpdatedAt, &output.DeleteTime)
	if err != nil {
		return nil, errors.Wrap(err, "Error deleting entry")
	}

	return output, nil
}

// DropEntries Drop all entries in the DB
// Only used internally for tests and such
func (w *Writer) DropEntries() error {
	rows, err := w.Db.Query("DELETE FROM entries")
	if rows != nil {
		defer rows.Close()
	}
	return errors.Wrap(err, "drop failed")
}
