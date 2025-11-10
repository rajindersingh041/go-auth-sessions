package order

import (
	"context"
	"database/sql"
)

// PostgresRepository implements Repository for PostgreSQL database
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL order repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ensureOrdersTable creates the orders table with the correct schema
func (r *PostgresRepository) ensureOrdersTable(ctx context.Context) error {
	// Check if the table exists with the old schema
	checkOldSchemaQuery := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'orders' AND column_name = 'item'`
	
	var oldColumn string
	err := r.db.QueryRowContext(ctx, checkOldSchemaQuery).Scan(&oldColumn)
	
	// If old schema exists, drop and recreate the table
	if err == nil {
		// Drop the old table to recreate with correct schema
		dropQuery := `DROP TABLE IF EXISTS orders`
		if _, err := r.db.ExecContext(ctx, dropQuery); err != nil {
			return err
		}
	}

	// Create the table with the correct schema
	createQuery := `
		CREATE TABLE IF NOT EXISTS orders (
			order_id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			product_id BIGINT NOT NULL,
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`
	_, err = r.db.ExecContext(ctx, createQuery)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, order *Order) error {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO orders (user_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4) RETURNING order_id"
	err := r.db.QueryRowContext(ctx, query, order.UserID, order.ProductID, order.Quantity, order.CreatedAt).Scan(&order.OrderID)
	return err
}

func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderID, &o.UserID, &o.ProductID, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *PostgresRepository) GetOrderByID(ctx context.Context, orderID uint64) (*Order, error) {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE order_id = $1"
	row := r.db.QueryRowContext(ctx, query, orderID)
	
	var o Order
	if err := row.Scan(&o.OrderID, &o.UserID, &o.ProductID, &o.Quantity, &o.CreatedAt); err != nil {
		return nil, err
	}
	return &o, nil
}