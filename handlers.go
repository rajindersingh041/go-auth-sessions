package main

type contextKey2 string

// // handlers.go
// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// )

// // A custom type for our context key
// type contextKey string
// const userContextKey = contextKey("username")

// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	var creds Credentials
// 	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if creds.Username == "" || creds.Password == "" {
// 		http.Error(w, "Username and password are required", http.StatusBadRequest)
// 		return
// 	}

// 	hashedPassword, err := hashPassword(creds.Password)
// 	if err != nil {
// 		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := createUser(creds.Username, hashedPassword); err != nil {
// 		http.Error(w, "Failed to create user (maybe user exists?)", http.StatusInternalServerError)
// 		return
// 	}

// 	userID, err := findUserID(creds.Username)
// 	if err != nil {
// 		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
// 		return
// 	}
// 	msg := fmt.Sprintf("New user created: %s (ID: %s)\n", creds.Username, string(userID))

// 	w.WriteHeader(http.StatusCreated)

// 	w.Write([]byte(msg))
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	var creds Credentials
// 	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Find the user in ClickHouse
// 	hashedPassword, err := findUserByUsername(creds.Username)
// 	if err != nil {
// 		http.Error(w, "Invalid username", http.StatusUnauthorized)
// 		return
// 	}

// 	// Check the password
// 	if !checkPasswordHash(creds.Password, hashedPassword) {
// 		http.Error(w, "Invalid password", http.StatusUnauthorized)
// 		return
// 	}

// 	// Generate the JWT using the simple method
// 	tokenString, err := generateJWTSimple(creds.Username)
// 	if err != nil {
// 		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send the token back
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"token": tokenString,
// 	})
// }

// func protectedHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get the username from the context (set by the middleware)
// 	username, ok := r.Context().Value(userContextKey).(string)
// 	if !ok {
// 		// This should not happen if middleware is set up correctly
// 		http.Error(w, "Unable to retrieve user from context", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte(fmt.Sprintf("Welcome to the protected area, %s!", username)))
// }

// // authMiddleware validates the JWT and adds user info to the context.
// func authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// 1. Get the "Authorization" header
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Authorization header required", http.StatusUnauthorized)
// 			return
// 		}

// 		// 2. Check for "Bearer " prefix
// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			http.Error(w, "Authorization header must be in format 'Bearer <token>'", http.StatusUnauthorized)
// 			return
// 		}
// 		tokenString := parts[1]

// 		// 3. Validate the token using the simple method
// 		claims, err := validateJWTSimple(tokenString)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
// 			return
// 		}

// 		// 4. Add user info to the request context
// 		ctx := context.WithValue(r.Context(), userContextKey, claims.Username)

// 		// 5. Call the next handler
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }