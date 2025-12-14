package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mirkosisko-dev/api/config"
	"github.com/mirkosisko-dev/api/internal/api/service/auth"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			auth.PermissionDenied(w)
			return
		}

		// Parse + validate
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Envs.JWTSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			auth.PermissionDenied(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			auth.PermissionDenied(w)
			return
		}

		// Extract "sub" (user id)
		sub, ok := claims["sub"]
		if !ok {
			auth.PermissionDenied(w)
			return
		}

		userID, ok := parseSubToInt(sub)
		if !ok {
			auth.PermissionDenied(w)
			return
		}

		// Put user id into context
		ctx := auth.WithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	reqToken, ok := auth.GetTokenFromRequest(r)
	if !ok {
		return ""
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	// Last char is " for some reason
	for len(reqToken) > 0 {
		_, size := utf8.DecodeLastRuneInString(reqToken)
		return reqToken[:len(reqToken)-size]
	}

	return reqToken
}

func parseSubToInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case string:
		n, err := strconv.Atoi(t)
		return n, err == nil
	default:
		return 0, false
	}
}
