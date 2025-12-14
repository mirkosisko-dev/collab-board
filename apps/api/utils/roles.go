package utils

import (
	"errors"

	"github.com/mirkosisko-dev/api/db/sqlc"
)

func CanInvite(role sqlc.OrganizationRole) bool {
	return role == sqlc.OrganizationRoleOwner ||
		role == sqlc.OrganizationRoleAdmin
}

func ParseOrganizationRole(s string) (sqlc.OrganizationRole, error) {
	switch s {
	case string(sqlc.OrganizationRoleOwner),
		string(sqlc.OrganizationRoleAdmin),
		string(sqlc.OrganizationRoleMember):
		return sqlc.OrganizationRole(s), nil
	default:
		return "", errors.New("invalid role")
	}
}
