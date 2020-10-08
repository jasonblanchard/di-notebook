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

// ListEntries lists entries in descending order of created_at
func (r *Reader) ListEntries(i *store.ListEntriesInput) (*store.ListEntriesOutputCollection, error) {
	var rows *sql.Rows
	var err error

	if i.After == 0 {
		rows, err = r.Db.Query(`
SELECT id, text, creator_id, created_at, updated_at
FROM entries
WHERE creator_id = $1
AND is_deleted = false
ORDER BY created_at DESC
LIMIT $2
		`, i.CreatorID, i.First)

		if err != nil {
			return nil, errors.Wrap(err, "Error running query")
		}
		defer rows.Close()
	} else {
		rows, err = r.Db.Query(`
SELECT id, text, creator_id, created_at, updated_at
FROM entries
WHERE creator_id = $1
AND is_deleted = false
AND id < $2
ORDER BY created_at DESC
LIMIT $3
		`, i.CreatorID, i.After, i.First)

		if err != nil {
			return nil, errors.Wrap(err, "Error running query")
		}
		defer rows.Close()
	}

	output := store.ListEntriesOutputCollection{}

	for rows.Next() {
		entry := store.ListEntriesOutput{}
		err := rows.Scan(&entry.ID, &entry.Text, &entry.CreatorID, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "Error scanning rows")
		}
		output = append(output, entry)
	}
	return &output, nil
}
