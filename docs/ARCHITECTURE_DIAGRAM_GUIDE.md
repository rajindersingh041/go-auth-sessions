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

## 📊 Service Dependencies

### Dependency Graph
```
┌─────────────┐    ┌─────────────┐
│    User     │    │   Product   │
│   Service   │    │   Service   │
└─────────────┘    └─────────────┘
       │                   │
       │                   │
       └───────┬───────────┘
               │
               ▼
    ┌─────────────────┐
    │     Order       │
    │    Service      │
    │  (depends on    │
    │ Product + User) │
    └─────────────────┘
               │
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
- Update tables to show:
  - `users` table
  - `products` table (new)
  - `orders` table with `product_id` instead of `item`
  - `invoices` table (new)
- Add foreign key relationship arrows

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