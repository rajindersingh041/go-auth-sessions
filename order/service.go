package order

import (
	"context"
	"fmt"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/product"
)

// Service defines the business logic interface for order operations
type Service interface {
	CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
	GetOrderByID(ctx context.Context, orderID uint64) (*Order, error)
}

// service implements the Service interface
type service struct {
	repo           Repository
	productService product.Service
}

// NewService creates a new order service
func NewService(repo Repository, productService product.Service) Service {
	return &service{
		repo:           repo,
		productService: productService,
	}
}

// CreateOrder creates a new order with validation
func (s *service) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*Order, error) {
	// Validate input
	if req.ProductID == 0 || req.Quantity <= 0 {
		return nil, fmt.Errorf("valid product ID and positive quantity are required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("valid user ID is required")
	}

	// Validate product exists and is in stock
	prod, err := s.productService.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if !prod.InStock {
		return nil, fmt.Errorf("product '%s' is currently out of stock", prod.Name)
	}

	// Create order
	order := &Order{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Return the created order - Note: In a real system, you'd get the ID from the database
	// For now, we'll return the order as is, but the OrderID won't be set
	return order, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *service) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	if userID == 0 {
		return nil, fmt.Errorf("valid user ID is required")
	}
	return s.repo.GetOrdersByUserID(ctx, userID)
}

// GetOrderByID retrieves a specific order by ID
func (s *service) GetOrderByID(ctx context.Context, orderID uint64) (*Order, error) {
	if orderID == 0 {
		return nil, fmt.Errorf("valid order ID is required")
	}
	return s.repo.GetOrderByID(ctx, orderID)
}