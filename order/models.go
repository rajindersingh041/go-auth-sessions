package order

import (
	"context"
)

// Order represents an order placed by a user
type Order struct {
	OrderID   uint64
	UserID    uint64
	Item      string
	Quantity  int
	CreatedAt string // or time.Time if you want to use time
}

// Repository defines the interface for order data operations
type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}