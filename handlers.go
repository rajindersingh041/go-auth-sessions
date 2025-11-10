package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// handleHealth is a simple health check endpoint
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	}
}

// handleRegister processes user registration
func (s *Server) handleRegister() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate input
		if req.Username == "" || req.Password == "" {
			respondError(w, http.StatusBadRequest, "Username and password are required")
			return
		}

		if len(req.Password) < 8 {
			respondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		// Hash password
		passwordHash, err := hashPassword(req.Password)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to process password")
			return
		}

		// Create user with timeout context
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := s.userRepo.Create(ctx, req.Username, passwordHash); err != nil {
			log.Printf("Failed to create user: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		respondJSON(w, http.StatusCreated, map[string]string{
			"message": "User registered successfully",
		})
	}
}

// handleLogin processes user authentication
func (s *Server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate input
		if req.Username == "" || req.Password == "" {
			respondError(w, http.StatusBadRequest, "Username and password are required")
			return
		}

		// Find user with timeout context
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		user, err := s.userRepo.FindByUsername(ctx, req.Username)
		if err != nil {
			log.Printf("Login failed for user %s: %v", req.Username, err)
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Verify password
		if !checkPasswordHash(req.Password, user.PasswordHash) {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Generate JWT token
		token, err := generateJWT(user.Username)
		if err != nil {
			log.Printf("Failed to generate token: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"token":   token,
			"message": "Login successful",
		})
	}
}

// handleProtected is an example protected endpoint
func (s *Server) handleProtected() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract username from context (set by authMiddleware)
		username, ok := r.Context().Value(contextKeyUsername).(string)
		if !ok {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"message":  "This is a protected endpoint",
			"username": username,
		})
	}
}
