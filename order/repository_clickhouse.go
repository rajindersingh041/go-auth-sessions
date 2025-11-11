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
func NewClickHouseRepository(db *sql.DB) OrderRepository {
	return &ClickHouseRepository{db: db}
}

func (r *ClickHouseRepository) Create(ctx context.Context, order *Order) error {
	// For ClickHouse, we need to generate an ID since it doesn't have auto-increment
	var maxID uint64
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(order_id), 0) FROM orders").Scan(&maxID)
	if err != nil {
		maxID = 0 // Start with ID 1 if no orders exist
	}
	
	order.OrderID = maxID + 1
	
	// Insert order
	orderQuery := "INSERT INTO orders (order_id, user_id, subtotal, tax, total, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, orderQuery, order.OrderID, order.UserID, order.Subtotal, order.Tax, order.Total, order.Status, order.CreatedAt)
	if err != nil {
		return err
	}

	// Insert order items
	return r.CreateOrderItems(ctx, order.OrderID, order.Items)
}

func (r *ClickHouseRepository) CreateOrderItems(ctx context.Context, orderID uint64, items []OrderItem) error {
	itemQuery := "INSERT INTO order_items (order_id, product_id, quantity, unit_price, total) VALUES (?, ?, ?, ?, ?)"
	for _, item := range items {
		_, err := r.db.ExecContext(ctx, itemQuery, orderID, item.ProductID, item.Quantity, item.UnitPrice, item.Total)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ClickHouseRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	query := "SELECT order_id, user_id, subtotal, tax, total, status, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderID, &o.UserID, &o.Subtotal, &o.Tax, &o.Total, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		
		// Load order items
		items, err := r.getOrderItems(ctx, o.OrderID)
		if err != nil {
			return nil, err
		}
		o.Items = items
		
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *ClickHouseRepository) GetOrderByID(ctx context.Context, orderID uint64) (*Order, error) {
	query := "SELECT order_id, user_id, subtotal, tax, total, status, created_at FROM orders WHERE order_id = ?"
	row := r.db.QueryRowContext(ctx, query, orderID)
	
	var o Order
	if err := row.Scan(&o.OrderID, &o.UserID, &o.Subtotal, &o.Tax, &o.Total, &o.Status, &o.CreatedAt); err != nil {
		return nil, err
	}
	
	// Load order items
	items, err := r.getOrderItems(ctx, o.OrderID)
	if err != nil {
		return nil, err
	}
	o.Items = items
	
	return &o, nil
}

// getOrderItems retrieves all items for a specific order
func (r *ClickHouseRepository) getOrderItems(ctx context.Context, orderID uint64) ([]OrderItem, error) {
	query := "SELECT product_id, quantity, unit_price, total FROM order_items WHERE order_id = ? ORDER BY product_id"
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.UnitPrice, &item.Total); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}