package invoice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/helper"
)

// Handler handles HTTP requests for invoice operations
type Handler struct {
	service    InvoiceService
	jwtManager auth.JWTManager
}

// NewHandler creates a new invoice handler
func NewHandler(service InvoiceService, jwtManager auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers all invoice-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux, jwtManager auth.JWTManager) {
	// All invoice routes require authentication
	mux.Handle("POST /invoices", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleCreateInvoice())))
	mux.Handle("GET /invoices/", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleGetInvoice())))
	mux.Handle("GET /invoices/user/", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleGetUserInvoices())))
	mux.Handle("PUT /invoices/", auth.WithJWTAuth(jwtManager, http.HandlerFunc(h.handleUpdateInvoiceStatus())))
}


// handleCreateInvoice handles requests to create an invoice from an order
// URL pattern: POST /invoices
func (h *Handler) handleCreateInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse request body
		var req CreateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Create the invoice
		invoice, err := h.service.CreateInvoiceFromOrder(ctx, req.OrderID)
		if err != nil {
			log.Printf("Invoice creation failed: %v", err)
			if strings.Contains(err.Error(), "valid order ID is required") {
				helper.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create invoice: %v", err))
			return
		}

		helper.RespondJSON(w, http.StatusCreated, invoice)
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
			helper.RespondError(w, http.StatusBadRequest, "Invoice ID or order ID required in path")
			return
		}

		var invoice *Invoice
		var err error

		// Check if it's an order lookup
		if strings.HasPrefix(path, "order/") {
			orderIDStr := strings.TrimPrefix(path, "order/")
			orderID, parseErr := strconv.ParseUint(orderIDStr, 10, 64)
			if parseErr != nil {
				helper.RespondError(w, http.StatusBadRequest, "Invalid order ID")
				return
			}
			invoice, err = h.service.GetInvoiceByOrderID(ctx, orderID)
		} else {
			// Assume it's an invoice ID
			invoiceID, parseErr := strconv.ParseUint(path, 10, 64)
			if parseErr != nil {
				helper.RespondError(w, http.StatusBadRequest, "Invalid invoice ID")
				return
			}
			invoice, err = h.service.GetInvoiceByID(ctx, invoiceID)
		}

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				helper.RespondError(w, http.StatusNotFound, "Invoice not found")
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, "Failed to retrieve invoice")
			return
		}

		helper.RespondJSON(w, http.StatusOK, invoice)
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
			helper.RespondError(w, http.StatusBadRequest, "User ID required in path")
			return
		}

		userID, err := strconv.ParseUint(path, 10, 64)
		if err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Get invoices
		invoices, err := h.service.GetInvoicesByUserID(ctx, userID)
		if err != nil {
			helper.RespondError(w, http.StatusInternalServerError, "Failed to retrieve invoices")
			return
		}

		helper.RespondJSON(w, http.StatusOK, map[string]interface{}{
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
			helper.RespondError(w, http.StatusBadRequest, "Invoice ID required in path")
			return
		}

		invoiceID, err := strconv.ParseUint(path, 10, 64)
		if err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid invoice ID")
			return
		}

		// Parse request body
		var req UpdateInvoiceStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Update status
		if err := h.service.UpdateInvoiceStatus(ctx, invoiceID, req.Status); err != nil {
			if strings.Contains(err.Error(), "invalid status") {
				helper.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, "Failed to update invoice status")
			return
		}

		helper.RespondJSON(w, http.StatusOK, map[string]string{
			"message": "Invoice status updated successfully",
		})
	}
}
