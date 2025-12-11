-- name: CreateDocument :one
INSERT INTO documents (organization_id, title, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDocument :one
SELECT * FROM documents
WHERE id = $1;

-- name: ListDocuments :many
SELECT * FROM documents
WHERE organization_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateDocument :one
UPDATE documents
SET title = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1;

