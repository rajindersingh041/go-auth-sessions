package product

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PostgresRepository implements Repository for PostgreSQL database
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL product repository
func NewPostgresRepository(db *sql.DB) ProductRepository {
	return &PostgresRepository{db: db}
}

// ensureProductsTable creates the products table if it doesn't exist
func (r *PostgresRepository) ensureProductsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS products (
			product_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			category TEXT NOT NULL,
			in_stock BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW()
		)`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, product *Product) error {
	if err := r.ensureProductsTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO products (name, description, price, category, in_stock, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Category, product.InStock, product.CreatedAt)
	return err
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Product, error) {
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
		var createdAt time.Time
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Category, &p.InStock, &createdAt); err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt.Format(time.RFC3339)
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, productID uint64) (*Product, error) {
	if err := r.ensureProductsTable(ctx); err != nil {
		return nil, err
	}
	var product Product
	var createdAt time.Time
	query := "SELECT product_id, name, description, price, category, in_stock, created_at FROM products WHERE product_id = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ProductID, &product.Name, &product.Description, &product.Price, &product.Category, &product.InStock, &createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	product.CreatedAt = createdAt.Format(time.RFC3339)
	return &product, nil
}

func (r *PostgresRepository) GetByCategory(ctx context.Context, category string) ([]Product, error) {
	if err := r.ensureProductsTable(ctx); err != nil {
		return nil, err
	}
	query := "SELECT product_id, name, description, price, category, in_stock, created_at FROM products WHERE category = $1 ORDER BY name"
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var createdAt time.Time
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Category, &p.InStock, &createdAt); err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt.Format(time.RFC3339)
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresRepository) UpdateStock(ctx context.Context, productID uint64, inStock bool) error {
	if err := r.ensureProductsTable(ctx); err != nil {
		return err
	}
	query := "UPDATE products SET in_stock = $1 WHERE product_id = $2"
	_, err := r.db.ExecContext(ctx, query, inStock, productID)
	return err
}

func (r *PostgresRepository) SeedSampleProducts(ctx context.Context) error {
	if err := r.ensureProductsTable(ctx); err != nil {
		return err
	}

	// Check if products already exist
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
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