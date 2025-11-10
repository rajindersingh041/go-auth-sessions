package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ClickHouseUserRepository implements UserRepository for ClickHouse
// (moved from repository.go)
type ClickHouseUserRepository struct {
	db *sql.DB
}

func NewClickHouseUserRepository(db *sql.DB) UserRepository {
	return &ClickHouseUserRepository{db: db}
}

func (r *ClickHouseUserRepository) Create(ctx context.Context, username, passwordHash string) error {
	query := "INSERT INTO users (username, password_hash) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *ClickHouseUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
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

func (r *ClickHouseUserRepository) FindUserID(ctx context.Context, username string) (uint64, error) {
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

func (r *ClickHouseUserRepository) UserExists(ctx context.Context, username string) (bool, error) {
	var count int
	query := "SELECT count() FROM users WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// ClickHouseOrderRepository implements OrderRepository for ClickHouse
type ClickHouseOrderRepository struct {
	db *sql.DB
}

func NewClickHouseOrderRepository(db *sql.DB) OrderRepository {
	return &ClickHouseOrderRepository{db: db}
}

func (r *ClickHouseOrderRepository) Create(ctx context.Context, order *Order) error {
	query := "INSERT INTO orders (user_id, item, quantity, created_at) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Item, order.Quantity, order.CreatedAt)
	return err
}

func (r *ClickHouseOrderRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	query := "SELECT order_id, user_id, item, quantity, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderID, &o.UserID, &o.Item, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
