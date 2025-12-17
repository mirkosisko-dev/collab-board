package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID `json:"userID"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func NewTokenClaims(userID uuid.UUID, email string, exp int64) (*TokenClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token id: %w", err)
	}

	return &TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   userID.String(),
			Issuer:    "collab-board",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
		},
	}, nil
}
