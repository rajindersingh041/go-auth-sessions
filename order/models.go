package order

import (
	"context"
)

// Order represents an order placed by a user
type Order struct {
	OrderID   uint64
	UserID    uint64
	ProductID uint64  // Changed from Item to ProductID
	Quantity  int
	CreatedAt string // or time.Time if you want to use time
}

// Repository defines the interface for order data operations
type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
	GetOrderByID(ctx context.Context, orderID uint64) (*Order, error)
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Order now includes ProductID instead of Item
type OrderWithProduct struct {
	OrderID   uint64  `json:"order_id"`
	UserID    uint64  `json:"user_id"`
	ProductID uint64  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	CreatedAt string  `json:"created_at"`
	// Additional product info for responses
	ProductName  string  `json:"product_name,omitempty"`
	ProductPrice float64 `json:"product_price,omitempty"`
}