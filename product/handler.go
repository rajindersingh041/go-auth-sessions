package product

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/auth"
)

// Handler handles HTTP requests for product operations
type Handler struct {
	service    Service
	jwtManager auth.JWTManager
}

// NewHandler creates a new product handler
func NewHandler(service Service, jwtManager auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers all product-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Public routes (no authentication required)
	mux.HandleFunc("GET /products", h.handleGetAllProducts())
	mux.HandleFunc("GET /products/", h.handleGetProductByIDOrCategory())
	
	// Protected routes (authentication required)
	mux.Handle("POST /products", h.requireAuth(http.HandlerFunc(h.handleCreateProduct())))
	mux.Handle("PUT /products/", h.requireAuth(http.HandlerFunc(h.handleUpdateProductStock())))
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

// handleGetAllProducts handles requests to get all products (public)
func (h *Handler) handleGetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		products, err := h.service.GetAllProducts(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch products")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"products": products,
			"count":    len(products),
		})
	}
}

// handleGetProductByIDOrCategory handles GET /products/{id} or GET /products/category/{category}
func (h *Handler) handleGetProductByIDOrCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/products/")
		
		ctx := r.Context()
		
		// Check if it's a category request: /products/category/{category}
		if strings.HasPrefix(path, "category/") {
			category := strings.TrimPrefix(path, "category/")
			if category == "" {
				respondError(w, http.StatusBadRequest, "Category name is required")
				return
			}
			
			products, err := h.service.GetProductsByCategory(ctx, category)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "Failed to fetch products by category")
				return
			}
			
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"products": products,
				"category": category,
				"count":    len(products),
			})
			return
		}

		// Otherwise, treat as product ID
		productIDStr := path
		if productIDStr == "" {
			respondError(w, http.StatusBadRequest, "Product ID is required")
			return
		}

		productID, err := strconv.ParseUint(productIDStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		product, err := h.service.GetProductByID(ctx, productID)
		if err != nil {
			if err.Error() == "product not found" {
				respondError(w, http.StatusNotFound, "Product not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to fetch product")
			return
		}

		respondJSON(w, http.StatusOK, product)
	}
}

// handleCreateProduct handles requests to create a new product (protected)
func (h *Handler) handleCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		if err := h.service.CreateProduct(ctx, req); err != nil {
			if strings.Contains(err.Error(), "required") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to create product")
			return
		}

		respondJSON(w, http.StatusCreated, map[string]string{
			"message": "Product created successfully",
		})
	}
}

// handleUpdateProductStock handles requests to update product stock (protected)
func (h *Handler) handleUpdateProductStock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract product ID from URL: /products/{id}/stock
		path := strings.TrimPrefix(r.URL.Path, "/products/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 || parts[1] != "stock" {
			respondError(w, http.StatusBadRequest, "Invalid URL format. Use /products/{id}/stock")
			return
		}

		productIDStr := parts[0]
		productID, err := strconv.ParseUint(productIDStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		var req struct {
			InStock bool `json:"in_stock"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		if err := h.service.UpdateProductStock(ctx, productID, req.InStock); err != nil {
			if strings.Contains(err.Error(), "required") {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to update product stock")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"message": "Product stock updated successfully",
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