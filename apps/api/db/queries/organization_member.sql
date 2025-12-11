-- name: CreateOrganizationMember :one
INSERT INTO organization_members (organization_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrganizationMember :one
SELECT * FROM organization_members
WHERE id = $1;

-- name: GetOrganizationMemberByOrgAndUser :one
SELECT * FROM organization_members
WHERE organization_id = $1 AND user_id = $2;

-- name: ListOrganizationMembers :many
SELECT * FROM organization_members
WHERE organization_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateOrganizationMember :one
UPDATE organization_members
SET role = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrganizationMember :exec
DELETE FROM organization_members
WHERE id = $1;

