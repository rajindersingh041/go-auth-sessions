package product

import (
	"context"
)

// Product represents a product in the database
type Product struct {
	ProductID   uint64  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
	CreatedAt   string  `json:"created_at"`
}

// Repository defines the interface for product data operations
type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, productID uint64) (*Product, error)
	GetByCategory(ctx context.Context, category string) ([]Product, error)
	UpdateStock(ctx context.Context, productID uint64, inStock bool) error
	SeedSampleProducts(ctx context.Context) error
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
}