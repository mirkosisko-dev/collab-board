package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mirkosisko-dev/api/utils"
)

type AuthKey struct{}

func CreateAccessToken(userID uuid.UUID, email string, secret string, durationSeconds int64) (string, *TokenClaims, error) {
	exp := time.Now().Add(time.Second * time.Duration(durationSeconds)).Unix()

	claims, err := NewTokenClaims(userID, email, exp)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	at, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return at, claims, nil
}

func CreateRefreshToken(userID uuid.UUID, email string, secret string, durationSeconds int64) (string, *TokenClaims, error) {
	exp := time.Now().Add(time.Second * time.Duration(durationSeconds)).Unix()

	claims, err := NewTokenClaims(userID, email, exp)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return rt, claims, nil
}

func ValidateToken(t string, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(t, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GetTokenFromRequest(r *http.Request) (string, bool) {
	token := strings.TrimSpace(r.Header.Get("Authorization"))

	if token == "" || token == `""` {
		return "", false
	}

	return token, true
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, AuthKey{}, userID)
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(AuthKey{}).(uuid.UUID)
	return userID, ok
}
