package types

import (
	"time"

	"github.com/google/uuid"
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

type UserRes struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type LoginUserRes struct {
	SessionID             string  `json:"session_id"`
	AccessToken           string  `json:"access_token"`
	AccessTokenExpiresAt  string  `json:"access_token_expires_at"`
	RefreshToken          string  `json:"refresh_token"`
	RefreshTokenExpiresAt string  `json:"refresh_token_expires_at"`
	User                  UserRes `json:"user"`
}

type CreateInvitePayload struct {
	InvitedUserID uuid.UUID                     `json:"invitedUserId"`
	Role          sqlc.OrganizationRole         `json:"role"`
	Status        sqlc.OrganizationInviteStatus `json:"status"`
	ExpiresAt     time.Time                     `json:"expiresAt,omitempty"`
}

type RenewAccessTokenPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type RenewAccessTokenRes struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
}
