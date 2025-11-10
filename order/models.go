package order

import (
	"context"
)

// OrderItem represents a single product within an order
type OrderItem struct {
	ProductID uint64  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price,omitempty"`
	Total     float64 `json:"total,omitempty"`
}

// Order represents an order placed by a user (can contain multiple products)
type Order struct {
	OrderID   uint64      `json:"order_id"`
	UserID    uint64      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Subtotal  float64     `json:"subtotal"`
	Tax       float64     `json:"tax"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
}

// Repository defines the interface for order data operations
type Repository interface {
	Create(ctx context.Context, order *Order) error
	CreateOrderItems(ctx context.Context, orderID uint64, items []OrderItem) error
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
	GetOrderByID(ctx context.Context, orderID uint64) (*Order, error)
}

// CreateOrderRequest represents the request to create an order with multiple products
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

// OrderItemRequest represents a product to add to an order
type OrderItemRequest struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Legacy single product order request (for backward compatibility)
type CreateSingleOrderRequest struct {
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