package main

import (
	"net/http"
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