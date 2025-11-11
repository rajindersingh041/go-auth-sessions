package user

import (
	"context"
	"fmt"
)

// UserService defines the business logic interface for user operations
type UserService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) error
	AuthenticateUser(ctx context.Context, req LoginRequest) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
}

// PasswordHasher interface for password hashing operations
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) error
}

// userService implements the UserService interface
type userService struct {
	repo           UserRepository
	passwordHasher PasswordHasher
}

// NewUserService creates a new user service
func NewUserService(repo UserRepository, passwordHasher PasswordHasher) UserService {
       return &userService{
	       repo:           repo,
	       passwordHasher: passwordHasher,
       }
}

// CreateUser creates a new user with validation and password hashing
func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) error {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return fmt.Errorf("username and password are required")
	}

	// Check if user already exists
	exists, err := s.repo.UserExists(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	return s.repo.Create(ctx, req.Username, hashedPassword)
}

// AuthenticateUser authenticates a user and returns user info if successful
func (s *userService) AuthenticateUser(ctx context.Context, req LoginRequest) (*User, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// Find user
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("authentication failed")
	}
	if user == nil {
		return nil, fmt.Errorf("authentication failed")
	}

	// Check password
	if err := s.passwordHasher.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, fmt.Errorf("authentication failed")
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	return s.repo.FindByUsername(ctx, username)
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	return s.repo.FindByID(ctx, userID)
}