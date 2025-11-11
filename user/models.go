package user

import (
	"context"
)

// User represents a user in the database
type User struct {
	UserID       uint64
	Username     string
	EmailID		 string
	PasswordHash string
}

// Repository defines the interface for user data operations
// why Context?
// Context is used to carry deadlines, cancellation signals,
// and other request-scoped values across API boundaries and between processes.
// This allows for better control over request lifecycles
// and resource management in database operations.
type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, userID uint64) (*User, error)
	FindUserID(ctx context.Context, username string) (uint64, error)
	UserExists(ctx context.Context, username string) (bool, error)
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}