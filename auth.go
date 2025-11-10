package main

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)


type contextKey string
const contextKeyUsername contextKey = "username"

// hashPassword creates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compares a password with a hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


// generateJWT creates a JWT token for the user using golang-jwt/jwt
func generateJWT(username string) (string, error) {
       return generateJWTSimple(username)
}




// verifyJWT validates and decodes a JWT token using golang-jwt/jwt
func verifyJWT(token string) (*JWTClaims, error) {
       return validateJWTSimple(token)
}



// authMiddleware validates JWT token and extracts user information using golang-jwt/jwt
func (s *Server) authMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	       // Extract token from Authorization header
	       authHeader := r.Header.Get("Authorization")
	       if authHeader == "" {
		       respondError(w, http.StatusUnauthorized, "Missing authorization header")
		       return
	       }

	       // Check if header has Bearer prefix
	       parts := strings.Split(authHeader, " ")
	       if len(parts) != 2 || parts[0] != "Bearer" {
		       respondError(w, http.StatusUnauthorized, "Invalid authorization header format")
		       return
	       }

	       token := parts[1]

	       // Verify token using validateJWTSimple
	       claims, err := verifyJWT(token)
	       if err != nil {
		       respondError(w, http.StatusUnauthorized, "Invalid or expired token")
		       return
	       }

	       // Add username to request context
	       ctx := context.WithValue(r.Context(), contextKeyUsername, claims.Username)
	       next.ServeHTTP(w, r.WithContext(ctx))
       })
}