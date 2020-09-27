package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jasonblanchard/di-notebook/store"
	"github.com/pkg/errors"
)

// WriterInput input for creating a store writer
type WriterInput struct {
	User     string
	Password string
	Dbname   string
	Host     string
}

// NewWriter Create a new reader store
func NewWriter(i *WriterInput) (store.Writer, error) {
	r := &Writer{}

	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", i.User, i.Password, i.Host, i.Dbname)
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return r, errors.Wrap(err, "Database connetion failed")
	}

	err = connection.Ping()
	if err != nil {
		return r, errors.Wrap(err, "Database ping failed")
	}
	connection.SetMaxIdleConns(0) // Let Aurora sleep if no connections present.

	r.Db = connection

	return r, nil
}

// Writer postgres write connection
type Writer struct {
	Db *sql.DB
}

// CreateEntry Creates an entry in the DB
func (w *Writer) CreateEntry(text string, creatorID string) (int, error) {
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
