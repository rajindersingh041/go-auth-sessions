package main

import (
	"net/http"
)

// Server holds all dependencies and HTTP handlers for the application.
type Server struct {
	userRepository UserRepository      // Handles user data operations
	passwordHasher PasswordHasher      // Handles password hashing and verification
	jwtManager     JWTManager          // Handles JWT creation and validation
	httpMux        *http.ServeMux      // HTTP request multiplexer
}

// NewServer creates a new Server instance with all routes and dependencies configured.
func NewServer(userRepository UserRepository, passwordHasher PasswordHasher, jwtManager JWTManager) *Server {
    srv := &Server{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		jwtManager:     jwtManager,
		httpMux:        http.NewServeMux(),
	}
	srv.registerRoutes()
    return srv
}

// Router returns the configured HTTP handler with all middleware applied.
func (srv *Server) Router() http.Handler {
	// Apply global middleware
	handler := srv.loggingMiddleware(srv.recoveryMiddleware(srv.httpMux))
	return handler
}

// registerRoutes sets up all HTTP routes for the server.
func (srv *Server) registerRoutes() {
	// Health check endpoint
	srv.httpMux.HandleFunc("GET /health", srv.handleHealth())

	// Authentication endpoints
	srv.httpMux.HandleFunc("POST /register", srv.handleRegister())
	srv.httpMux.HandleFunc("POST /login", srv.handleLogin())

	// Protected endpoint (requires authentication)
	srv.httpMux.Handle("GET /protected", srv.authMiddleware(http.HandlerFunc(srv.handleProtected())))
}