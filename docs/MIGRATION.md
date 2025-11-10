# Database Schema and Migration Guide

## Complete Database Schema (Current)

This document describes the current database schema for all four domains: **User**, **Product**, **Order**, and **Invoice**.

## Users Table

### PostgreSQL
```sql
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### ClickHouse
```sql
CREATE TABLE users (
    user_id UInt64,
    username String,
    password_hash String,
    created_at String
) ENGINE = MergeTree()
ORDER BY user_id;
```

## Products Table

### PostgreSQL
```sql
CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category TEXT NOT NULL,
    in_stock BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### ClickHouse
```sql
CREATE TABLE products (
    product_id UInt64,
    name String,
    description String,
    price Decimal64(2),
    category String,
    in_stock UInt8,
    created_at String
) ENGINE = MergeTree()
ORDER BY product_id;
```

## Orders Table (Current Schema)

### PostgreSQL
```sql
CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    items JSONB NOT NULL,
    subtotal DECIMAL(10,2) DEFAULT 0.00,
    tax DECIMAL(10,2) DEFAULT 0.00,
    total DECIMAL(10,2) DEFAULT 0.00,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### ClickHouse
```sql
CREATE TABLE orders (
    order_id UInt64,
    user_id UInt64,
    items String,
    subtotal Decimal64(2),
    tax Decimal64(2),
    total Decimal64(2),
    status String,
    created_at String
) ENGINE = MergeTree()
ORDER BY order_id;
```

## Invoices Table

### PostgreSQL
```sql
CREATE TABLE IF NOT EXISTS invoices (
    invoice_id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    username TEXT NOT NULL,
    invoice_number TEXT NOT NULL UNIQUE,
    items JSONB,
    subtotal DECIMAL(10,2) DEFAULT 0.00,
    tax DECIMAL(10,2) DEFAULT 0.00,
    total DECIMAL(10,2) DEFAULT 0.00,
    status TEXT NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL
);
```

### ClickHouse
```sql
CREATE TABLE invoices (
    invoice_id UInt64,
    order_id UInt64,
    user_id UInt64,
    username String,
    invoice_number String,
    items String,
    subtotal Decimal64(2),
    tax Decimal64(2),
    total Decimal64(2),
    status String,
    created_at String,
    due_date String
) ENGINE = MergeTree()
ORDER BY invoice_id;
```

## Migration from Legacy Schema

### Legacy Orders Table Migration

If you have an existing orders table with the old schema, run these commands:

**PostgreSQL:**
```sql
-- Drop the old orders table (CAUTION: This deletes existing data)
DROP TABLE IF EXISTS orders;

-- Create the new orders table with updated schema
    order_id UInt64,
    user_id UInt64,
    product_id UInt64,
    quantity UInt32,
    created_at String
) ENGINE = MergeTree()
ORDER BY order_id;
```

## Products Table Creation

### PostgreSQL
```sql
CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category TEXT NOT NULL,
    in_stock BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### ClickHouse
```sql
CREATE TABLE products (
    product_id UInt64,
    name String,
    description String,
    price Decimal64(2),
    category String,
    in_stock UInt8,
    created_at String
) ENGINE = MergeTree()
ORDER BY product_id;
```

## Invoices Table Creation

### PostgreSQL
```sql
CREATE TABLE IF NOT EXISTS invoices (
    invoice_id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    username TEXT NOT NULL,
    invoice_number TEXT NOT NULL UNIQUE,
    items JSONB,
    subtotal DECIMAL(10,2) DEFAULT 0.00,
    tax DECIMAL(10,2) DEFAULT 0.00,
    total DECIMAL(10,2) DEFAULT 0.00,
    status TEXT NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL
);
```

### ClickHouse
```sql
CREATE TABLE invoices (
    invoice_id UInt64,
    order_id UInt64,
    user_id UInt64,
    username String,
    invoice_number String,
    items String,
    subtotal Decimal64(2),
    tax Decimal64(2),
    total Decimal64(2),
    status String,
    created_at String,
    due_date String
) ENGINE = MergeTree()
ORDER BY invoice_id;
```

## Database Relationships

```sql
-- Foreign Key Relationships (PostgreSQL)
-- Note: ClickHouse doesn't enforce foreign keys but the relationships exist logically

users.user_id → orders.user_id
users.user_id → invoices.user_id
orders.order_id → invoices.order_id

-- JSON Structure in orders.items and invoices.items:
{
  "product_id": 123,
  "product_name": "Product Name", 
  "quantity": 2,
  "unit_price": 99.99,
  "total": 199.98
}
```

## Auto-Migration Features

### PostgreSQL
- **Automatic Table Creation**: Tables are created automatically if they don't exist
- **Schema Evolution**: The application handles schema changes gracefully
- **Sample Data**: Products are automatically seeded on first startup

### ClickHouse  
- **Manual Schema Setup**: Tables need to be created manually using the SQL above
- **No Foreign Keys**: Relationships are maintained at the application level
- **Performance Optimized**: Uses appropriate engines and ordering for analytics

## Environment Configuration

```env
# Choose your database driver
DB_DRIVER=postgres          # or clickhouse

# PostgreSQL Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=mysecretpassword
POSTGRES_DB=authdb

# ClickHouse Configuration  
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=MyPassword2025
CLICKHOUSE_DB=default
```

## Migration Notes

- **Data Loss Warning**: Schema migrations may require dropping tables with existing data
- **Backup First**: Always backup your data before running migrations
- **Development Mode**: Tables are auto-created in development environments
- **Production**: Use proper migration tools (golang-migrate, Flyway, etc.) for production deployments