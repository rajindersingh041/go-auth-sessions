package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// UserRepository handles all database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) error {
	query := "INSERT INTO users (username, password_hash) VALUES (?, ?)"

	_, err := r.db.ExecContext(ctx, query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByUsername retrieves a user by their username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	query := "SELECT user_id, username, password_hash FROM users WHERE username = ? LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// FindUserID retrieves only the user ID by username
func (r *UserRepository) FindUserID(ctx context.Context, username string) (uint64, error) {
	var userID uint64

	query := "SELECT user_id FROM users WHERE username = ? LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to query user ID: %w", err)
	}

	return userID, nil
}

// UserExists checks if a user with the given username exists
func (r *UserRepository) UserExists(ctx context.Context, username string) (bool, error) {
	var count int

	query := "SELECT count() FROM users WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}
