package middleware

import (
	"context"
	"net/http"
	"oravue_backend/internal/config"
	"oravue_backend/pkg/jwt"
	"strings"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

var jwtSecretKey = []byte("your_secure_secret_key") // Use a secure key loaded from config

func AuthMiddleware(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Ensure it's a Bearer token
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]

			// Call VerifyJWT function
			claims, err := jwt.VerifyJWT(tokenString, config)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract user_id from claims
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "invalid user_id in token claims", http.StatusUnauthorized)
				return
			}

			// Attach user_id to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
