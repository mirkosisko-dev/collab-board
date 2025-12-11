-- name: CreateMessage :one
INSERT INTO messages (board_id, organization_id, user_id, document_id, content)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1;

-- name: ListMessages :many
SELECT * FROM messages
WHERE board_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListMessagesByOrganization :many
SELECT * FROM messages
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdateMessage :one
UPDATE messages
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;

