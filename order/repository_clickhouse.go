package order

import (
	"context"
	"database/sql"
)

// ClickHouseRepository implements Repository for ClickHouse database
type ClickHouseRepository struct {
	db *sql.DB
}

// NewClickHouseRepository creates a new ClickHouse order repository
func NewClickHouseRepository(db *sql.DB) Repository {
	return &ClickHouseRepository{db: db}
}

func (r *ClickHouseRepository) Create(ctx context.Context, order *Order) error {
	// For ClickHouse, we need to generate an ID since it doesn't have auto-increment
	// Get the next ID by finding the max existing ID
	var maxID uint64
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(order_id), 0) FROM orders").Scan(&maxID)
	if err != nil {
		// If table doesn't exist yet, start with ID 1
		maxID = 0
	}
	
	order.OrderID = maxID + 1
	query := "INSERT INTO orders (order_id, user_id, product_id, quantity, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, order.OrderID, order.UserID, order.ProductID, order.Quantity, order.CreatedAt)
	return err
}

func (r *ClickHouseRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	query := "SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC"
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

func (r *ClickHouseRepository) GetOrderByID(ctx context.Context, orderID uint64) (*Order, error) {
	query := "SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE order_id = ?"
	row := r.db.QueryRowContext(ctx, query, orderID)
	
	var o Order
	if err := row.Scan(&o.OrderID, &o.UserID, &o.ProductID, &o.Quantity, &o.CreatedAt); err != nil {
		return nil, err
	}
	return &o, nil
}