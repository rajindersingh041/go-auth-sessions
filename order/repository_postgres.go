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

// ensureOrdersTable creates the orders table if it doesn't exist
func (r *PostgresRepository) ensureOrdersTable(ctx context.Context) error {
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

func (r *PostgresRepository) Create(ctx context.Context, order *Order) error {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO orders (user_id, item, quantity, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Item, order.Quantity, order.CreatedAt)
	return err
}

func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
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