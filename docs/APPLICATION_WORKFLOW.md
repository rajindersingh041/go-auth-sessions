# Application Workflow Diagram

## Request Flow Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                   CLIENT                                        │
│                        (Web Browser, Mobile App, API Client)                   │
└─────────────────────────────┬───────────────────────────────────────────────────┘
                              │
                              ▼ HTTP Requests
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              MIDDLEWARE STACK                                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────────┐  │
│  │ Logging         │  │ Recovery        │  │ CORS & Security Headers         │  │
│  │ Middleware      │  │ Middleware      │  │ (Future Enhancement)            │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────────────┘  │
└─────────────────────────────┬───────────────────────────────────────────────────┘
                              │
                              ▼ Route to Handler
┌─────────────────────────────────────────────────────────────────────────────────┐
│                               HTTP HANDLERS                                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │ User Handler    │  │ Product Handler │  │ Order Handler   │  │ Invoice     │  │
│  │ • /register     │  │ • /products     │  │ • /orders       │  │ Handler     │  │
│  │ • /login        │  │ (Public & JWT)  │  │ (JWT required)  │  │ • /invoices │  │
│  │ (Public)        │  │                 │  │                 │  │ (JWT req.)  │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘  │
│                                           │                                      │
│  ┌─────────────────────────────────────────────────────────────────────────────┐  │
│  │                      Health Handler - /health                              │  │
│  └─────────────────────────────────────────────────────────────────────────────┘  │
└──────────┬──────────────────┬──────────────────┬──────────────────┬─────────────┘
           │                  │                                         │
           ▼                  ▼                                         ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           DEPENDENCY INJECTION CONTAINER                       │
│  ┌─────────────────────────────────────────────────────────────────────────────┐│
│  │ Services: UserService • ProductService • OrderService • InvoiceService    ││
│  │ Repositories: User • Product • Order • Invoice (PostgreSQL & ClickHouse)  ││
│  │ Auth Components: JWTManager • PasswordHasher                              ││
│  │ Database: Connection Pool with Driver Selection                            ││
│  └─────────────────────────────────────────────────────────────────────────────┘│
└─────────────────────────────┬───────────────────────────────────────────────────┘
                              │
                              ▼ Inject Dependencies
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              SERVICE LAYER                                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │ User Service    │  │ Product Service │  │ Order Service   │  │ Invoice     │  │
│  │ • CreateUser()  │  │ • GetProducts() │  │ • CreateOrder() │  │ Service     │  │
│  │ • Authenticate()│  │ • CreateProd()  │  │ • GetOrders()   │  │ • Generate  │  │
│  │ • Validation    │  │ • UpdateStock() │  │ • Validation    │  │ • Track     │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘  │
│                                           │         ↑                   ↑        │
│  ┌─────────────────────────────────────────────────────────────────────────────┐  │
│  │                      Auth Components                                       │  │
│  │  • JWT Generation & Validation  • Password Hashing  • Token Middleware    │  │
│  └─────────────────────────────────────────────────────────────────────────────┘  │
└──────────┬──────────────────┬──────────────────┬──────────────────────────────────┘
           │                  │
           ▼                  ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            REPOSITORY LAYER                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐│
│  │ User Repo   │  │ Product     │  │ Order Repo  │  │ Invoice Repository      ││
│  │ Interface   │  │ Repository  │  │ Interface   │  │ Interface               ││
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────────┘│
│         │                │                │                       │            │
│         ▼                ▼                ▼                       ▼            │
│  ┌─────────────────────────────────────────────────────────────────────────────┐│
│  │              Database Implementation Selection                              ││
│  │  ┌─────────────────────────────────────────────────────────────────────────┐││
│  │  │ PostgreSQL Implementations  │  ClickHouse Implementations              │││
│  │  │ • user/repository_postgres   │  • user/repository_clickhouse           │││
│  │  │ • product/repository_postgres│  • product/repository_clickhouse        │││
│  │  │ • order/repository_postgres  │  • order/repository_clickhouse          │││
│  │  │ • invoice/repository_postgres│  • invoice/repository_clickhouse        │││
│  │  └─────────────────────────────────────────────────────────────────────────┘││
│  │                     Selected via DB_DRIVER environment variable             ││
│  └─────────────────────────────────────────────────────────────────────────────┘│
└──────────┬──────────────────┬──────────────────┬──────────────────────────────────┘
           │                  │
           ▼                  ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              DATABASE LAYER                                    │
│  ┌─────────────────────────────────┐  ┌─────────────────────────────────────┐  │
│  │        ClickHouse               │  │         PostgreSQL                  │  │
│  │  ┌─────────────────────────────┐│  │  ┌─────────────────────────────────┐│  │
│  │  │ users table                 ││  │  │ users table                     ││  │
│  │  │ - user_id (UInt64)          ││  │  │ - user_id (SERIAL)              ││  │
│  │  │ - username (String)         ││  │  │ - username (TEXT UNIQUE)        ││  │
│  │  │ - password_hash (String)    ││  │  │ - password_hash (TEXT)          ││  │
│  │  │ - created_at (String)       ││  │  │ - created_at (TIMESTAMP)        ││  │
│  │  └─────────────────────────────┘│  │  └─────────────────────────────────┘│  │
│  │  ┌─────────────────────────────┐│  │  ┌─────────────────────────────────┐│  │
│  │  │ products table              ││  │  │ products table                  ││  │
│  │  │ - product_id (UInt64)       ││  │  │ - product_id (SERIAL)           ││  │
│  │  │ - name (String)             ││  │  │ - name (TEXT)                   ││  │
│  │  │ - description (String)      ││  │  │ - description (TEXT)            ││  │
│  │  │ - price (Float64)           ││  │  │ - price (DECIMAL)               ││  │
│  │  │ - category (String)         ││  │  │ - category (TEXT)               ││  │
│  │  │ - in_stock (Bool)           ││  │  │ - in_stock (BOOLEAN)            ││  │
│  │  │ - created_at (String)       ││  │  │ - created_at (TIMESTAMP)        ││  │
│  │  └─────────────────────────────┘│  │  └─────────────────────────────────┘│  │
│  │  ┌─────────────────────────────┐│  │  ┌─────────────────────────────────┐│  │
│  │  │ orders table                ││  │  │ orders table                    ││  │
│  │  │ - order_id (UInt64)         ││  │  │ - order_id (SERIAL)             ││  │
│  │  │ - user_id (UInt64)          ││  │  │ - user_id (BIGINT)              ││  │
│  │  │ - items_json (String)       ││  │  │ - items_json (JSON)             ││  │
│  │  │ - subtotal (Float64)        ││  │  │ - subtotal (DECIMAL)            ││  │
│  │  │ - tax (Float64)             ││  │  │ - tax (DECIMAL)                 ││  │
│  │  │ - total (Float64)           ││  │  │ - total (DECIMAL)               ││  │
│  │  │ - status (String)           ││  │  │ - status (TEXT)                 ││  │
│  │  │ - created_at (String)       ││  │  │ - created_at (TIMESTAMP)        ││  │
│  │  └─────────────────────────────┘│  │  └─────────────────────────────────┘│  │
│  │  ┌─────────────────────────────┐│  │  ┌─────────────────────────────────┐│  │
│  │  │ invoices table              ││  │  │ invoices table                  ││  │
│  │  │ - invoice_id (UInt64)       ││  │  │ - invoice_id (SERIAL)           ││  │
│  │  │ - order_id (UInt64)         ││  │  │ - order_id (BIGINT)             ││  │
│  │  │ - user_id (UInt64)          ││  │  │ - user_id (BIGINT)              ││  │
│  │  │ - invoice_number (String)   ││  │  │ - invoice_number (TEXT)         ││  │
│  │  │ - items_json (String)       ││  │  │ - items_json (JSON)             ││  │
│  │  │ - subtotal (Float64)        ││  │  │ - subtotal (DECIMAL)            ││  │
│  │  │ - tax (Float64)             ││  │  │ - tax (DECIMAL)                 ││  │
│  │  │ - total (Float64)           ││  │  │ - total (DECIMAL)               ││  │
│  │  │ - status (String)           ││  │  │ - status (TEXT)                 ││  │
│  │  │ - created_at (String)       ││  │  │ - created_at (TIMESTAMP)        ││  │
│  │  │ - due_date (String)         ││  │  │ - due_date (TIMESTAMP)          ││  │
│  │  └─────────────────────────────┘│  │  └─────────────────────────────────┘│  │
│  │         Port: 9000              │  │         Port: 5432                  │  │
│  └─────────────────────────────────┘  └─────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Typical Request Flow Examples

### 1. User Registration Flow
```
POST /register → User Handler → User Service → User Repository → Database
                                    ↓
                               Password Hashing
                               Input Validation
                               Duplicate Check
```

### 2. User Login Flow
```
POST /login → User Handler → User Service → User Repository → Database
                                ↓              ↓
                          Password Verify  Find User
                                ↓
                           JWT Generation
                                ↓
                          Return Token
```

### 3. Create Order Flow
```
POST /orders/{username} → Order Handler → Order Service + User Service → Repositories → Database
                                             ↓              ↓
                                       Validate Order   Find User ID
                                             ↓
                                        Create Order
```

### 4. Get Orders Flow
```
GET /orders/{username} → Order Handler → Order Service + User Service → Repositories → Database
                                            ↓              ↓
                                      Get Orders     Find User ID
                                            ↓
                                    Return Order List
```

### 5. Product Catalog Flow (Public)
```
GET /products → Product Handler → Product Service → Product Repository → Database
                    ↓                    ↓
              No JWT Required      Get All Products
                    ↓                    ↓
             Return Product List   With Stock Status
```

### 6. Create Product Flow (Protected)
```
POST /products → Product Handler → JWT Middleware → Product Service → Repository → Database
                      ↓                  ↓               ↓
               Check Authorization  Validate Token  Validate Product Data
                      ↓                              ↓
                Create Product                Store Product
```

### 7. Invoice Generation Flow
```
POST /invoices → Invoice Handler → JWT Middleware → Invoice Service → Multiple Services
                      ↓                ↓                  ↓
               Check Authorization  Validate Token   Get Order Details
                                                           ↓
                                               Get Product Details
                                                           ↓
                                               Get User Details
                                                           ↓
                                            Generate Invoice Number
                                                           ↓
                                              Create Complete Invoice
```

### 8. Cross-Service Data Flow
```
Order Creation:
User → Order Service → Product Service (validate stock) → Order Repository
                              ↓
                        Update Product Stock

Invoice Generation:
User → Invoice Service → Order Service (get order)
                              ↓
                        Product Service (get details)
                              ↓
                        User Service (get user info)
                              ↓
                        Invoice Repository (store)
```

## Configuration & Environment

```
Environment Variables → Container → Repository Factory → Database Implementation

DB_DRIVER=postgres    →  PostgreSQL Repositories
DB_DRIVER=clickhouse  →  ClickHouse Repositories

Other Config:
- PORT (server port)
- JWT_SECRET (token signing)
- Database connection strings
```

## Key Architecture Benefits

1. **Separation of Concerns**: Each layer has a single responsibility
2. **Dependency Inversion**: High-level modules don't depend on low-level modules
3. **Database Agnostic**: Easy to switch between databases
4. **Testability**: Each layer can be unit tested independently
5. **Scalability**: Easy to add new domains and features
6. **Maintainability**: Clean, organized code structure

## File Organization

```
├── main.go                     # Application entry point & server setup
├── container.go               # Dependency injection container
├── db.go                      # Database initialization
├── middleware.go              # HTTP middleware functions
├── user/                      # User domain
│   ├── models.go             # User entities & DTOs
│   ├── service.go            # User business logic
│   ├── handler.go            # User HTTP handlers
│   ├── repository_clickhouse.go
│   └── repository_postgres.go
├── product/                   # Product domain
│   ├── models.go             # Product entities & DTOs
│   ├── service.go            # Product business logic
│   ├── handler.go            # Product HTTP handlers
│   ├── repository_clickhouse.go
│   └── repository_postgres.go
├── order/                     # Order domain
│   ├── models.go             # Order entities & DTOs
│   ├── service.go            # Order business logic
│   ├── handler.go            # Order HTTP handlers
│   ├── repository_clickhouse.go
│   └── repository_postgres.go
├── invoice/                   # Invoice domain
│   ├── models.go             # Invoice entities & DTOs
│   ├── service.go            # Invoice business logic
│   ├── handler.go            # Invoice HTTP handlers
│   ├── repository_clickhouse.go
│   └── repository_postgres.go
└── auth/                      # Authentication utilities
    └── auth.go               # JWT & password hashing
```