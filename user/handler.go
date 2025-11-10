package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/auth"
)

// Handler handles HTTP requests for user operations
type Handler struct {
	service    Service
	jwtManager auth.JWTManager
}

// NewHandler creates a new user handler
func NewHandler(service Service, jwtManager auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers all user-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.handleRegister())
	mux.HandleFunc("POST /login", h.handleLogin())
}

// handleRegister handles user registration requests
func (h *Handler) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		if err := h.service.CreateUser(ctx, req); err != nil {
			// Check for specific error types to return appropriate status codes
			if err.Error() == "user already exists" {
				respondError(w, http.StatusConflict, err.Error())
				return
			}
			if err.Error() == "username and password are required" {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		respondJSON(w, http.StatusCreated, map[string]string{
			"message": "User created successfully",
		})
	}
}

// handleLogin handles user login requests
func (h *Handler) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		user, err := h.service.AuthenticateUser(ctx, req)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Generate JWT token
		token, err := h.jwtManager.GenerateToken(user.Username)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Login successful",
			"token":   token,
			"user": map[string]interface{}{
				"id":       user.UserID,
				"username": user.Username,
			},
		})
	}
}

// Helper functions for HTTP responses
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error":     message,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}