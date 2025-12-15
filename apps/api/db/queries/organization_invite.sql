-- name: CreateOrganizationInvite :one
INSERT INTO organization_invite (organization_id, invited_by_user_id, invited_user_id, role, status, expires_at, responded_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetOrganizationInvite :one
SELECT * FROM organization_invite
WHERE id = $1;

-- name: ListOrganizationInvites :many
SELECT * FROM organization_invite 
WHERE invited_user_id = $1;

-- name: UpdateOrganizationInvite :one
UPDATE organization_invite
SET expires_at = $2 
WHERE id = $1
RETURNING *;

-- name: DeleteOrganizationInvite :exec
DELETE FROM organization_invite
WHERE id = $1;

