package main

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string
const contextKeyUsername contextKey = "username"

// PasswordHasher defines the interface for password hashing and checking
type PasswordHasher interface {
	// Hash hashes the given password and returns the hashed value.
	Hash(password string) (string, error)
	Check(password, hash string) bool
}

// BcryptPasswordHasher is a concrete implementation using bcrypt
type BcryptPasswordHasher struct{}
	// Hash hashes the given password using bcrypt.

func (BcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (BcryptPasswordHasher) Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWTManager defines the interface for JWT operations
type JWTManager interface {
	// Generate creates a new JWT token for the given username.
	Generate(username string) (string, error)
	Validate(token string) (string, error)
}

// SimpleJWTManager is a concrete implementation using jwt_simple.go
type SimpleJWTManager struct{}
	// Generate creates a new JWT token for the given username.

func (SimpleJWTManager) Generate(username string) (string, error) {
	return generateJWTSimple(username)
}

func (SimpleJWTManager) Validate(token string) (string, error) {
// contextKeyUsername is the key used to store the username in request context.
	claims, err := validateJWTSimple(token)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
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
