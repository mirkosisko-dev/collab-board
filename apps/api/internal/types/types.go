package types

import (
	"time"

	"github.com/mirkosisko-dev/api/db/sqlc"
)

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateInvitePayload struct {
	InvitedUserID int32                         `json:"invitedUserId"`
	Role          sqlc.OrganizationRole         `json:"role"`
	Status        sqlc.OrganizationInviteStatus `json:"status"`
	ExpiresAt     time.Time                     `json:"expiresAt,omitempty"`
}
