-- name: GetEntry :one
SELECT id, text, creator_id, created_at, updated_at
FROM entries
WHERE id = $1
AND delete_time IS NULL;

-- name: GetDeletedEntry :one
SELECT id, text, creator_id, created_at, updated_at
FROM entries
WHERE id = $1
AND delete_time IS NOT NULL;

-- name: CreateEntry :one
INSERT INTO entries (text, creator_id, created_at)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateEntry :one
UPDATE entries
SET text = $1, updated_at = $2
WHERE id = $3
AND delete_time is null
RETURNING id, text, creator_id, created_at, updated_at;

-- name: DeleteEntry :one
UPDATE entries
SET delete_time = $1
WHERE id = $2
RETURNING id, text, creator_id, created_at, updated_at, delete_time;

-- name: UndeleteEntry :one
UPDATE entries
SET delete_time = $1
WHERE id = $2
RETURNING id;