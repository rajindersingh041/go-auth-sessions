// jwt_simple.go - Much simpler JWT implementation using golang-jwt/jwt
package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key - in production, get this from environment variables
var jwtSecretKey = []byte("motu-munni-dobo")

// Claims structure for JWT payload
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// generateJWTSimple creates a signed JWT token - much simpler!
func generateJWTSimple(username string) (string, error) {
	// Create claims with expiration time
	claims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWTSimple parses and validates a JWT token - much simpler!
func validateJWTSimple(tokenString string) (*JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid and get claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}