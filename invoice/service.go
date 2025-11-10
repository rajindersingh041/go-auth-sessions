package invoice

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/order"
	"github.com/rajindersingh041/go-auth-sessions/product"
	"github.com/rajindersingh041/go-auth-sessions/user"
)

// Service defines the business logic interface for invoice operations
type Service interface {
	CreateInvoiceFromOrder(ctx context.Context, orderID uint64) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, invoiceID uint64) (*Invoice, error)
	GetInvoiceByOrderID(ctx context.Context, orderID uint64) (*Invoice, error)
	GetInvoicesByUserID(ctx context.Context, userID uint64) ([]Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, invoiceID uint64, status string) error
}

// service implements the Service interface
type service struct {
	repo           Repository
	orderService   order.Service
	productService product.Service
	userService    user.Service
}

// NewService creates a new invoice service
func NewService(repo Repository, orderService order.Service, productService product.Service, userService user.Service) Service {
	return &service{
		repo:           repo,
		orderService:   orderService,
		productService: productService,
		userService:    userService,
	}
}

// CreateInvoiceFromOrder creates an invoice from an existing order
func (s *service) CreateInvoiceFromOrder(ctx context.Context, orderID uint64) (*Invoice, error) {
	if orderID == 0 {
		return nil, fmt.Errorf("valid order ID is required")
	}

	// Check if invoice already exists for this order
	existingInvoice, err := s.repo.GetByOrderID(ctx, orderID)
	if err == nil && existingInvoice != nil {
		// If invoice exists, populate it with complete details and return
		return s.populateInvoiceDetails(ctx, existingInvoice)
	}

	// Get order details
	orderDetails, err := s.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order details: %w", err)
	}

	// Get product details
	productDetails, err := s.productService.GetProductByID(ctx, orderDetails.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product details: %w", err)
	}

	// Get user details 
	userDetails, err := s.userService.GetUserByID(ctx, orderDetails.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	// Calculate invoice totals
	itemPrice := productDetails.Price
	quantity := float64(orderDetails.Quantity)
	subtotal := itemPrice * quantity
	tax := subtotal * 0.1 // 10% tax
	total := subtotal + tax

	// Create invoice items
	items := []InvoiceItem{
		{
			ProductID:   productDetails.ProductID,
			ProductName: productDetails.Name,
			Description: productDetails.Description,
			Quantity:    orderDetails.Quantity,
			UnitPrice:   itemPrice,
			TotalPrice:  itemPrice * quantity,
		},
	}

	// Generate invoice number
	invoiceNumber := s.generateInvoiceNumber()
	
	// Create invoice structure
	invoice := &Invoice{
		OrderID:       orderID,
		UserID:        orderDetails.UserID,
		Username:      userDetails.Username,
		InvoiceNumber: invoiceNumber,
		Items:         items,
		Subtotal:      subtotal,
		Tax:           tax,
		Total:         total,
		Status:        "draft",
		CreatedAt:     time.Now().Format(time.RFC3339),
		DueDate:       time.Now().AddDate(0, 0, 30).Format(time.RFC3339), // 30 days from now
	}

	// Create the invoice
	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Return the created invoice (now has ID populated by Create method)
	return invoice, nil
}

// GetInvoiceByID retrieves an invoice by its ID
func (s *service) GetInvoiceByID(ctx context.Context, invoiceID uint64) (*Invoice, error) {
	if invoiceID == 0 {
		return nil, fmt.Errorf("valid invoice ID is required")
	}
	invoice, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	return s.populateInvoiceDetails(ctx, invoice)
}

// GetInvoiceByOrderID retrieves an invoice by order ID
func (s *service) GetInvoiceByOrderID(ctx context.Context, orderID uint64) (*Invoice, error) {
	if orderID == 0 {
		return nil, fmt.Errorf("valid order ID is required")
	}
	invoice, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return s.populateInvoiceDetails(ctx, invoice)
}

// GetInvoicesByUserID retrieves all invoices for a user
func (s *service) GetInvoicesByUserID(ctx context.Context, userID uint64) ([]Invoice, error) {
	if userID == 0 {
		return nil, fmt.Errorf("valid user ID is required")
	}
	invoices, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Populate details for each invoice
	populatedInvoices := make([]Invoice, len(invoices))
	for i, invoice := range invoices {
		populated, err := s.populateInvoiceDetails(ctx, &invoice)
		if err != nil {
			// If we can't populate details, use the original invoice
			populatedInvoices[i] = invoice
		} else {
			populatedInvoices[i] = *populated
		}
	}
	return populatedInvoices, nil
}

// UpdateInvoiceStatus updates the status of an invoice
func (s *service) UpdateInvoiceStatus(ctx context.Context, invoiceID uint64, status string) error {
	if invoiceID == 0 {
		return fmt.Errorf("valid invoice ID is required")
	}
	
	validStatuses := map[string]bool{
		"draft":     true,
		"sent":      true,
		"paid":      true,
		"cancelled": true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s. Valid statuses are: draft, sent, paid, cancelled", status)
	}
	
	return s.repo.UpdateStatus(ctx, invoiceID, status)
}

// generateInvoiceNumber generates a unique invoice number
func (s *service) generateInvoiceNumber() string {
	timestamp := time.Now().Unix()
	return "INV-" + strconv.FormatInt(timestamp, 10)
}

// populateInvoiceDetails populates an invoice with complete order, product, and user details
func (s *service) populateInvoiceDetails(ctx context.Context, invoice *Invoice) (*Invoice, error) {
	// If invoice is already populated (has user info), return as is
	if invoice.Username != "" && len(invoice.Items) > 0 {
		return invoice, nil
	}

	// Get order details
	orderDetails, err := s.orderService.GetOrderByID(ctx, invoice.OrderID)
	if err != nil {
		return invoice, nil // Return original invoice if we can't get order details
	}

	// Get product details
	productDetails, err := s.productService.GetProductByID(ctx, orderDetails.ProductID)
	if err != nil {
		return invoice, nil // Return original invoice if we can't get product details
	}

	// Get user details
	userDetails, err := s.userService.GetUserByID(ctx, orderDetails.UserID)
	if err != nil {
		return invoice, nil // Return original invoice if we can't get user details
	}

	// Populate missing fields
	if invoice.UserID == 0 {
		invoice.UserID = orderDetails.UserID
	}
	if invoice.Username == "" {
		invoice.Username = userDetails.Username
	}
	if len(invoice.Items) == 0 {
		itemPrice := productDetails.Price
		quantity := float64(orderDetails.Quantity)
		invoice.Items = []InvoiceItem{
			{
				ProductID:   productDetails.ProductID,
				ProductName: productDetails.Name,
				Description: productDetails.Description,
				Quantity:    orderDetails.Quantity,
				UnitPrice:   itemPrice,
				TotalPrice:  itemPrice * quantity,
			},
		}
	}
	if invoice.Subtotal == 0 {
		itemPrice := productDetails.Price
		quantity := float64(orderDetails.Quantity)
		invoice.Subtotal = itemPrice * quantity
		invoice.Tax = invoice.Subtotal * 0.1 // 10% tax
		invoice.Total = invoice.Subtotal + invoice.Tax
	}

	return invoice, nil
}