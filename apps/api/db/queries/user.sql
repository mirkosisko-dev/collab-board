-- name: CreateUser :one
INSERT INTO users (name, email, password_hash)
VALUES (
  sqlc.arg(name),
  LOWER(sqlc.arg(email)),
  sqlc.arg(password_hash)
)
RETURNING id, name, email, password_hash, created_at;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = LOWER($3)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
