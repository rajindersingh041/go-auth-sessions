package main

import "context"

// "context"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindUserID(ctx context.Context, username string) (uint64, error)
	UserExists(ctx context.Context, username string) (bool, error)
}

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
}
