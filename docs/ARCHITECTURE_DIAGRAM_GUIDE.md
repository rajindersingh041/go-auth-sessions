# Updated Architecture Diagram for go-auth-sessions

## Clean Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer                            │
├─────────────────┬─────────────────┬─────────────────────────┤
│   User Handler  │  Order Handler  │    Health Endpoint      │
│   - Register    │   - Create      │    - Status Check       │
│   - Login       │   - Get Orders  │                         │
└─────────────────┴─────────────────┴─────────────────────────┘
         │                  │                  │
         ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                            │
├─────────────────┬─────────────────┬─────────────────────────┤
│   User Service  │  Order Service  │    Auth Components      │
│   - Create User │  - Create Order │    - JWT Manager        │
│   - Authenticate│  - Get Orders   │    - Password Hasher    │
│   - Get User    │                 │                         │
└─────────────────┴─────────────────┴─────────────────────────┘
         │                  │                  │
         ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  Repository Layer                           │
├─────────────────┬─────────────────┬─────────────────────────┤
│  User Repository│ Order Repository│   Database Abstraction  │
│  - ClickHouse   │  - ClickHouse   │   - Connection Pool     │
│  - PostgreSQL   │  - PostgreSQL   │   - Transaction Mgmt    │
└─────────────────┴─────────────────┴─────────────────────────┘
         │                  │
         ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                     Database Layer                          │
├─────────────────────────┬───────────────────────────────────┤
│      ClickHouse         │         PostgreSQL                │
│   - users table         │      - users table                │
│   - orders table        │      - orders table               │
└─────────────────────────┴───────────────────────────────────┘
```

## Component Interactions

### 1. Client Request Flow
```
Client → HTTP Handler → Service → Repository → Database
```

### 2. User Registration Flow
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