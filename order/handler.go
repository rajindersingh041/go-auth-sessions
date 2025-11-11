package order

import (
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
	service     OrderService
	userService user.UserService
}

// NewHandler creates a new order handler
func NewHandler(service OrderService, userService user.UserService) *Handler {
	return &Handler{
		service:     service,
		userService: userService,
	}
}

// RegisterRoutes registers all order-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux, jwtManager auth.JWTManager) {
	// Register routes with or without authentication as needed
	mux.Handle("GET /orders", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleGetOrders())))
	mux.Handle("POST /orders", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleCreateOrder())))
	mux.Handle("POST /orders/single", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleCreateSingleOrder())))
	mux.Handle("GET /orders/", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleGetOrdersByUsername())))
	mux.Handle("POST /orders/", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleCreateOrderLegacy())))
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

// handleCreateOrderLegacyPath handles requests to create a new order (legacy username-based path)
// URL pattern: POST /orders/{username}
func (h *Handler) handleCreateOrderLegacyPath() http.HandlerFunc {
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
		// Get username from context (set by WithJWTAuth middleware)
		username, ok := r.Context().Value(auth.UsernameContextKey).(string)
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

// handleCreateOrder handles requests to create a new order with multiple products
// URL pattern: POST /orders (uses JWT token to identify user)
func (h *Handler) handleCreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get username from context (set by WithJWTAuth middleware)
		username, ok := r.Context().Value(auth.UsernameContextKey).(string)
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

		// Parse request body for multi-product order
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
			"total_amount": order.Total,
			"items_count": len(order.Items),
		})
	}
}

// handleCreateSingleOrder handles requests to create a single-product order (legacy compatibility)
// URL pattern: POST /orders/single (uses JWT token to identify user)
func (h *Handler) handleCreateSingleOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get username from context (set by WithJWTAuth middleware)
		username, ok := r.Context().Value(auth.UsernameContextKey).(string)
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

		// Parse request body for single product order
		var req CreateSingleOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Create the order
		order, err := h.service.CreateSingleOrder(ctx, user.UserID, req)
		if err != nil {
			// Log the actual error for debugging
			log.Printf("Single order creation failed: %v", err)
			
			// Check if it's a validation error
			if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "out of stock") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
			return
		}

		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Order created successfully",
			"order":   order,
		})
	}
}

// handleCreateOrderLegacy handles the legacy order creation (same as handleCreateOrderLegacyPath)
func (h *Handler) handleCreateOrderLegacy() http.HandlerFunc {
	return h.handleCreateOrderLegacyPath()
}