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

// InvoiceService defines the business logic interface for invoice operations
type InvoiceService interface {
	CreateInvoiceFromOrder(ctx context.Context, orderID uint64) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, invoiceID uint64) (*Invoice, error)
	GetInvoiceByOrderID(ctx context.Context, orderID uint64) (*Invoice, error)
	GetInvoicesByUserID(ctx context.Context, userID uint64) ([]Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, invoiceID uint64, status string) error
}

// invoiceService implements the InvoiceService interface
type invoiceService struct {
	repo           InvoiceRepository
	orderService   order.OrderService
	productService product.ProductService
	userService    user.UserService
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(repo InvoiceRepository, orderService order.OrderService, productService product.ProductService, userService user.UserService) InvoiceService {
	return &invoiceService{
		repo:           repo,
		orderService:   orderService,
		productService: productService,
		userService:    userService,
	}
}

// CreateInvoiceFromOrder creates an invoice from an existing order
func (s *invoiceService) CreateInvoiceFromOrder(ctx context.Context, orderID uint64) (*Invoice, error) {
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

	// Get user details 
	userDetails, err := s.userService.GetUserByID(ctx, orderDetails.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	// Create invoice items from order items
	var invoiceItems []InvoiceItem
	for _, orderItem := range orderDetails.Items {
		// Get product details for each item
		productDetails, err := s.productService.GetProductByID(ctx, orderItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product details for product %d: %w", orderItem.ProductID, err)
		}

		invoiceItem := InvoiceItem{
			ProductID:   orderItem.ProductID,
			ProductName: productDetails.Name,
			Description: productDetails.Description,
			Quantity:    orderItem.Quantity,
			UnitPrice:   orderItem.UnitPrice,
			TotalPrice:  orderItem.Total,
		}
		invoiceItems = append(invoiceItems, invoiceItem)
	}

	// Use order totals (already calculated)
	subtotal := orderDetails.Subtotal
	tax := orderDetails.Tax
	total := orderDetails.Total

	// Generate invoice number
	invoiceNumber := s.generateInvoiceNumber()
	
	// Create invoice structure
	invoice := &Invoice{
		OrderID:       orderID,
		UserID:        orderDetails.UserID,
		Username:      userDetails.Username,
		InvoiceNumber: invoiceNumber,
		Items:         invoiceItems,
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
func (s *invoiceService) GetInvoiceByID(ctx context.Context, invoiceID uint64) (*Invoice, error) {
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
func (s *invoiceService) GetInvoiceByOrderID(ctx context.Context, orderID uint64) (*Invoice, error) {
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
func (s *invoiceService) GetInvoicesByUserID(ctx context.Context, userID uint64) ([]Invoice, error) {
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
func (s *invoiceService) UpdateInvoiceStatus(ctx context.Context, invoiceID uint64, status string) error {
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
func (s *invoiceService) generateInvoiceNumber() string {
	timestamp := time.Now().Unix()
	return "INV-" + strconv.FormatInt(timestamp, 10)
}

// populateInvoiceDetails populates an invoice with complete order, product, and user details
func (s *invoiceService) populateInvoiceDetails(ctx context.Context, invoice *Invoice) (*Invoice, error) {
	// If invoice is already populated (has user info), return as is
	if invoice.Username != "" && len(invoice.Items) > 0 {
		return invoice, nil
	}

	// Get order details
	orderDetails, err := s.orderService.GetOrderByID(ctx, invoice.OrderID)
	if err != nil {
		return invoice, nil // Return original invoice if we can't get order details
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
	
	// Populate invoice items from order items
	if len(invoice.Items) == 0 {
		var invoiceItems []InvoiceItem
		
		for _, orderItem := range orderDetails.Items {
			// Get product details for each item
			productDetails, err := s.productService.GetProductByID(ctx, orderItem.ProductID)
			if err != nil {
				continue // Skip this item if we can't get product details
			}

			invoiceItem := InvoiceItem{
				ProductID:   orderItem.ProductID,
				ProductName: productDetails.Name,
				Description: productDetails.Description,
				Quantity:    orderItem.Quantity,
				UnitPrice:   orderItem.UnitPrice,
				TotalPrice:  orderItem.Total,
			}
			invoiceItems = append(invoiceItems, invoiceItem)
		}
		invoice.Items = invoiceItems
	}
	
	// Use order totals if invoice totals are not set
	if invoice.Subtotal == 0 {
		invoice.Subtotal = orderDetails.Subtotal
		invoice.Tax = orderDetails.Tax
		invoice.Total = orderDetails.Total
	}

	return invoice, nil
}