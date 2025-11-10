package main

import (
	"context"
	"database/sql"
)

// PostgresUserRepository implements UserRepository for Postgres (stub)
type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}


// ensureUsersTable creates the users table if it doesn't exist (id SERIAL PRIMARY KEY, username unique, password_hash)
func (r *PostgresUserRepository) ensureUsersTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	)`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresUserRepository) Create(ctx context.Context, username, passwordHash string) error {
	if err := r.ensureUsersTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, username, passwordHash)
	return err
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
    if err := r.ensureUsersTable(ctx); err != nil {
		return nil, err
	}
	var user User
	query := "SELECT user_id, username, password_hash FROM users WHERE username = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.UserID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindUserID(ctx context.Context, username string) (uint64, error) {
	if err := r.ensureUsersTable(ctx); err != nil {
		return 0, err
	}
	var userID uint64
	query := "SELECT user_id FROM users WHERE username = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func (r *PostgresUserRepository) UserExists(ctx context.Context, username string) (bool, error) {
	if err := r.ensureUsersTable(ctx); err != nil {
		return false, err
	}
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PostgresOrderRepository implements OrderRepository for Postgres
type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) OrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) ensureOrdersTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS orders (
			order_id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			item TEXT NOT NULL,
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *Order) error {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO orders (user_id, item, quantity, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Item, order.Quantity, order.CreatedAt)
	return err
}

func (r *PostgresOrderRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT order_id, user_id, item, quantity, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC"
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
