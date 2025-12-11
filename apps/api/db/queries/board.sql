-- name: CreateBoard :one
INSERT INTO boards (organization_id, name, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoard :one
SELECT * FROM boards
WHERE id = $1;

-- name: ListBoards :many
SELECT * FROM boards
WHERE organization_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateBoard :one
UPDATE boards
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteBoard :exec
DELETE FROM boards
WHERE id = $1;

