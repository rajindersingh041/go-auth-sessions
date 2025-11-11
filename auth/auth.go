package auth // Package auth provides authentication-related functionalities including password hashing and JWT token management.

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey is a type for context keys used in authentication
// UsernameContextKey is the context key for the username in requests
// Exported for use in other modules

type ContextKey string

const UsernameContextKey ContextKey = "username"

// WithJWTAuth is a generic HTTP middleware for JWT authentication.
//
// It validates the Authorization header for a Bearer token, verifies the JWT using the provided JWTManager,
// and injects the username into the request context using UsernameContextKey. If authentication fails,
// it responds with HTTP 401 Unauthorized and does not call the next handler.
//
// Usage:
//
//   mux.Handle("GET /protected", auth.WithJWTAuth(jwtManager, http.HandlerFunc(protectedHandler)))
//
// In your handler, retrieve the username from context:
//
//   username, ok := r.Context().Value(auth.UsernameContextKey).(string)
//   if !ok || username == "" {
//       // handle unauthenticated
//   }
//
// This middleware can be reused in any module that requires JWT authentication.
func WithJWTAuth(jwtManager JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required. Please provide JWT token.", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header must start with 'Bearer '. Format: 'Bearer <token>'", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "JWT token is required. Please login first.", http.StatusUnauthorized)
			return
		}
		username, err := jwtManager.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired JWT token. Please login again.", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, UsernameContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
