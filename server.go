package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Server holds dependencies and HTTP handlers
type Server struct {
	userRepo *UserRepository
	mux      *http.ServeMux
}

// NewServer creates a new server instance with all routes configured
func NewServer(userRepo *UserRepository) *Server {
	s := &Server{
		userRepo: userRepo,
		mux:      http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

// Router returns the configured HTTP handler with all middleware applied
func (s *Server) Router() http.Handler {
	// Apply global middleware
	handler := s.loggingMiddleware(s.recoveryMiddleware(s.mux))
	return handler
}

// registerRoutes sets up all HTTP routes
func (s *Server) registerRoutes() {
	// Health check endpoint
	s.mux.HandleFunc("GET /health", s.handleHealth())

	// Authentication endpoints
	s.mux.HandleFunc("POST /register", s.handleRegister())
	s.mux.HandleFunc("POST /login", s.handleLogin())

	// Protected endpoint (requires authentication)
	s.mux.Handle("GET /protected", s.authMiddleware(http.HandlerFunc(s.handleProtected())))
}

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

// loggingMiddleware logs all incoming requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log request details
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// recoveryMiddleware recovers from panics and returns a 500 error
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				respondError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// respondError writes a JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}