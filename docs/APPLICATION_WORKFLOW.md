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
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────────┐  │
│  │ User Handler    │  │ Order Handler   │  │ Health Handler                  │  │
│  │ /register       │  │ /orders/{user}  │  │ /health                         │  │
│  │ /login          │  │ POST & GET      │  │                                 │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────────────┘  │
└──────────┬──────────────────┬─────────────────────────────────────────┬─────────┘
           │                  │                                         │
           ▼                  ▼                                         ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           DEPENDENCY INJECTION CONTAINER                       │
│  ┌─────────────────────────────────────────────────────────────────────────────┐│
│  │ • UserService    • OrderService    • JWTManager    • PasswordHasher       ││
│  │ • UserRepository • OrderRepository • Database Connection                   ││
│  └─────────────────────────────────────────────────────────────────────────────┘│
└─────────────────────────────┬───────────────────────────────────────────────────┘
                              │
                              ▼ Inject Dependencies
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              SERVICE LAYER                                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────────┐  │
│  │ User Service    │  │ Order Service   │  │ Auth Components                 │  │
│  │ • CreateUser()  │  │ • CreateOrder() │  │ • JWT Generation                │  │
│  │ • Authenticate()│  │ • GetOrders()   │  │ • Password Hashing              │  │
│  │ • Validation    │  │ • Validation    │  │ • Token Validation              │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────────────┘  │
└──────────┬──────────────────┬─────────────────────────────────────────────────────┘
           │                  │
           ▼                  ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            REPOSITORY LAYER                                    │
│  ┌─────────────────┐  ┌─────────────────┐                                      │
│  │ User Repository │  │ Order Repository│     ◄─── Database Abstraction        │
│  │ Interface       │  │ Interface       │          Layer                       │
│  └─────────────────┘  └─────────────────┘                                      │
│           │                      │                                             │
│           ▼                      ▼                                             │
│  ┌─────────────────┐  ┌─────────────────┐                                      │
│  │ ClickHouse Impl │  │ PostgreSQL Impl │     ◄─── Implementation             │
│  │ repository_     │  │ repository_     │          Selection via              │
│  │ clickhouse.go   │  │ postgres.go     │          DB_DRIVER env var          │
│  └─────────────────┘  └─────────────────┘                                      │
└──────────┬──────────────────┬─────────────────────────────────────────────────────┘
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
│  │  └─────────────────────────────┘│  │  └─────────────────────────────────┘│  │
│  │  ┌─────────────────────────────┐│  │  ┌─────────────────────────────────┐│  │
│  │  │ orders table                ││  │  │ orders table                    ││  │
│  │  │ - order_id (UInt64)         ││  │  │ - order_id (SERIAL)             ││  │
│  │  │ - user_id (UInt64)          ││  │  │ - user_id (BIGINT)              ││  │
│  │  │ - item (String)             ││  │  │ - item (TEXT)                   ││  │
│  │  │ - quantity (Int32)          ││  │  │ - quantity (INT)                ││  │
│  │  │ - created_at (String)       ││  │  │ - created_at (TIMESTAMP)        ││  │
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
├── order/                     # Order domain
│   ├── models.go             # Order entities & DTOs
│   ├── service.go            # Order business logic
│   ├── handler.go            # Order HTTP handlers
│   ├── repository_clickhouse.go
│   └── repository_postgres.go
└── auth/                      # Authentication utilities
    └── auth.go               # JWT & password hashing
```