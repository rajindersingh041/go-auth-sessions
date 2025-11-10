package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/user"
)

// Handler handles HTTP requests for order operations
type Handler struct {
	service     Service
	userService user.Service
	jwtManager  auth.JWTManager
}

// NewHandler creates a new order handler
func NewHandler(service Service, userService user.Service, jwtManager auth.JWTManager) *Handler {
	return &Handler{
		service:     service,
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// RegisterRoutes registers all order-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// All order routes require authentication
	mux.Handle("GET /orders", h.requireAuth(http.HandlerFunc(h.handleGetOrders())))
	mux.Handle("POST /orders", h.requireAuth(http.HandlerFunc(h.handleCreateOrderModern())))
	// Legacy routes for backward compatibility
	mux.Handle("GET /orders/", h.requireAuth(http.HandlerFunc(h.handleGetOrdersByUsername())))
	mux.Handle("POST /orders/", h.requireAuth(http.HandlerFunc(h.handleCreateOrder())))
}

// Context key for username
type contextKey string
const usernameContextKey contextKey = "username"

// requireAuth is a middleware that checks for valid JWT token
func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header required. Please provide JWT token.")
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondError(w, http.StatusUnauthorized, "Authorization header must start with 'Bearer '. Format: 'Bearer <token>'")
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			respondError(w, http.StatusUnauthorized, "JWT token is required. Please login first.")
			return
		}

		// Validate token
		username, err := h.jwtManager.ValidateToken(token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid or expired JWT token. Please login again.")
			return
		}

		// Add username to request context for use in handlers  
		ctx := r.Context()
		ctx = context.WithValue(ctx, usernameContextKey, username)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleGetOrdersByUsername handles requests to fetch orders by username
// URL pattern: GET /orders/{username}
func (h *Handler) handleGetOrdersByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract username from URL path: /orders/{username}
		username := ""
		if len(r.URL.Path) > len("/orders/") {
			username = r.URL.Path[len("/orders/"):]
		}
		if username == "" {
			respondError(w, http.StatusBadRequest, "Username required in path")
			return
		}

		ctx := r.Context()
		
		// Get user to validate existence and get user ID
		user, err := h.userService.GetUserByUsername(ctx, username)
		if err != nil || user == nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}

		// Fetch orders for the user
		orders, err := h.service.GetOrdersByUserID(ctx, user.UserID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch orders")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"orders": orders,
			"count":  len(orders),
		})
	}
}

// handleCreateOrder handles requests to create a new order
// URL pattern: POST /orders/{username}
func (h *Handler) handleCreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract username from URL path: /orders/{username}
		username := ""
		if len(r.URL.Path) > len("/orders/") {
			username = r.URL.Path[len("/orders/"):]
		}
		if username == "" {
			respondError(w, http.StatusBadRequest, "Username required in path")
			return
		}

		ctx := r.Context()
		
		// Get user to validate existence and get user ID
		user, err := h.userService.GetUserByUsername(ctx, username)
		if err != nil || user == nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}

		// Parse request body
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Create the order
		order, err := h.service.CreateOrder(ctx, user.UserID, req)
		if err != nil {
			// Log the actual error for debugging
			log.Printf("Order creation failed: %v", err)
			
			if strings.Contains(err.Error(), "valid product ID and positive quantity are required") ||
			strings.Contains(err.Error(), "product not found") ||
			strings.Contains(err.Error(), "out of stock") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			// Return the actual error message for debugging
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
			return
		}

		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Order created successfully",
			"order":   order,
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

// handleGetOrders handles requests to get orders for authenticated user
// URL pattern: GET /orders (uses JWT token to identify user)
func (h *Handler) handleGetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get username from context (set by requireAuth middleware)
		username, ok := r.Context().Value(usernameContextKey).(string)
		if !ok || username == "" {
			respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		ctx := r.Context()
		
		// Get user to validate existence and get user ID
		user, err := h.userService.GetUserByUsername(ctx, username)
		if err != nil || user == nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}

		// Fetch orders for the user
		orders, err := h.service.GetOrdersByUserID(ctx, user.UserID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch orders")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"orders": orders,
			"count":  len(orders),
		})
	}
}

// handleCreateOrderModern handles requests to create a new order using JWT authentication
// URL pattern: POST /orders (uses JWT token to identify user)
func (h *Handler) handleCreateOrderModern() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get username from context (set by requireAuth middleware)
		username, ok := r.Context().Value(usernameContextKey).(string)
		if !ok || username == "" {
			respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		ctx := r.Context()
		
		// Get user to validate existence and get user ID
		user, err := h.userService.GetUserByUsername(ctx, username)
		if err != nil || user == nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}

		// Parse request body
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Create the order
		order, err := h.service.CreateOrder(ctx, user.UserID, req)
		if err != nil {
			// Log the actual error for debugging
			log.Printf("Order creation failed: %v", err)
			
			// Check if it's a validation error (contains specific messages)
			if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "out of stock") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			// Return the actual error message for debugging
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
			return
		}

		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Order created successfully",
			"order":   order,
		})
	}
}