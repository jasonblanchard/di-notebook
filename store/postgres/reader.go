package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jasonblanchard/di-notebook/store"
	"github.com/pkg/errors"
)

// ReaderInput input for creating a store reader
type ReaderInput struct {
	User     string
	Password string
	Dbname   string
	Host     string
}

// NewReader Create a new reader store
func NewReader(i *ReaderInput) (store.Reader, error) {
	r := &Reader{}

	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", i.User, i.Password, i.Host, i.Dbname)
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return r, errors.Wrap(err, "Database connetion failed")
	}

	err = connection.Ping()
	if err != nil {
		return r, errors.Wrap(err, "ping failed")
	}
	connection.SetMaxIdleConns(0) // Let Aurora sleep if no connections present.

	r.Db = connection

	return r, nil
}

// Reader postgres reader
type Reader struct {
	Db *sql.DB
}

// GetEntry gets an entry
func (r *Reader) GetEntry(id int) (*store.GetEntryOutput, error) {
	row := r.Db.QueryRow(`
SELECT id, text, creator_id, created_at, updated_at
FROM entries
WHERE id = $1
AND is_deleted = false
`, id)

	output := &store.GetEntryOutput{}
	err := row.Scan(&output.ID, &output.Text, &output.CreatorID, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}

	return output, nil
}
