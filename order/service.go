package order

import (
	"context"
	"fmt"
	"time"
)

// Service defines the business logic interface for order operations
type Service interface {
	CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) error
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new order service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateOrder creates a new order with validation
func (s *service) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) error {
	// Validate input
	if req.Item == "" || req.Quantity <= 0 {
		return fmt.Errorf("item and positive quantity are required")
	}
	if userID == 0 {
		return fmt.Errorf("valid user ID is required")
	}

	// Create order
	order := &Order{
		UserID:    userID,
		Item:      req.Item,
		Quantity:  req.Quantity,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return s.repo.Create(ctx, order)
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *service) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	if userID == 0 {
		return nil, fmt.Errorf("valid user ID is required")
	}
	return s.repo.GetOrdersByUserID(ctx, userID)
}