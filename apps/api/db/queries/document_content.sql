-- name: CreateDocumentContent :one
INSERT INTO document_content (document_id, ydoc_state)
VALUES ($1, $2)
RETURNING *;

-- name: GetDocumentContent :one
SELECT * FROM document_content
WHERE document_id = $1;

-- name: UpdateDocumentContent :one
UPDATE document_content
SET ydoc_state = $2, updated_at = NOW()
WHERE document_id = $1
RETURNING *;

-- name: DeleteDocumentContent :exec
DELETE FROM document_content
WHERE document_id = $1;

