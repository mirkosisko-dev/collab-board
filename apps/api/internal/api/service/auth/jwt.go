package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mirkosisko-dev/api/config"
	"github.com/mirkosisko-dev/api/utils"
)

type contextKey string

const userKey contextKey = "sub"

func CreateJWT(userID uuid.UUID) (string, error) {
	exp := time.Now().Add(time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID.String(),
			"iss": "collab-board",
			"exp": exp,
			"iat": time.Now().Unix(),
		})

	jwt, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func CreateRefreshToken(userID int) (string, error) {
	exp := time.Now().Add(time.Second * time.Duration(config.Envs.RefreshTokenExpirationInSeconds)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"iss": "collab-board",
			"exp": exp,
			"iat": time.Now().Unix(),
		})

	rt, err := token.SignedString([]byte(config.Envs.RefreshTokenSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTokenFromRequest(r *http.Request) (string, bool) {
	token := strings.TrimSpace(r.Header.Get("Authorization"))

	if token == "" || token == `""` {
		return "", false
	}

	return token, true
}

func ValidateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userKey).(uuid.UUID)
	return userID, ok
}
