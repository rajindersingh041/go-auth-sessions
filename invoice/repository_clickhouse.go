package invoice

import (
	"context"
	"database/sql"
	"encoding/json"
)

// ClickHouseRepository implements Repository for ClickHouse database
type ClickHouseRepository struct {
	db *sql.DB
}

// NewClickHouseRepository creates a new ClickHouse invoice repository
func NewClickHouseRepository(db *sql.DB) InvoiceRepository {
	return &ClickHouseRepository{db: db}
}

func (r *ClickHouseRepository) Create(ctx context.Context, invoice *Invoice) error {
	// For ClickHouse, we'll generate a simple incremental ID
	// First get the current max ID
	var maxID uint64
	maxQuery := "SELECT COALESCE(MAX(invoice_id), 0) FROM invoices"
	err := r.db.QueryRowContext(ctx, maxQuery).Scan(&maxID)
	if err != nil {
		maxID = 0 // Start from 0 if no invoices exist
	}
	
	// Set the next ID
	invoice.InvoiceID = maxID + 1

	// Convert items to JSON string
	itemsJSON, err := json.Marshal(invoice.Items)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO invoices (invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err = r.db.ExecContext(ctx, query, 
		invoice.InvoiceID,
		invoice.OrderID, 
		invoice.UserID, 
		invoice.Username,
		invoice.InvoiceNumber, 
		string(itemsJSON),
		invoice.Subtotal,
		invoice.Tax,
		invoice.Total,
		invoice.Status, 
		invoice.CreatedAt, 
		invoice.DueDate)
	return err
}

func (r *ClickHouseRepository) GetByID(ctx context.Context, invoiceID uint64) (*Invoice, error) {
	query := `
		SELECT invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date 
		FROM invoices WHERE invoice_id = ?`
	
	return r.scanInvoice(ctx, query, invoiceID)
}

func (r *ClickHouseRepository) GetByOrderID(ctx context.Context, orderID uint64) (*Invoice, error) {
	query := `
		SELECT invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date 
		FROM invoices WHERE order_id = ?`
	
	return r.scanInvoice(ctx, query, orderID)
}

func (r *ClickHouseRepository) GetByUserID(ctx context.Context, userID uint64) ([]Invoice, error) {
	query := `
		SELECT invoice_id, order_id, user_id, username, invoice_number, items, subtotal, tax, total, status, created_at, due_date 
		FROM invoices WHERE user_id = ? ORDER BY created_at DESC`
	
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

func (r *ClickHouseRepository) UpdateStatus(ctx context.Context, invoiceID uint64, status string) error {
	// Note: ClickHouse doesn't support UPDATE operations on all table engines
	// For production, you might need to use ReplacingMergeTree or insert a new record
	query := `ALTER TABLE invoices UPDATE status = ? WHERE invoice_id = ?`
	_, err := r.db.ExecContext(ctx, query, status, invoiceID)
	return err
}

// Helper method to scan a single invoice
func (r *ClickHouseRepository) scanInvoice(ctx context.Context, query string, arg interface{}) (*Invoice, error) {
	row := r.db.QueryRowContext(ctx, query, arg)
	return r.scanInvoiceFromRow(row)
}

// Helper method to scan invoice from a single row
func (r *ClickHouseRepository) scanInvoiceFromRow(row *sql.Row) (*Invoice, error) {
	var invoice Invoice
	var itemsJSON string

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
	if itemsJSON != "" {
		if err := json.Unmarshal([]byte(itemsJSON), &invoice.Items); err != nil {
			return nil, err
		}
	}

	return &invoice, nil
}

// Helper method to scan invoice from rows
func (r *ClickHouseRepository) scanInvoiceFromRows(rows *sql.Rows) (*Invoice, error) {
	var invoice Invoice
	var itemsJSON string

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
	if itemsJSON != "" {
		if err := json.Unmarshal([]byte(itemsJSON), &invoice.Items); err != nil {
			return nil, err
		}
	}

	return &invoice, nil
}