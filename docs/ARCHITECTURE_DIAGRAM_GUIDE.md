# Architecture Diagram Guide

This document describes the current architecture of the Go Auth Sessions application and provides guidance for updating the DrawIO diagrams.

## 🏗️ Current Architecture (Updated)

### System Overview
```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Router                              │
│                     (main.go)                                   │  
├─────────────────────────────────────────────────────────────────┤
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐   │
│  │    User    │ │  Product   │ │   Order    │ │  Invoice   │   │
│  │  Handler   │ │  Handler   │ │  Handler   │ │  Handler   │   │
│  │            │ │            │ │            │ │            │   │
│  │ ┌────────┐ │ │ ┌────────┐ │ │ ┌────────┐ │ │ ┌────────┐ │   │
│  │ │   JWT  │ │ │ │   JWT  │ │ │ │   JWT  │ │ │ │   JWT  │ │   │
│  │ │  Auth  │ │ │ │  Auth  │ │ │ │  Auth  │ │ │ │  Auth  │ │   │
│  │ └────────┘ │ │ └────────┘ │ │ └────────┘ │ │ └────────┘ │   │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐   │
│  │    User    │ │  Product   │ │   Order    │ │  Invoice   │   │
│  │  Service   │ │  Service   │ │  Service   │ │  Service   │   │
│  │            │ │            │ │     │      │ │            │   │
│  │            │ │            │ │     │      │ │            │   │
│  │            │ │            │ │   Depends  │ │  Depends   │   │
│  │            │ │            │ │     on     │ │    on      │   │
│  │            │ │            │ │  Product   │ │ Order +    │   │
│  │            │ │            │ │            │ │ Product +  │   │
│  │            │ │            │ │            │ │ User       │   │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐   │
│  │    User    │ │  Product   │ │   Order    │ │  Invoice   │   │
│  │Repository  │ │Repository  │ │Repository  │ │Repository  │   │
│  │            │ │            │ │            │ │            │   │
│  │  ┌──────┐  │ │  ┌──────┐  │ │  ┌──────┐  │ │  ┌──────┐  │   │
│  │  │ PG   │  │ │  │ PG   │  │ │  │ PG   │  │ │  │ PG   │  │   │
│  │  └──────┘  │ │  └──────┘  │ │  └──────┘  │ │  └──────┘  │   │
│  │  ┌──────┐  │ │  ┌──────┐  │ │  ┌──────┐  │ │  ┌──────┐  │   │
│  │  │ CH   │  │ │  │ CH   │  │ │  │ CH   │  │ │  │ CH   │  │   │
│  │  └──────┘  │ │  └──────┘  │ │  └──────┘  │ │  └──────┘  │   │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│                   Database Layer                                │
│    ┌─────────────────┐         ┌─────────────────┐             │
│    │   PostgreSQL    │         │   ClickHouse    │             │
│    │                 │         │                 │             │
│    │ ┌─────────────┐ │         │ ┌─────────────┐ │             │
│    │ │    users    │ │         │ │    users    │ │             │
│    │ │  products   │ │         │ │  products   │ │             │
│    │ │   orders    │ │         │ │   orders    │ │             │
│    │ │  invoices   │ │         │ │  invoices   │ │             │
│    │ └─────────────┘ │         │ └─────────────┘ │             │
│    └─────────────────┘         └─────────────────┘             │
└─────────────────────────────────────────────────────────────────┘
```

## 📊 Complete Service Dependencies (Current)

### Service Dependency Graph
```
┌─────────────┐    ┌─────────────┐
│    User     │    │   Product   │
│   Service   │    │   Service   │
│ • Auth      │    │ • Catalog   │
│ • Profiles  │    │ • Stock Mgmt│
└─────────────┘    └─────────────┘
       │                   │
       │                   │
       └───────┬───────────┘
               │
               ▼
    ┌─────────────────┐
    │     Order       │
    │    Service      │
    │  • Validates    │
    │    Products     │
    │  • Calculates   │
    │    Totals       │
    └─────────────────┘
               │
               ▼
    ┌─────────────────┐
    │    Invoice      │
    │    Service      │
    │  • Aggregates   │
    │    Order +      │  
    │    Product +    │
    │    User Data    │
    └─────────────────┘
```

### API Endpoints Overview (Current)

**Public Endpoints (No JWT required):**
```
GET  /health                     - System health check
POST /register                   - User registration  
POST /login                      - User authentication
GET  /products                   - List all products
GET  /products/{id}              - Get product by ID
GET  /products/category/{name}   - Get products by category
```

**Protected Endpoints (JWT required):**
```
POST /products                   - Create product
PUT  /products/{id}              - Update product stock

POST /orders                     - Create order (multiple products)
POST /orders/single              - Create order (single product)  
GET  /orders                     - Get user's orders
GET  /orders/{username}          - Get orders by username (legacy)

POST /invoices                   - Generate invoice from order
GET  /invoices/{id}              - Get invoice by ID
GET  /invoices/order/{id}        - Get invoice by order ID
GET  /invoices/user/{id}         - Get invoices by user ID
PUT  /invoices/{id}              - Update invoice status
```
               ▼
    ┌─────────────────┐
    │    Invoice      │
    │    Service      │
    │  (depends on    │
    │Order+Product+   │
    │     User)       │
    └─────────────────┘
```

### Data Flow
```
1. User Registration/Login
   User → Auth → JWT Token

2. Product Management
   User → Product Service → Database
   
3. Order Creation
   User → Order Service → validates Product → creates Order
   
4. Invoice Generation  
   User → Invoice Service → gets Order → gets Product → gets User → creates Invoice
```

## 🔐 Security Architecture

### JWT Authentication Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │    Auth     │    │  Protected  │
│             │    │  Service    │    │   Service   │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       │ 1. Login          │                   │
       ├──────────────────►│                   │
       │                   │                   │
       │ 2. JWT Token      │                   │
       │◄──────────────────┤                   │
       │                   │                   │
       │ 3. API Request    │                   │
       │   + JWT Token     │                   │
       ├───────────────────┼──────────────────►│
       │                   │                   │
       │                   │ 4. Validate JWT   │
       │                   │◄──────────────────┤
       │                   │                   │
       │                   │ 5. JWT Valid      │
       │                   ├──────────────────►│
       │                   │                   │
       │ 6. Response       │                   │
       │◄──────────────────┼───────────────────┤
       │                   │                   │
```

## 🗄️ Database Schema Relationships

### Entity Relationship Diagram
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    users    │    │  products   │    │   orders    │    │  invoices   │
├─────────────┤    ├─────────────┤    ├─────────────┤    ├─────────────┤
│ user_id (PK)│    │product_id(PK│    │order_id (PK)│    │invoice_id(PK│
│ username    │    │ name        │    │ user_id (FK)│    │ order_id(FK)│
│password_hash│    │ description │    │product_id(FK│    │ user_id (FK)│
│ created_at  │    │ price       │    │ quantity    │    │invoice_num  │
└─────────────┘    │ category    │    │ created_at  │    │ items (JSON)│
       │           │ in_stock    │    └─────────────┘    │ subtotal    │
       │           │ created_at  │           │           │ tax         │
       │           └─────────────┘           │           │ total       │
       │                  │                 │           │ status      │
       │                  │                 │           │ created_at  │
       └──────────────────┼─────────────────┼───────────│ due_date    │
                          │                 │           └─────────────┘
                          └─────────────────┘
```

### Foreign Key Relationships
- `orders.user_id` → `users.user_id`
- `orders.product_id` → `products.product_id`  
- `invoices.order_id` → `orders.order_id`
- `invoices.user_id` → `users.user_id`

## 📚 DrawIO Update Instructions

To update the existing DrawIO diagrams (`go-auth-sessions-workflow.drawio` and `go-auth-sessions-clean-architecture.drawio`):

### 1. Service Layer Updates
- Add **Product Service** box between User and Order
- Add **Invoice Service** box after Order Service
- Update dependency arrows to show:
  - Order → Product (validation)
  - Invoice → Order + Product + User (data aggregation)

### 2. Handler Layer Updates
- Add **Product Handler** with JWT middleware
- Add **Invoice Handler** with JWT middleware
- Show all handlers connecting to the same JWT authentication component

### 3. Repository Layer Updates
- Add **Product Repository** with PostgreSQL and ClickHouse implementations
- Add **Invoice Repository** with PostgreSQL and ClickHouse implementations
- Show repository pattern for all four domains

### 4. Database Schema Updates  
- Update tables to show complete current schema:
  - `users` table (user_id, username, password_hash, created_at)
  - `products` table (product_id, name, description, price, category, in_stock, created_at)
  - `orders` table (order_id, user_id, items_json, subtotal, tax, total, status, created_at)
  - `invoices` table (invoice_id, order_id, user_id, invoice_number, items_json, subtotal, tax, total, status, created_at, due_date)
- Add relationship arrows showing data flow
- Show JSON structure for items_json fields

### 5. API Endpoint Documentation
Update endpoint documentation to show:
- Public endpoints (products listing, health)
- Protected endpoints (orders, invoices, product management)
- JWT requirement indicators

### 6. Data Flow Diagrams
Create/update flow showing:
1. User login → JWT token
2. Product browsing (public)
3. Order creation (protected, validates product)
4. Invoice generation (protected, aggregates data)

## 🎨 Recommended DrawIO Elements

### Colors
- **User Service**: Blue (#4285F4)
- **Product Service**: Green (#34A853) 
- **Order Service**: Orange (#FF9800)
- **Invoice Service**: Purple (#9C27B0)
- **Auth/JWT**: Red (#EA4335)
- **Database**: Gray (#757575)

### Shapes
- **Services**: Rounded rectangles
- **Handlers**: Rectangles with thicker borders
- **Repositories**: Hexagons
- **Databases**: Cylinders
- **Dependencies**: Dashed arrows
- **Data Flow**: Solid arrows
- **JWT Protection**: Shield icons

### Layers
- **Presentation Layer**: HTTP Handlers + Middleware
- **Business Layer**: Services + Domain Logic
- **Data Layer**: Repositories + Database Access
- **Infrastructure**: PostgreSQL + ClickHouse

This comprehensive architecture documentation should provide clear guidance for updating the DrawIO diagrams to reflect the current four-service architecture with proper dependencies and security implementation.
```
POST /register → User Handler → User Service → User Repository → Database
```

### 3. Order Creation Flow
```
POST /orders/{username} → Order Handler → Order Service + User Service → Order Repository + User Repository → Database
```

### 4. Database Switch Flow
```
Environment Variable (DB_DRIVER) → Container → Repository Factory → Specific Implementation
```

## Key Architecture Benefits

1. **Domain Isolation**: Each domain (user, order) is self-contained
2. **Interface Segregation**: Clean interfaces between layers
3. **Dependency Injection**: All dependencies injected via container
4. **Database Agnostic**: Easy switching between database implementations
5. **Testability**: Each layer can be mocked and tested independently
6. **Scalability**: Easy to add new domains and features

## Middleware Stack
```
Request → Logging Middleware → Recovery Middleware → Handler → Response
```

## Container Dependencies
```
Container
├── UserService (depends on UserRepository, PasswordHasher)
├── OrderService (depends on OrderRepository)
├── JWTManager
├── PasswordHasher
└── Database Connection
```

## To Update DrawIO Diagram:

1. Open the existing diagram in Draw.io
2. Replace the monolithic structure with the layered architecture shown above
3. Add separate boxes for each domain (User, Order)
4. Show the dependency injection container
5. Illustrate the database abstraction layer
6. Add arrows showing request flow
7. Include middleware components
8. Save as both .drawio and .png formats in the docs/ folder

The new diagram should emphasize:
- Clean separation of concerns
- Dependency inversion principle
- Domain-driven design
- Database agnostic architecture
- Scalable and maintainable structure

## 📝 Current Architecture Status: **COMPLETE IMPLEMENTATION**

### ✅ Fully Implemented Features
- **4 Complete Domains**: User, Product, Order, Invoice  
- **Dual Database Support**: PostgreSQL + ClickHouse implementations for all domains
- **JWT Authentication**: Consistent middleware across all protected endpoints
- **Clean Architecture**: Repository pattern, service layer, dependency injection
- **Business Logic**: Complete workflow from user registration to invoice generation
- **Testing Suite**: Comprehensive integration testing scripts
- **Documentation**: Complete API documentation and developer guides

### 🏗️ Architecture Highlights
1. **Repository Pattern**: Interface-based database abstraction for all domains
2. **Service Layer**: Business logic with proper dependency injection
3. **Handler Pattern**: HTTP handlers with consistent JWT middleware  
4. **Container Pattern**: Centralized dependency management and service wiring
5. **Environment Configuration**: Database selection via DB_DRIVER environment variable

### 📊 Current API Surface
- **6 Public Endpoints**: Health, auth, product catalog
- **12 Protected Endpoints**: Complete CRUD operations for orders and invoices
- **Multi-format Support**: JSON APIs with comprehensive error handling
- **Cross-service Integration**: Invoice service aggregates data from all other services

This architecture demonstrates **production-ready clean architecture principles** with complete domain separation, comprehensive testing, and enterprise-grade patterns.