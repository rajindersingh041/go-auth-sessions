# Database Cleanup and Migration

## PostgreSQL Migration

If you have an existing orders table with the old schema, run these commands to clean it up:

```sql
-- Drop the old orders table completely (this will delete existing order data)
DROP TABLE IF EXISTS orders;

-- Create the new orders table with the correct schema
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## ClickHouse Migration

For ClickHouse, you'll need to create the orders table manually:

```sql
-- Drop the old orders table if it exists
DROP TABLE IF EXISTS orders;

-- Create the new orders table
CREATE TABLE orders (
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

## Notes

- **Data Loss Warning**: The migration will delete existing order data since we're changing the schema from `item` (text) to `product_id` (integer).
- **Auto Migration**: The PostgreSQL repository now includes automatic migration logic that will handle the schema update.
- **Manual Setup**: For ClickHouse, you may need to run the table creation commands manually since ClickHouse has limited ALTER support.