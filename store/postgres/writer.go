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

// DeleteEntry delete entry instance
func (w *Writer) DeleteEntry(id int) error {
	row := w.Db.QueryRow(`
UPDATE entries
SET is_deleted = TRUE
WHERE id = $1
`, id)

	if row.Err() != nil {
		return errors.Wrap(row.Err(), "Delete failed")
	}

	return nil
}

// DropEntries Drop all entries in the DB
func (w *Writer) DropEntries() error {
	rows, err := w.Db.Query("DELETE FROM entries")
	if rows != nil {
		defer rows.Close()
	}
	return errors.Wrap(err, "drop failed")
}
