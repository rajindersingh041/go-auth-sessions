# Go Auth Sessions - Clean Architecture

A scalable Go web application demonstrating clean architecture principles with user authentication and order management.

## 🏗️ Architecture Overview

This project follows **Domain-Driven Design (DDD)** and **Clean Architecture** principles:

```
├── user/                    # User domain
│   ├── models.go           # User entities and DTOs
│   ├── service.go          # Business logic layer
│   ├── handler.go          # HTTP transport layer
│   ├── repository_*.go     # Data access layer
├── order/                   # Order domain  
│   ├── models.go           # Order entities and DTOs
│   ├── service.go          # Business logic layer
│   ├── handler.go          # HTTP transport layer
│   ├── repository_*.go     # Data access layer
├── auth/                    # Authentication utilities
│   └── auth.go             # JWT & password hashing
├── container.go            # Dependency injection
├── main.go                 # Application entry point
├── db.go                   # Database initialization
└── middleware.go           # HTTP middleware
```

## 🎯 Key Design Principles

### 1. **Domain Separation**
- Each domain (`user`, `order`) is self-contained
- Clear boundaries between business logic
- Easy to add new domains without affecting others

### 2. **Layered Architecture**
```
┌─────────────────┐
│   HTTP Layer    │ ← Handlers (transport)
├─────────────────┤
│  Service Layer  │ ← Business logic
├─────────────────┤
│Repository Layer │ ← Data access
└─────────────────┘
```

### 3. **Dependency Injection**
- All dependencies are injected via interfaces
- Easy to mock for testing
- Supports multiple database implementations (ClickHouse, PostgreSQL)

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

### Authentication
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

### Orders
```bash
# Create an order for a user
curl -X POST http://localhost:8080/orders/alice \
  -H "Content-Type: application/json" \
  -d '{"item":"Book","quantity":2}'

# Get all orders for a user
curl -X GET http://localhost:8080/orders/alice
```

### Health Check
```bash
curl -X GET http://localhost:8080/health
```

## 🗄️ Database Support

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

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific domain tests
go test ./user/...
go test ./order/...
```

## 📈 Scaling Considerations

### Adding New Domains
```
1. Create new domain folder: /product
2. Add models.go, service.go, handler.go
3. Implement repository interfaces for each database
4. Register routes in main.go
5. Add to dependency injection container
```

### Database Migration
- Each repository handles its own table creation
- Production: Use proper migration tools (golang-migrate, etc.)
- Tables are auto-created in development mode

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