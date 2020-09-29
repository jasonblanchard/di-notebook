package postgres

import (
	"database/sql"

	"github.com/jasonblanchard/di-notebook/store"
	"github.com/pkg/errors"
)

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
