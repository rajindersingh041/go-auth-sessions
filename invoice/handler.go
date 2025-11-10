package invoice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rajindersingh041/go-auth-sessions/auth"
)

// Handler handles HTTP requests for invoice operations
type Handler struct {
	service    Service
	jwtManager auth.JWTManager
}

// NewHandler creates a new invoice handler
func NewHandler(service Service, jwtManager auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers all invoice-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// All invoice routes require authentication
	mux.Handle("POST /invoices", h.requireAuth(http.HandlerFunc(h.handleCreateInvoice())))
	mux.Handle("GET /invoices/", h.requireAuth(http.HandlerFunc(h.handleGetInvoice())))
	mux.Handle("GET /invoices/user/", h.requireAuth(http.HandlerFunc(h.handleGetUserInvoices())))
	mux.Handle("PUT /invoices/", h.requireAuth(http.HandlerFunc(h.handleUpdateInvoiceStatus())))
}

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
		// You can add username to context if needed
		_ = username

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleCreateInvoice handles requests to create an invoice from an order
// URL pattern: POST /invoices
func (h *Handler) handleCreateInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse request body
		var req CreateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Create the invoice
		invoice, err := h.service.CreateInvoiceFromOrder(ctx, req.OrderID)
		if err != nil {
			log.Printf("Invoice creation failed: %v", err)
			if strings.Contains(err.Error(), "valid order ID is required") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create invoice: %v", err))
			return
		}

		respondJSON(w, http.StatusCreated, invoice)
	}
}

// handleGetInvoice handles requests to get an invoice by ID or order ID
// URL patterns: 
// - GET /invoices/{id} - Get by invoice ID
// - GET /invoices/order/{order_id} - Get by order ID
func (h *Handler) handleGetInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse URL path
		path := strings.TrimPrefix(r.URL.Path, "/invoices/")
		if path == "" {
			respondError(w, http.StatusBadRequest, "Invoice ID or order ID required in path")
			return
		}

		var invoice *Invoice
		var err error

		// Check if it's an order lookup
		if strings.HasPrefix(path, "order/") {
			orderIDStr := strings.TrimPrefix(path, "order/")
			orderID, parseErr := strconv.ParseUint(orderIDStr, 10, 64)
			if parseErr != nil {
				respondError(w, http.StatusBadRequest, "Invalid order ID")
				return
			}
			invoice, err = h.service.GetInvoiceByOrderID(ctx, orderID)
		} else {
			// Assume it's an invoice ID
			invoiceID, parseErr := strconv.ParseUint(path, 10, 64)
			if parseErr != nil {
				respondError(w, http.StatusBadRequest, "Invalid invoice ID")
				return
			}
			invoice, err = h.service.GetInvoiceByID(ctx, invoiceID)
		}

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				respondError(w, http.StatusNotFound, "Invoice not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to retrieve invoice")
			return
		}

		respondJSON(w, http.StatusOK, invoice)
	}
}

// handleGetUserInvoices handles requests to get all invoices for a user
// URL pattern: GET /invoices/user/{user_id}
func (h *Handler) handleGetUserInvoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract user ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/invoices/user/")
		if path == "" {
			respondError(w, http.StatusBadRequest, "User ID required in path")
			return
		}

		userID, err := strconv.ParseUint(path, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Get invoices
		invoices, err := h.service.GetInvoicesByUserID(ctx, userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve invoices")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"invoices": invoices,
			"count":    len(invoices),
		})
	}
}

// handleUpdateInvoiceStatus handles requests to update invoice status
// URL pattern: PUT /invoices/{id}
func (h *Handler) handleUpdateInvoiceStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract invoice ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/invoices/")
		if path == "" {
			respondError(w, http.StatusBadRequest, "Invoice ID required in path")
			return
		}

		invoiceID, err := strconv.ParseUint(path, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid invoice ID")
			return
		}

		// Parse request body
		var req UpdateInvoiceStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Update status
		if err := h.service.UpdateInvoiceStatus(ctx, invoiceID, req.Status); err != nil {
			if strings.Contains(err.Error(), "invalid status") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to update invoice status")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"message": "Invoice status updated successfully",
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":     message,
		"timestamp": "2025-11-10T22:41:03+01:00", // You can use time.Now() here
	})
}