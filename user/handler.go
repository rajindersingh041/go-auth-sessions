package user

import (
	"encoding/json"
	"net/http"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/helper"
)

// Handler handles HTTP requests for user operations
type Handler struct {
	service    UserService
	jwtManager auth.JWTManager
}

// NewHandler creates a new user handler
func NewHandler(service UserService, jwtManager auth.JWTManager) *Handler {
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
			helper.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		if err := h.service.CreateUser(ctx, req); err != nil {
			// Check for specific error types to return appropriate status codes
			if err.Error() == "user already exists" {
				helper.RespondError(w, http.StatusConflict, err.Error())
				return
			}
			if err.Error() == "username and password are required" {
				helper.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		helper.RespondJSON(w, http.StatusCreated, map[string]string{
			"message": "User created successfully",
		})
	}
}

// handleLogin handles user login requests
func (h *Handler) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		user, err := h.service.AuthenticateUser(ctx, req)
		if err != nil {
			helper.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Generate JWT token
		token, err := h.jwtManager.GenerateToken(user.Username)
		if err != nil {
			helper.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		helper.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Login successful",
			"token":   token,
			"user": map[string]interface{}{
				"id":       user.UserID,
				"username": user.Username,
			},
		})
	}
}
