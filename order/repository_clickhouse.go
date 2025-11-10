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
	query := "INSERT INTO orders (user_id, item, quantity, created_at) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Item, order.Quantity, order.CreatedAt)
	return err
}

func (r *ClickHouseRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
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