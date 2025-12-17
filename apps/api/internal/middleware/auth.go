package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mirkosisko-dev/api/internal/handlers/auth"
)

func AuthenticationMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)
			if tokenString == "" {
				auth.PermissionDenied(w)
				return
			}

			token, err := auth.ValidateToken(tokenString, jwtSecret)
			if err != nil || token == nil || !token.Valid {
				auth.PermissionDenied(w)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				auth.PermissionDenied(w)
				return
			}

			sub, ok := claims[string(auth.UserKey)].(string)
			if !ok || sub == "" {
				auth.PermissionDenied(w)
				return
			}

			userID, err := uuid.Parse(sub)
			if err != nil {
				auth.PermissionDenied(w)
				return
			}

			ctx := auth.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	reqToken, ok := auth.GetTokenFromRequest(r)
	if !ok {
		return ""
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	return reqToken
}
