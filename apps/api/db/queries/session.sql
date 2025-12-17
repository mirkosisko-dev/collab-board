-- name: CreateSesion :one
INSERT INTO sessions (user_id, refresh_token, is_revoked, expires_at)
VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1;

-- name: ListSessions :many
SELECT * FROM sessions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateSession :one
UPDATE sessions
SET is_revoked = $2, expires_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = true
WHERE id = $1;
