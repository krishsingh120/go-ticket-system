package middleware

import (
	"context"
	"net/http"
	"strings"
	"ticket-system/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

const UserIDKey ContextKey = "userID"

// JWTMiddleware validates the JWT token and extracts the user ID into the request context.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.JSONError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.JSONError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return utils.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			utils.JSONError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.JSONError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			utils.JSONError(w, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}

		userID := int(userIDFloat)

		// Set the user ID in the request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
