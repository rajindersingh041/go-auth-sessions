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

// ensureOrdersTable creates the orders table with the new multi-item schema
func (r *PostgresRepository) ensureOrdersTable(ctx context.Context) error {
	// Create the orders table
	createOrdersQuery := `
		CREATE TABLE IF NOT EXISTS orders (
			order_id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			subtotal DECIMAL(10,2) NOT NULL DEFAULT 0.00,
			tax DECIMAL(10,2) NOT NULL DEFAULT 0.00,
			total DECIMAL(10,2) NOT NULL DEFAULT 0.00,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`
	if _, err := r.db.ExecContext(ctx, createOrdersQuery); err != nil {
		return err
	}

	// Create the order_items table
	createOrderItemsQuery := `
		CREATE TABLE IF NOT EXISTS order_items (
			item_id SERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
			product_id BIGINT NOT NULL,
			quantity INT NOT NULL,
			unit_price DECIMAL(10,2) NOT NULL,
			total DECIMAL(10,2) NOT NULL
		)`
	_, err := r.db.ExecContext(ctx, createOrderItemsQuery)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, order *Order) error {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := "INSERT INTO orders (user_id, subtotal, tax, total, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING order_id"
	err = tx.QueryRowContext(ctx, orderQuery, order.UserID, order.Subtotal, order.Tax, order.Total, order.Status, order.CreatedAt).Scan(&order.OrderID)
	if err != nil {
		return err
	}

	// Insert order items
	if err := r.createOrderItemsInTx(ctx, tx, order.OrderID, order.Items); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (r *PostgresRepository) CreateOrderItems(ctx context.Context, orderID uint64, items []OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.createOrderItemsInTx(ctx, tx, orderID, items); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) createOrderItemsInTx(ctx context.Context, tx *sql.Tx, orderID uint64, items []OrderItem) error {
	itemQuery := "INSERT INTO order_items (order_id, product_id, quantity, unit_price, total) VALUES ($1, $2, $3, $4, $5)"
	for _, item := range items {
		_, err := tx.ExecContext(ctx, itemQuery, orderID, item.ProductID, item.Quantity, item.UnitPrice, item.Total)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT order_id, user_id, subtotal, tax, total, status, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC"
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

func (r *PostgresRepository) GetOrderByID(ctx context.Context, orderID uint64) (*Order, error) {
	if err := r.ensureOrdersTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT order_id, user_id, subtotal, tax, total, status, created_at FROM orders WHERE order_id = $1"
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
func (r *PostgresRepository) getOrderItems(ctx context.Context, orderID uint64) ([]OrderItem, error) {
	query := "SELECT product_id, quantity, unit_price, total FROM order_items WHERE order_id = $1 ORDER BY item_id"
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