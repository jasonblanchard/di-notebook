package postgres

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Writer postgres write connection
type Writer struct {
	Db *sql.DB
}

// CreateEntry Creates an entry in the DB
func (w *Writer) CreateEntry(text string, creatorID string) (int, error) {
	// TODO: Allow passing in all values
	now := time.Now()

	row := w.Db.QueryRow(`
INSERT INTO entries (text, creator_id, created_at)
VALUES ($1, $2, $3)
RETURNING id
`, text, creatorID, now)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return id, errors.Wrap(err, "Insert failed")
	}

	return id, nil
}

// DropEntries Drop all entries in the DB
func (w *Writer) DropEntries() error {
	rows, err := w.Db.Query("DELETE FROM entries")
	if rows != nil {
		defer rows.Close()
	}
	return errors.Wrap(err, "drop failed")
}
