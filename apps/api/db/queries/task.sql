-- name: CreateTask :one
INSERT INTO tasks (board_id, column_id, assignee_id, title, description, position, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE board_id = $1
ORDER BY position
LIMIT $2
OFFSET $3;

-- name: ListTasksByColumn :many
SELECT * FROM tasks
WHERE column_id = $1
ORDER BY position
LIMIT $2
OFFSET $3;

-- name: UpdateTask :one
UPDATE tasks
SET column_id = $2, assignee_id = $3, title = $4, description = $5, position = $6
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;

