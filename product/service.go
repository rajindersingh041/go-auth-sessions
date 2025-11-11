package product

import (
	"context"
	"fmt"
	"time"
)

// ProductService defines the business logic interface for product operations
type ProductService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) error
	GetAllProducts(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, productID uint64) (*Product, error)
	GetProductsByCategory(ctx context.Context, category string) ([]Product, error)
	UpdateProductStock(ctx context.Context, productID uint64, inStock bool) error
	InitializeSampleProducts(ctx context.Context) error
}

// productService implements the ProductService interface
type productService struct {
	repo ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// CreateProduct creates a new product with validation
func (s *productService) CreateProduct(ctx context.Context, req CreateProductRequest) error {
	// Validate input
	if req.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if req.Price < 0 {
		return fmt.Errorf("product price must be non-negative")
	}
	if req.Category == "" {
		return fmt.Errorf("product category is required")
	}

	// Create product
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		InStock:     req.InStock,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	return s.repo.Create(ctx, product)
}

// GetAllProducts retrieves all products
func (s *productService) GetAllProducts(ctx context.Context) ([]Product, error) {
	return s.repo.GetAll(ctx)
}

// GetProductByID retrieves a specific product by ID
func (s *productService) GetProductByID(ctx context.Context, productID uint64) (*Product, error) {
	if productID == 0 {
		return nil, fmt.Errorf("valid product ID is required")
	}
	return s.repo.GetByID(ctx, productID)
}

// GetProductsByCategory retrieves products by category
func (s *productService) GetProductsByCategory(ctx context.Context, category string) ([]Product, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}
	return s.repo.GetByCategory(ctx, category)
}

// UpdateProductStock updates the stock status of a product
func (s *productService) UpdateProductStock(ctx context.Context, productID uint64, inStock bool) error {
	if productID == 0 {
		return fmt.Errorf("valid product ID is required")
	}
	return s.repo.UpdateStock(ctx, productID, inStock)
}

// InitializeSampleProducts creates sample products if none exist
func (s *productService) InitializeSampleProducts(ctx context.Context) error {
	return s.repo.SeedSampleProducts(ctx)
}