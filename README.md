# Go Auth Sessions - Clean Architecture Microservices

A scalable Go web application demonstrating clean architecture principles with complete business domain separation including user management, product catalog, order processing, and invoice generation.

## 🏗️ Architecture Overview

This project follows **Domain-Driven Design (DDD)** and **Clean Architecture** principles with complete service separation:

```
├── user/                    # User Management Domain
│   ├── models.go           # User entities and DTOs
│   ├── service.go          # User business logic
│   ├── handler.go          # User HTTP endpoints
│   ├── repository_postgres.go
│   └── repository_clickhouse.go
├── product/                 # Product Catalog Domain
│   ├── models.go           # Product entities and DTOs  
│   ├── service.go          # Product business logic
│   ├── handler.go          # Product HTTP endpoints (JWT protected)
│   ├── repository_postgres.go
│   └── repository_clickhouse.go
├── order/                   # Order Management Domain
│   ├── models.go           # Order entities and DTOs
│   ├── service.go          # Order business logic + validation
│   ├── handler.go          # Order HTTP endpoints (JWT protected)
│   ├── repository_postgres.go
│   └── repository_clickhouse.go
├── invoice/                 # Invoice Generation Domain
│   ├── models.go           # Invoice entities and DTOs
│   ├── service.go          # Invoice business logic
│   ├── handler.go          # Invoice HTTP endpoints (JWT protected)
│   ├── repository_postgres.go
│   └── repository_clickhouse.go
├── auth/                    # Authentication & Authorization
│   └── auth.go             # JWT management & password hashing
├── container.go            # Dependency injection container
├── main.go                 # Application entry point & service wiring
├── db.go                   # Database connection management
└── MIGRATION.md            # Database schema & migration guide
```

## 🎯 Key Design Principles

### 1. **Domain Separation**
- Each domain (`user`, `product`, `order`, `invoice`) is completely self-contained
- Clear boundaries between business contexts
- Zero coupling between domains (only through interfaces)
- Easy to extract domains into separate microservices

### 2. **Layered Architecture**
```
┌─────────────────┐
│   HTTP Layer    │ ← Handlers (transport) + JWT middleware
├─────────────────┤
│  Service Layer  │ ← Business logic + domain validation
├─────────────────┤
│Repository Layer │ ← Data access + database abstraction
└─────────────────┘
```

### 3. **Dependency Injection**
- All dependencies injected via interfaces in `container.go`
- Easy to mock for testing and service isolation
- Supports multiple database implementations (ClickHouse, PostgreSQL)
- Service dependencies properly wired (e.g., orders depend on products)

### 4. **Security & Authentication**
- JWT-based authentication with middleware pattern
- Protected endpoints require `Authorization: Bearer <token>`
- Domain-level authorization (each service handles its own auth)
- Password hashing with bcrypt

### 5. **Database Flexibility**
- Multi-database support (PostgreSQL + ClickHouse)
- Repository pattern abstracts database specifics
- Automatic schema migration for PostgreSQL
- Environment-driven database selection

### 4. **Database Agnostic**
- Repository pattern abstracts database specifics  
- Switch between ClickHouse and PostgreSQL via `DB_DRIVER` env var
- Easy to add new database implementations

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- Docker & Docker Compose
- Make (optional)

### Environment Setup
```bash
# Copy and configure environment variables
cp .env.example .env

# Key configurations:
DB_DRIVER=postgres        # or clickhouse
PORT=8080
JWT_SECRET=your-secret-key
```

### Database Setup
```bash
# Start databases (both ClickHouse and PostgreSQL)
docker compose up -d

# Check database status
docker compose ps
```

### Run Application
```bash
# Install dependencies
go mod tidy

# Run the application
go run .
```

## 📡 API Endpoints

### 🔐 Authentication (Public)
```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'

# Login and get JWT token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

### 📦 Products
```bash
# Get all products (Public)
curl -X GET http://localhost:8080/products

# Get product by ID (Public)
curl -X GET http://localhost:8080/products/2

# Create a new product (Protected - JWT required)
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{"name":"Laptop","description":"Gaming laptop","price":1299.99,"category":"Electronics"}'

# Update product stock (Protected - JWT required)
curl -X PUT http://localhost:8080/products/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{"in_stock":false}'
```

### 🛒 Orders (Protected - JWT required)
```bash
# Create an order (uses product_id, not item name)
curl -X POST http://localhost:8080/orders/alice \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{"product_id":2,"quantity":1}'

# Get all orders for a user
curl -X GET http://localhost:8080/orders/alice \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 🧾 Invoices (Protected - JWT required)
```bash
# Create invoice from order
curl -X POST http://localhost:8080/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{"order_id":123}'

# Get invoice by ID
curl -X GET http://localhost:8080/invoices/456 \
  -H "Authorization: Bearer <your_jwt_token>"

# Get invoice by order ID
curl -X GET http://localhost:8080/invoices/order/123 \
  -H "Authorization: Bearer <your_jwt_token>"

# Get all invoices for a user
curl -X GET http://localhost:8080/invoices/user/1 \
  -H "Authorization: Bearer <your_jwt_token>"

# Update invoice status
curl -X PUT http://localhost:8080/invoices/456 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{"status":"paid"}'
```

### ⚡ Health Check
```bash
curl -X GET http://localhost:8080/health
```

## � Service Architecture & Data Flow

### Service Dependencies
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │   Product   │    │    Order    │    │   Invoice   │
│   Service   │    │   Service   │    │   Service   │    │   Service   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │                   │
       │                   │                   │                   │
       └───────────────────┼───────────────────┼───────────────────┘
                           │                   │
                           │                   │
              ┌────────────▼───────────────────▼────────────┐
              │           Service Dependencies              │
              │  • Orders validate against Products        │
              │  • Invoices need Orders, Products, Users   │
              │  • All services use JWT from Auth          │
              └────────────────────────────────────────────┘
```

### Business Flow
1. **User Registration/Login** → JWT token generation
2. **Product Catalog** → Sample products seeded automatically  
3. **Order Creation** → Validates product exists and is in stock
4. **Invoice Generation** → Creates detailed invoice from order
5. **Order/Invoice Management** → Status tracking and updates

### Database Schema Relationships
```sql
users (user_id, username, password_hash)
  ↓
orders (order_id, user_id, product_id, quantity) 
  ↓                            ↑
invoices (invoice_id, order_id, user_id, items_json)
                                ↑
products (product_id, name, price, in_stock)
```

## �🗄️ Database Support

### ClickHouse Configuration
```env
DB_DRIVER=clickhouse
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=MyPassword2025
```

### PostgreSQL Configuration  
```env
DB_DRIVER=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=mysecretpassword
POSTGRES_DB=authdb
```

### Auto-Migration Features
- **PostgreSQL**: Automatic schema detection and migration
- **Sample Data**: Products are automatically seeded on startup
- **Schema Updates**: Old `item` columns automatically converted to `product_id`

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific domain tests
go test ./user/...
go test ./order/...
go test ./product/...
go test ./invoice/...
```
2. Add models.go, service.go, handler.go  
3. Implement repository interfaces for each database
4. Register routes in main.go
5. Add to dependency injection container
```

See **[SERVICE_GUIDE.md](SERVICE_GUIDE.md)** for detailed step-by-step instructions.

## 🏛️ Complete Service Architecture

### Current Services Overview
```
🔐 Authentication Service (Public)
   ├── POST /register     - User registration
   └── POST /login        - JWT token generation

📦 Product Catalog Service (Mixed Access)  
   ├── GET  /products     - List all products (Public)
   ├── GET  /products/{id} - Get product details (Public)
   ├── POST /products     - Create product (Protected)
   └── PUT  /products/{id} - Update stock (Protected)

🛒 Order Management Service (Protected)
   ├── POST /orders/{username} - Create order with product validation
   └── GET  /orders/{username} - Get user orders

🧾 Invoice Generation Service (Protected)
   ├── POST /invoices           - Generate invoice from order
   ├── GET  /invoices/{id}      - Get invoice by ID
   ├── GET  /invoices/order/{id} - Get invoice by order ID
   ├── GET  /invoices/user/{id}  - Get all user invoices
   └── PUT  /invoices/{id}      - Update invoice status

⚡ System Health
   └── GET  /health       - Health check endpoint
```

### Service Interaction Flow
```
1. User Registration/Login → JWT Token
2. Browse Products (Public) → Product Catalog  
3. Create Order (Protected) → Order Service validates against Product Service
4. Generate Invoice (Protected) → Invoice Service pulls data from Order + Product + User
5. Track Status → Invoice status updates
```

### Database Migration
- Each repository handles its own table creation
- Production: Use proper migration tools (golang-migrate, etc.)
- Tables are auto-created in development mode
- See [MIGRATION.md](MIGRATION.md) for schema details

### Monitoring & Observability
- Structured logging with context
- HTTP request/response logging middleware  
- Easy to add metrics (Prometheus) and tracing (Jaeger)

## 🛠️ Development Guidelines

### Code Organization
- **Models**: Data structures and DTOs
- **Services**: Business logic and validation
- **Handlers**: HTTP request/response handling
- **Repositories**: Database operations

### Adding New Features
1. Start with domain models
2. Implement business logic in service layer
3. Add repository methods if needed
4. Create HTTP handlers
5. Register routes
6. Add tests

### Error Handling
- Domain-specific errors in service layer
- HTTP status codes in handler layer
- Structured error responses with timestamps

---

### How to Use Each File

- **main.go**: Run this file to start the application. It wires up the server and database.
- **server.go**: Used by `main.go` to create and configure the HTTP server.
- **handlers.go**: Called by the server to process incoming HTTP requests.
- **middleware.go**: Used by the server to wrap handlers with cross-cutting concerns (logging, error recovery, authentication).
- **response.go**: Used by handlers and middleware to send consistent JSON responses.
- **db.go**: Called by `main.go` to initialize the database connection.
- **models.go**: Imported wherever user data structures are needed.
- **repository.go**: Used by handlers to interact with the database via the repository pattern.
- **jwt_simple.go**: Used by authentication handlers and middleware for token management.