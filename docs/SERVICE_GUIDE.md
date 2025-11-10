# Adding New Services - Developer Guide

This guide explains how to add new business domains/services to the Go Auth Sessions application following the established clean architecture patterns.

## 📋 Overview

The application follows a **domain-driven design** with each service being completely self-contained. To add a new service, you'll need to create the complete domain structure and wire it into the dependency injection system.

## 🚀 Step-by-Step Guide

### Step 1: Create Domain Directory Structure

Create a new directory for your domain (e.g., `shipping/`):

```bash
mkdir shipping
```

### Step 2: Create Domain Models (`models.go`)

Define your entities, DTOs, and repository interface:

```go
package shipping

import "context"

// Shipping represents a shipping record in the database
type Shipping struct {
    ShippingID  uint64  `json:"shipping_id"`
    OrderID     uint64  `json:"order_id"`
    Address     string  `json:"address"`
    TrackingID  string  `json:"tracking_id"`
    Status      string  `json:"status"` // "pending", "shipped", "delivered"
    Cost        float64 `json:"cost"`
    CreatedAt   string  `json:"created_at"`
    DeliveredAt string  `json:"delivered_at,omitempty"`
}

// Repository defines the interface for shipping data operations
type Repository interface {
    Create(ctx context.Context, shipping *Shipping) error
    GetByID(ctx context.Context, shippingID uint64) (*Shipping, error)
    GetByOrderID(ctx context.Context, orderID uint64) (*Shipping, error)
    UpdateStatus(ctx context.Context, shippingID uint64, status string) error
}

// CreateShippingRequest represents the request to create a shipping record
type CreateShippingRequest struct {
    OrderID uint64 `json:"order_id"`
    Address string `json:"address"`
}

// UpdateShippingStatusRequest represents the request to update shipping status
type UpdateShippingStatusRequest struct {
    Status string `json:"status"`
}
```

### Step 3: Create Business Service (`service.go`)

Implement your business logic and service interface:

```go
package shipping

import (
    "context"
    "fmt"
    "time"
    "math/rand"
    "strconv"

    "github.com/rajindersingh041/go-auth-sessions/order"
)

// Service defines the business logic interface for shipping operations
type Service interface {
    CreateShipping(ctx context.Context, req CreateShippingRequest) (*Shipping, error)
    GetShippingByID(ctx context.Context, shippingID uint64) (*Shipping, error)
    GetShippingByOrderID(ctx context.Context, orderID uint64) (*Shipping, error)
    UpdateShippingStatus(ctx context.Context, shippingID uint64, status string) error
}

// service implements the Service interface
type service struct {
    repo         Repository
    orderService order.Service // Dependency on order service
}

// NewService creates a new shipping service
func NewService(repo Repository, orderService order.Service) Service {
    return &service{
        repo:         repo,
        orderService: orderService,
    }
}

// CreateShipping creates a new shipping record
func (s *service) CreateShipping(ctx context.Context, req CreateShippingRequest) (*Shipping, error) {
    // Validate input
    if req.OrderID == 0 || req.Address == "" {
        return nil, fmt.Errorf("order ID and address are required")
    }

    // Validate order exists
    order, err := s.orderService.GetOrderByID(ctx, req.OrderID)
    if err != nil {
        return nil, fmt.Errorf("order not found")
    }

    // Generate tracking ID
    trackingID := s.generateTrackingID()

    // Calculate shipping cost (business logic)
    cost := s.calculateShippingCost(order.Quantity)

    // Create shipping record
    shipping := &Shipping{
        OrderID:    req.OrderID,
        Address:    req.Address,
        TrackingID: trackingID,
        Status:     "pending",
        Cost:       cost,
        CreatedAt:  time.Now().Format(time.RFC3339),
    }

    if err := s.repo.Create(ctx, shipping); err != nil {
        return nil, fmt.Errorf("failed to create shipping record: %w", err)
    }

    return shipping, nil
}

// GetShippingByID retrieves a shipping record by ID
func (s *service) GetShippingByID(ctx context.Context, shippingID uint64) (*Shipping, error) {
    if shippingID == 0 {
        return nil, fmt.Errorf("valid shipping ID is required")
    }
    return s.repo.GetByID(ctx, shippingID)
}

// GetShippingByOrderID retrieves a shipping record by order ID
func (s *service) GetShippingByOrderID(ctx context.Context, orderID uint64) (*Shipping, error) {
    if orderID == 0 {
        return nil, fmt.Errorf("valid order ID is required")
    }
    return s.repo.GetByOrderID(ctx, orderID)
}

// UpdateShippingStatus updates the status of a shipping record
func (s *service) UpdateShippingStatus(ctx context.Context, shippingID uint64, status string) error {
    validStatuses := map[string]bool{
        "pending":   true,
        "shipped":   true,
        "delivered": true,
        "cancelled": true,
    }
    
    if !validStatuses[status] {
        return fmt.Errorf("invalid status: %s", status)
    }
    
    return s.repo.UpdateStatus(ctx, shippingID, status)
}

// Business logic helpers
func (s *service) generateTrackingID() string {
    return "SHIP-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(rand.Intn(1000))
}

func (s *service) calculateShippingCost(quantity int) float64 {
    baseCost := 5.99
    return baseCost + (float64(quantity) * 2.50)
}
```

### Step 4: Create Repository Implementations

#### PostgreSQL Repository (`repository_postgres.go`)

```go
package shipping

import (
    "context"
    "database/sql"
)

// PostgresRepository implements Repository for PostgreSQL database
type PostgresRepository struct {
    db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL shipping repository
func NewPostgresRepository(db *sql.DB) Repository {
    return &PostgresRepository{db: db}
}

// ensureShippingTable creates the shipping table if it doesn't exist
func (r *PostgresRepository) ensureShippingTable(ctx context.Context) error {
    query := `
        CREATE TABLE IF NOT EXISTS shipping (
            shipping_id SERIAL PRIMARY KEY,
            order_id BIGINT NOT NULL,
            address TEXT NOT NULL,
            tracking_id TEXT NOT NULL UNIQUE,
            status TEXT NOT NULL DEFAULT 'pending',
            cost DECIMAL(10,2) NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            delivered_at TIMESTAMP
        )`
    _, err := r.db.ExecContext(ctx, query)
    return err
}

func (r *PostgresRepository) Create(ctx context.Context, shipping *Shipping) error {
    if err := r.ensureShippingTable(ctx); err != nil {
        return err
    }
    
    query := `INSERT INTO shipping (order_id, address, tracking_id, status, cost, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING shipping_id`
    err := r.db.QueryRowContext(ctx, query, 
        shipping.OrderID, 
        shipping.Address,
        shipping.TrackingID,
        shipping.Status,
        shipping.Cost,
        shipping.CreatedAt).Scan(&shipping.ShippingID)
    return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, shippingID uint64) (*Shipping, error) {
    if err := r.ensureShippingTable(ctx); err != nil {
        return nil, err
    }
    
    query := `SELECT shipping_id, order_id, address, tracking_id, status, cost, created_at, delivered_at 
              FROM shipping WHERE shipping_id = $1`
    row := r.db.QueryRowContext(ctx, query, shippingID)
    
    var shipping Shipping
    var deliveredAt sql.NullString
    err := row.Scan(
        &shipping.ShippingID,
        &shipping.OrderID,
        &shipping.Address,
        &shipping.TrackingID,
        &shipping.Status,
        &shipping.Cost,
        &shipping.CreatedAt,
        &deliveredAt,
    )
    
    if deliveredAt.Valid {
        shipping.DeliveredAt = deliveredAt.String
    }
    
    return &shipping, err
}

func (r *PostgresRepository) GetByOrderID(ctx context.Context, orderID uint64) (*Shipping, error) {
    if err := r.ensureShippingTable(ctx); err != nil {
        return nil, err
    }
    
    query := `SELECT shipping_id, order_id, address, tracking_id, status, cost, created_at, delivered_at 
              FROM shipping WHERE order_id = $1`
    row := r.db.QueryRowContext(ctx, query, orderID)
    
    var shipping Shipping
    var deliveredAt sql.NullString
    err := row.Scan(
        &shipping.ShippingID,
        &shipping.OrderID,
        &shipping.Address,
        &shipping.TrackingID,
        &shipping.Status,
        &shipping.Cost,
        &shipping.CreatedAt,
        &deliveredAt,
    )
    
    if deliveredAt.Valid {
        shipping.DeliveredAt = deliveredAt.String
    }
    
    return &shipping, err
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, shippingID uint64, status string) error {
    if err := r.ensureShippingTable(ctx); err != nil {
        return err
    }
    
    var query string
    if status == "delivered" {
        query = `UPDATE shipping SET status = $1, delivered_at = NOW() WHERE shipping_id = $2`
    } else {
        query = `UPDATE shipping SET status = $1 WHERE shipping_id = $2`
    }
    
    _, err := r.db.ExecContext(ctx, query, status, shippingID)
    return err
}
```

#### ClickHouse Repository (`repository_clickhouse.go`)

```go
package shipping

import (
    "context"
    "database/sql"
)

// ClickHouseRepository implements Repository for ClickHouse database
type ClickHouseRepository struct {
    db *sql.DB
}

// NewClickHouseRepository creates a new ClickHouse shipping repository
func NewClickHouseRepository(db *sql.DB) Repository {
    return &ClickHouseRepository{db: db}
}

func (r *ClickHouseRepository) Create(ctx context.Context, shipping *Shipping) error {
    // Generate ID for ClickHouse
    var maxID uint64
    err := r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(shipping_id), 0) FROM shipping").Scan(&maxID)
    if err != nil {
        maxID = 0
    }
    
    shipping.ShippingID = maxID + 1
    query := `INSERT INTO shipping (shipping_id, order_id, address, tracking_id, status, cost, created_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
    _, err = r.db.ExecContext(ctx, query, 
        shipping.ShippingID,
        shipping.OrderID,
        shipping.Address,
        shipping.TrackingID,
        shipping.Status,
        shipping.Cost,
        shipping.CreatedAt)
    return err
}

func (r *ClickHouseRepository) GetByID(ctx context.Context, shippingID uint64) (*Shipping, error) {
    query := `SELECT shipping_id, order_id, address, tracking_id, status, cost, created_at, delivered_at 
              FROM shipping WHERE shipping_id = ?`
    row := r.db.QueryRowContext(ctx, query, shippingID)
    
    var shipping Shipping
    err := row.Scan(
        &shipping.ShippingID,
        &shipping.OrderID,
        &shipping.Address,
        &shipping.TrackingID,
        &shipping.Status,
        &shipping.Cost,
        &shipping.CreatedAt,
        &shipping.DeliveredAt,
    )
    
    return &shipping, err
}

func (r *ClickHouseRepository) GetByOrderID(ctx context.Context, orderID uint64) (*Shipping, error) {
    query := `SELECT shipping_id, order_id, address, tracking_id, status, cost, created_at, delivered_at 
              FROM shipping WHERE order_id = ?`
    row := r.db.QueryRowContext(ctx, query, orderID)
    
    var shipping Shipping
    err := row.Scan(
        &shipping.ShippingID,
        &shipping.OrderID,
        &shipping.Address,
        &shipping.TrackingID,
        &shipping.Status,
        &shipping.Cost,
        &shipping.CreatedAt,
        &shipping.DeliveredAt,
    )
    
    return &shipping, err
}

func (r *ClickHouseRepository) UpdateStatus(ctx context.Context, shippingID uint64, status string) error {
    // Note: ClickHouse UPDATE is limited, consider using ReplacingMergeTree
    var deliveredAt string
    if status == "delivered" {
        deliveredAt = time.Now().Format(time.RFC3339)
    }
    
    query := `ALTER TABLE shipping UPDATE status = ?, delivered_at = ? WHERE shipping_id = ?`
    _, err := r.db.ExecContext(ctx, query, status, deliveredAt, shippingID)
    return err
}
```

### Step 5: Create HTTP Handler (`handler.go`)

```go
package shipping

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "github.com/rajindersingh041/go-auth-sessions/auth"
)

// Handler handles HTTP requests for shipping operations
type Handler struct {
    service    Service
    jwtManager auth.JWTManager
}

// NewHandler creates a new shipping handler
func NewHandler(service Service, jwtManager auth.JWTManager) *Handler {
    return &Handler{
        service:    service,
        jwtManager: jwtManager,
    }
}

// RegisterRoutes registers all shipping-related routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
    // All shipping routes require authentication
    mux.Handle("POST /shipping", h.requireAuth(http.HandlerFunc(h.handleCreateShipping())))
    mux.Handle("GET /shipping/", h.requireAuth(http.HandlerFunc(h.handleGetShipping())))
    mux.Handle("PUT /shipping/", h.requireAuth(http.HandlerFunc(h.handleUpdateShippingStatus())))
}

// requireAuth middleware for JWT authentication
func (h *Handler) requireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            respondError(w, http.StatusUnauthorized, "Authorization header required")
            return
        }

        if !strings.HasPrefix(authHeader, "Bearer ") {
            respondError(w, http.StatusUnauthorized, "Authorization header must start with 'Bearer '")
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == "" {
            respondError(w, http.StatusUnauthorized, "JWT token is required")
            return
        }

        username, err := h.jwtManager.ValidateToken(token)
        if err != nil {
            respondError(w, http.StatusUnauthorized, "Invalid or expired JWT token")
            return
        }

        _ = username // Add to context if needed
        next.ServeHTTP(w, r.WithContext(r.Context()))
    })
}

// handleCreateShipping handles requests to create shipping
func (h *Handler) handleCreateShipping() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req CreateShippingRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondError(w, http.StatusBadRequest, "Invalid request body")
            return
        }

        shipping, err := h.service.CreateShipping(r.Context(), req)
        if err != nil {
            log.Printf("Shipping creation failed: %v", err)
            respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create shipping: %v", err))
            return
        }

        respondJSON(w, http.StatusCreated, shipping)
    }
}

// handleGetShipping handles requests to get shipping by ID or order ID
func (h *Handler) handleGetShipping() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimPrefix(r.URL.Path, "/shipping/")
        if path == "" {
            respondError(w, http.StatusBadRequest, "Shipping ID or order ID required")
            return
        }

        var shipping *Shipping
        var err error

        if strings.HasPrefix(path, "order/") {
            orderIDStr := strings.TrimPrefix(path, "order/")
            orderID, parseErr := strconv.ParseUint(orderIDStr, 10, 64)
            if parseErr != nil {
                respondError(w, http.StatusBadRequest, "Invalid order ID")
                return
            }
            shipping, err = h.service.GetShippingByOrderID(r.Context(), orderID)
        } else {
            shippingID, parseErr := strconv.ParseUint(path, 10, 64)
            if parseErr != nil {
                respondError(w, http.StatusBadRequest, "Invalid shipping ID")
                return
            }
            shipping, err = h.service.GetShippingByID(r.Context(), shippingID)
        }

        if err != nil {
            respondError(w, http.StatusNotFound, "Shipping record not found")
            return
        }

        respondJSON(w, http.StatusOK, shipping)
    }
}

// handleUpdateShippingStatus handles requests to update shipping status
func (h *Handler) handleUpdateShippingStatus() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimPrefix(r.URL.Path, "/shipping/")
        shippingID, err := strconv.ParseUint(path, 10, 64)
        if err != nil {
            respondError(w, http.StatusBadRequest, "Invalid shipping ID")
            return
        }

        var req UpdateShippingStatusRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondError(w, http.StatusBadRequest, "Invalid request body")
            return
        }

        if err := h.service.UpdateShippingStatus(r.Context(), shippingID, req.Status); err != nil {
            respondError(w, http.StatusInternalServerError, "Failed to update shipping status")
            return
        }

        respondJSON(w, http.StatusOK, map[string]string{
            "message": "Shipping status updated successfully",
        })
    }
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{
        "error": message,
    })
}
```

### Step 6: Update Dependency Injection Container

Add your new service to `container.go`:

```go
// 1. Add import
import (
    // ... existing imports
    "github.com/rajindersingh041/go-auth-sessions/shipping"
)

// 2. Add to Container struct
type Container struct {
    // ... existing services
    ShippingService shipping.Service
    // ... existing components
}

// 3. Add repository creation
func NewContainer(db *sql.DB, dbDriver string) *Container {
    // ... existing repository creation
    var shippingRepo shipping.Repository

    switch dbDriver {
    case "clickhouse":
        // ... existing repos
        shippingRepo = shipping.NewClickHouseRepository(db)
    case "postgres":
        // ... existing repos  
        shippingRepo = shipping.NewPostgresRepository(db)
    default:
        log.Fatalf("Unsupported DB_DRIVER: %s", dbDriver)
    }

    // ... existing service creation
    shippingService := shipping.NewService(shippingRepo, orderService)

    return &Container{
        // ... existing services
        ShippingService: shippingService,
        // ... existing components
    }
}
```

### Step 7: Update Main Application

Add your handler to `main.go`:

```go
// 1. Add import
import (
    // ... existing imports
    "github.com/rajindersingh041/go-auth-sessions/shipping"
)

// 2. Create handler and update setupServer
func main() {
    // ... existing code

    // Create HTTP handlers
    // ... existing handlers
    shippingHandler := shipping.NewHandler(container.ShippingService, container.JWTManager)

    // Setup HTTP server with routes
    server := setupServer(userHandler, orderHandler, productHandler, invoiceHandler, shippingHandler)
    
    // ... rest of main
}

// 3. Update setupServer function signature
func setupServer(
    userHandler *user.Handler, 
    orderHandler *order.Handler, 
    productHandler *product.Handler, 
    invoiceHandler *invoice.Handler,
    shippingHandler *shipping.Handler,
) http.Handler {
    mux := http.NewServeMux()

    // Health check endpoint
    mux.HandleFunc("GET /health", handleHealth())

    // Register domain-specific routes
    userHandler.RegisterRoutes(mux)
    orderHandler.RegisterRoutes(mux)
    productHandler.RegisterRoutes(mux)
    invoiceHandler.RegisterRoutes(mux)
    shippingHandler.RegisterRoutes(mux) // Add your new service

    // Apply global middleware
    handler := globalLoggingMiddleware(globalRecoveryMiddleware(mux))
    return handler
}
```

### Step 8: Create Database Schema (Optional)

Add to `MIGRATION.md`:

```sql
-- PostgreSQL
CREATE TABLE IF NOT EXISTS shipping (
    shipping_id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    address TEXT NOT NULL,
    tracking_id TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'pending',
    cost DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMP
);

-- ClickHouse
CREATE TABLE shipping (
    shipping_id UInt64,
    order_id UInt64,
    address String,
    tracking_id String,
    status String,
    cost Decimal64(2),
    created_at String,
    delivered_at String
) ENGINE = MergeTree()
ORDER BY shipping_id;
```

## ✅ Verification Checklist

After implementing your new service, verify:

- [ ] Domain directory created with all required files
- [ ] Models defined with proper entities and DTOs
- [ ] Service interface and implementation created
- [ ] Both PostgreSQL and ClickHouse repositories implemented
- [ ] HTTP handler with JWT authentication middleware
- [ ] Service added to dependency injection container
- [ ] Handler registered in main.go setupServer function
- [ ] Database schema documented
- [ ] API endpoints tested with Postman/curl

## 🔧 Testing Your New Service

```bash
# Test the new shipping endpoints
curl -X POST http://localhost:8080/shipping \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{"order_id":123,"address":"123 Main St, City, State"}'

curl -X GET http://localhost:8080/shipping/order/123 \
  -H "Authorization: Bearer <jwt_token>"
```

## 📚 Best Practices

1. **Keep domains isolated** - No direct imports between business domains
2. **Use interfaces** - All dependencies should be interfaces for testability
3. **Handle errors gracefully** - Provide meaningful error messages
4. **Add authentication** - Use JWT middleware for protected endpoints
5. **Document your API** - Update API documentation with new endpoints
6. **Write tests** - Create unit tests for your service logic
7. **Consider relationships** - If your service depends on others, inject them properly

## 🔄 Next Steps

After adding your service:

1. Write comprehensive unit tests
2. Add integration tests for the HTTP endpoints
3. Update API documentation
4. Consider adding metrics and logging
5. Update the main README.md with your new endpoints
6. Create any necessary database migrations

This pattern ensures consistency across all services and maintains the clean architecture principles of the application.