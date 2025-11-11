package order

import (
	"context"
	"fmt"
	"time"

	"github.com/rajindersingh041/go-auth-sessions/product"
)

// OrderService defines the business logic interface for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*Order, error)
	CreateSingleOrder(ctx context.Context, userID uint64, req CreateSingleOrderRequest) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error)
	GetOrderByID(ctx context.Context, orderID uint64) (*Order, error)
}

// orderService implements the OrderService interface
type orderService struct {
	repo           OrderRepository
	productService product.ProductService
}

// NewOrderService creates a new order service
func NewOrderService(repo OrderRepository, productService product.ProductService) OrderService {
	return &orderService{
		repo:           repo,
		productService: productService,
	}
}

// CreateOrder creates a new order with multiple products
func (s *orderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*Order, error) {
	if userID == 0 {
		return nil, fmt.Errorf("valid user ID is required")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("at least one product is required")
	}

	var orderItems []OrderItem
	var subtotal float64

	// Validate each product and calculate totals
	for i, item := range req.Items {
		if item.ProductID == 0 || item.Quantity <= 0 {
			return nil, fmt.Errorf("item %d: valid product ID and positive quantity are required", i+1)
		}

		// Validate product exists and is in stock
		prod, err := s.productService.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("item %d: product not found", i+1)
		}
		if !prod.InStock {
			return nil, fmt.Errorf("item %d: product '%s' is currently out of stock", i+1, prod.Name)
		}

		// Calculate item total
		itemTotal := prod.Price * float64(item.Quantity)
		subtotal += itemTotal

		// Create order item
		orderItem := OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: prod.Price,
			Total:     itemTotal,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Calculate tax and total
	tax := subtotal * 0.1 // 10% tax
	total := subtotal + tax

	// Create order
	order := &Order{
		UserID:    userID,
		Items:     orderItems,
		Subtotal:  subtotal,
		Tax:       tax,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// CreateSingleOrder creates an order with a single product (for backward compatibility)
func (s *orderService) CreateSingleOrder(ctx context.Context, userID uint64, req CreateSingleOrderRequest) (*Order, error) {
	// Convert single order to multi-item order
	multiReq := CreateOrderRequest{
		Items: []OrderItemRequest{
			{
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			},
		},
	}
	return s.CreateOrder(ctx, userID, multiReq)
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *orderService) GetOrdersByUserID(ctx context.Context, userID uint64) ([]Order, error) {
	if userID == 0 {
		return nil, fmt.Errorf("valid user ID is required")
	}
	return s.repo.GetOrdersByUserID(ctx, userID)
}

// GetOrderByID retrieves a specific order by ID
func (s *orderService) GetOrderByID(ctx context.Context, orderID uint64) (*Order, error) {
	if orderID == 0 {
		return nil, fmt.Errorf("valid order ID is required")
	}
	return s.repo.GetOrderByID(ctx, orderID)
}

