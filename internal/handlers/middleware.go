package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/database"
	"github.com/quentinsteinke/mkvmender/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware validates API key and adds user to context
func AuthMiddleware(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get API key from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Expected format: "Bearer <api_key>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			apiKey := parts[1]

			// Validate API key
			user, err := db.GetUserByAPIKey(apiKey)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "invalid API key")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated user from context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware checks if user has admin or moderator role
// Must be used after AuthMiddleware
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			respondError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		// Check if user has admin or moderator role
		if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
			respondError(w, http.StatusForbidden, "insufficient permissions")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add proper logging
		next.ServeHTTP(w, r)
	})
}
