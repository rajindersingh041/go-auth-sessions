package product

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/helper"
)

// Handler handles HTTP requests for product operations
type Handler struct {
	service    ProductService
	jwtManager auth.JWTManager
}

// NewHandler creates a new product handler
func NewHandler(service ProductService, jwtManager auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers all product-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux, jwtManager auth.JWTManager) {
	// Public routes (no authentication required)
	mux.HandleFunc("GET /products", h.handleGetAllProducts())
	mux.HandleFunc("GET /products/", h.handleGetProductByIDOrCategory())
	
	// Protected routes (authentication required)
	mux.Handle("POST /products", auth.WithJWTAuth(h.jwtManager, http.HandlerFunc(h.handleCreateProduct())))
	mux.Handle("PUT /products/", auth.WithJWTAuth(h.jwtManager, http.HandlerFunc(h.handleUpdateProductStock())))
}


// handleGetAllProducts handles requests to get all products (public)
func (h *Handler) handleGetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		products, err := h.service.GetAllProducts(ctx)
		if err != nil {
			helper.RespondError(w, http.StatusInternalServerError, "Failed to fetch products")
			return
		}

		helper.RespondJSON(w, http.StatusOK, map[string]interface{}{
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
				helper.RespondError(w, http.StatusBadRequest, "Category name is required")
				return
			}
			
			products, err := h.service.GetProductsByCategory(ctx, category)
			if err != nil {
				helper.RespondError(w, http.StatusInternalServerError, "Failed to fetch products by category")
				return
			}
			
			helper.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"products": products,
				"category": category,
				"count":    len(products),
			})
			return
		}

		// Otherwise, treat as product ID
		productIDStr := path
		if productIDStr == "" {
			helper.RespondError(w, http.StatusBadRequest, "Product ID is required")
			return
		}

		productID, err := strconv.ParseUint(productIDStr, 10, 64)
		if err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		product, err := h.service.GetProductByID(ctx, productID)
		if err != nil {
			if err.Error() == "product not found" {
				helper.RespondError(w, http.StatusNotFound, "Product not found")
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, "Failed to fetch product")
			return
		}

		helper.RespondJSON(w, http.StatusOK, product)
	}
}

// handleCreateProduct handles requests to create a new product (protected)
func (h *Handler) handleCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		if err := h.service.CreateProduct(ctx, req); err != nil {
			if strings.Contains(err.Error(), "required") {
				helper.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, "Failed to create product")
			return
		}

		helper.RespondJSON(w, http.StatusCreated, map[string]string{
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
			helper.RespondError(w, http.StatusBadRequest, "Invalid URL format. Use /products/{id}/stock")
			return
		}

		productIDStr := parts[0]
		productID, err := strconv.ParseUint(productIDStr, 10, 64)
		if err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		var req struct {
			InStock bool `json:"in_stock"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		if err := h.service.UpdateProductStock(ctx, productID, req.InStock); err != nil {
			if strings.Contains(err.Error(), "required") {
				helper.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			helper.RespondError(w, http.StatusInternalServerError, "Failed to update product stock")
			return
		}

		helper.RespondJSON(w, http.StatusOK, map[string]string{
			"message": "Product stock updated successfully",
		})
	}
}

