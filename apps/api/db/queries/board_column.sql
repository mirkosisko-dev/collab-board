-- name: CreateBoardColumn :one
INSERT INTO board_column (board_id, name, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoardColumn :one
SELECT * FROM board_column
WHERE id = $1;

-- name: ListBoardColumns :many
SELECT * FROM board_column
WHERE board_id = $1
ORDER BY position
LIMIT $2
OFFSET $3;

-- name: UpdateBoardColumn :one
UPDATE board_column
SET name = $2, position = $3
WHERE id = $1
RETURNING *;

-- name: DeleteBoardColumn :exec
DELETE FROM board_column
WHERE id = $1;

