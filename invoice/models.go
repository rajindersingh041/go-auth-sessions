package invoice

import (
	"context"
)

// Invoice represents an invoice in the database
type Invoice struct {
	InvoiceID     uint64        `json:"invoice_id"`
	OrderID       uint64        `json:"order_id"`
	UserID        uint64        `json:"user_id"`
	Username      string        `json:"username"`
	InvoiceNumber string        `json:"invoice_number"`
	Items         []InvoiceItem `json:"items"`
	Subtotal      float64       `json:"subtotal"`
	Tax           float64       `json:"tax"`
	Total         float64       `json:"total"`
	Status        string        `json:"status"` // "draft", "sent", "paid", "cancelled"
	CreatedAt     string        `json:"created_at"`
	DueDate       string        `json:"due_date"`
}

// InvoiceItem represents an item in an invoice
type InvoiceItem struct {
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Description  string  `json:"description"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
}

// Repository defines the interface for invoice data operations
type Repository interface {
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, invoiceID uint64) (*Invoice, error)
	GetByOrderID(ctx context.Context, orderID uint64) (*Invoice, error)
	GetByUserID(ctx context.Context, userID uint64) ([]Invoice, error)
	UpdateStatus(ctx context.Context, invoiceID uint64, status string) error
}

// CreateInvoiceRequest represents the request to create an invoice
type CreateInvoiceRequest struct {
	OrderID uint64 `json:"order_id"`
}

// UpdateInvoiceStatusRequest represents the request to update invoice status
type UpdateInvoiceStatusRequest struct {
	Status string `json:"status"`
}