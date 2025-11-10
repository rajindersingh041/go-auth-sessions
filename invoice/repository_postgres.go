package invoice

import (
	"context"
	"database/sql"
	"encoding/json"
)

// PostgresRepository implements Repository for PostgreSQL database
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL invoice repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ensureInvoicesTable creates the invoices table if it doesn't exist
func (r *PostgresRepository) ensureInvoicesTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS invoices (
			invoice_id SERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			username TEXT NOT NULL,
			invoice_number TEXT NOT NULL UNIQUE,
			items JSONB,
			subtotal DECIMAL(10,2) DEFAULT 0.00,
			tax DECIMAL(10,2) DEFAULT 0.00,
			total DECIMAL(10,2) DEFAULT 0.00,
			status TEXT NOT NULL DEFAULT 'draft',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			due_date TIMESTAMP NOT NULL
		)`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, invoice *Invoice) error {
	if err := r.ensureInvoicesTable(ctx); err != nil {
		return err
	}

	// Convert items to JSON
	itemsJSON, err := json.Marshal(invoice.Items)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO invoices (order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING invoice_id`
	
	err = r.db.QueryRowContext(ctx, query, 
		invoice.OrderID, 
		invoice.UserID, 
		invoice.Username,
		invoice.InvoiceNumber, 
		itemsJSON,
		invoice.Subtotal,
		invoice.Tax,
		invoice.Total,
		invoice.Status, 
		invoice.CreatedAt, 
		invoice.DueDate).Scan(&invoice.InvoiceID)
		
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, invoiceID uint64) (*Invoice, error) {
	if err := r.ensureInvoicesTable(ctx); err != nil {
		return nil, err
	}

	query := `
		SELECT invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date 
		FROM invoices WHERE invoice_id = $1`
	
	return r.scanInvoice(ctx, query, invoiceID)
}

func (r *PostgresRepository) GetByOrderID(ctx context.Context, orderID uint64) (*Invoice, error) {
	if err := r.ensureInvoicesTable(ctx); err != nil {
		return nil, err
	}

	query := `
		SELECT invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date 
		FROM invoices WHERE order_id = $1`
	
	return r.scanInvoice(ctx, query, orderID)
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID uint64) ([]Invoice, error) {
	if err := r.ensureInvoicesTable(ctx); err != nil {
		return nil, err
	}

	query := `
		SELECT invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date 
		FROM invoices WHERE user_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		invoice, err := r.scanInvoiceFromRows(rows)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, *invoice)
	}
	return invoices, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, invoiceID uint64, status string) error {
	if err := r.ensureInvoicesTable(ctx); err != nil {
		return err
	}

	query := `UPDATE invoices SET status = $1 WHERE invoice_id = $2`
	_, err := r.db.ExecContext(ctx, query, status, invoiceID)
	return err
}

// Helper method to scan a single invoice
func (r *PostgresRepository) scanInvoice(ctx context.Context, query string, arg interface{}) (*Invoice, error) {
	row := r.db.QueryRowContext(ctx, query, arg)
	return r.scanInvoiceFromRow(row)
}

// Helper method to scan invoice from a single row
func (r *PostgresRepository) scanInvoiceFromRow(row *sql.Row) (*Invoice, error) {
	var invoice Invoice
	var itemsJSON []byte

	err := row.Scan(
		&invoice.InvoiceID,
		&invoice.OrderID,
		&invoice.UserID,
		&invoice.Username,
		&invoice.InvoiceNumber,
		&itemsJSON,
		&invoice.Subtotal,
		&invoice.Tax,
		&invoice.Total,
		&invoice.Status,
		&invoice.CreatedAt,
		&invoice.DueDate,
	)
	
	if err != nil {
		return nil, err
	}

	// Unmarshal items JSON
	if len(itemsJSON) > 0 {
		if err := json.Unmarshal(itemsJSON, &invoice.Items); err != nil {
			return nil, err
		}
	}

	return &invoice, nil
}

// Helper method to scan invoice from rows
func (r *PostgresRepository) scanInvoiceFromRows(rows *sql.Rows) (*Invoice, error) {
	var invoice Invoice
	var itemsJSON []byte

	err := rows.Scan(
		&invoice.InvoiceID,
		&invoice.OrderID,
		&invoice.UserID,
		&invoice.Username,
		&invoice.InvoiceNumber,
		&itemsJSON,
		&invoice.Subtotal,
		&invoice.Tax,
		&invoice.Total,
		&invoice.Status,
		&invoice.CreatedAt,
		&invoice.DueDate,
	)
	
	if err != nil {
		return nil, err
	}

	// Unmarshal items JSON
	if len(itemsJSON) > 0 {
		if err := json.Unmarshal(itemsJSON, &invoice.Items); err != nil {
			return nil, err
		}
	}

	return &invoice, nil
}