package order

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/user"
)

// Handler handles HTTP requests for order operations
type Handler struct {
	service     Service
	userService user.Service
}

// NewHandler creates a new order handler
func NewHandler(service Service, userService user.Service) *Handler {
	return &Handler{
		service:     service,
		userService: userService,
	}
}

// RegisterRoutes registers all order-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /orders/", h.handleGetOrdersByUsername())
	mux.HandleFunc("POST /orders/", h.handleCreateOrder())
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
		if err := h.service.CreateOrder(ctx, user.UserID, req); err != nil {
			if err.Error() == "item and positive quantity are required" {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to create order")
			return
		}

		respondJSON(w, http.StatusCreated, map[string]string{
			"message": "Order created successfully",
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