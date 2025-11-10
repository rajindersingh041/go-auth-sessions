package product

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ClickHouseRepository implements Repository for ClickHouse database
type ClickHouseRepository struct {
	db *sql.DB
}

// NewClickHouseRepository creates a new ClickHouse product repository
func NewClickHouseRepository(db *sql.DB) Repository {
	return &ClickHouseRepository{db: db}
}

// ensureProductsTable creates the products table if it doesn't exist
func (r *ClickHouseRepository) ensureProductsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS products (
			product_id UInt64 DEFAULT toUInt64(rand()),
			name String,
			description String,
			price Float64,
			category String,
			in_stock Bool,
			created_at String
		) ENGINE = MergeTree() 
		ORDER BY product_id
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *ClickHouseRepository) Create(ctx context.Context, product *Product) error {
	if err := r.ensureProductsTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO products (name, description, price, category, in_stock, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Category, product.InStock, product.CreatedAt)
	return err
}

func (r *ClickHouseRepository) GetAll(ctx context.Context) ([]Product, error) {
	if err := r.ensureProductsTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT product_id, name, description, price, category, in_stock, created_at FROM products ORDER BY name"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Category, &p.InStock, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ClickHouseRepository) GetByID(ctx context.Context, productID uint64) (*Product, error) {
	if err := r.ensureProductsTable(ctx); err != nil {
		return nil, err
	}
	var product Product
	query := "SELECT product_id, name, description, price, category, in_stock, created_at FROM products WHERE product_id = ? LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ProductID, &product.Name, &product.Description, &product.Price, &product.Category, &product.InStock, &product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *ClickHouseRepository) GetByCategory(ctx context.Context, category string) ([]Product, error) {
	if err := r.ensureProductsTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT product_id, name, description, price, category, in_stock, created_at FROM products WHERE category = ? ORDER BY name"
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Category, &p.InStock, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ClickHouseRepository) UpdateStock(ctx context.Context, productID uint64, inStock bool) error {
	if err := r.ensureProductsTable(ctx); err != nil {
		return err
	}
	query := "ALTER TABLE products UPDATE in_stock = ? WHERE product_id = ?"
	_, err := r.db.ExecContext(ctx, query, inStock, productID)
	return err
}

func (r *ClickHouseRepository) SeedSampleProducts(ctx context.Context) error {
	if err := r.ensureProductsTable(ctx); err != nil {
		return err
	}

	// Check if products already exist
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT count() FROM products").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Products already exist
	}

	// Sample products data
	sampleProducts := []Product{
		{Name: "MacBook Pro 16\"", Description: "High-performance laptop for professionals", Price: 2499.99, Category: "Electronics", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "iPhone 15 Pro", Description: "Latest smartphone with advanced features", Price: 999.99, Category: "Electronics", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "Wireless Headphones", Description: "Premium noise-cancelling headphones", Price: 299.99, Category: "Electronics", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "Coffee Maker", Description: "Automatic drip coffee maker", Price: 89.99, Category: "Appliances", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "Office Chair", Description: "Ergonomic office chair with lumbar support", Price: 199.99, Category: "Furniture", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "Bluetooth Speaker", Description: "Portable wireless speaker", Price: 49.99, Category: "Electronics", InStock: false, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "Desk Lamp", Description: "LED desk lamp with adjustable brightness", Price: 39.99, Category: "Furniture", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "Water Bottle", Description: "Insulated stainless steel water bottle", Price: 24.99, Category: "Accessories", InStock: true, CreatedAt: time.Now().Format(time.RFC3339)},
	}

	for _, product := range sampleProducts {
		if err := r.Create(ctx, &product); err != nil {
			return fmt.Errorf("failed to seed product %s: %w", product.Name, err)
		}
	}

	return nil
}