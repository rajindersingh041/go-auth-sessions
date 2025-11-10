package main

import (
	"database/sql"
	"log"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/invoice"
	"github.com/rajindersingh041/go-auth-sessions/order"
	"github.com/rajindersingh041/go-auth-sessions/product"
	"github.com/rajindersingh041/go-auth-sessions/user"
)

// Container holds all services and dependencies
type Container struct {
	// Services
	UserService    user.Service
	OrderService   order.Service
	ProductService product.Service
	InvoiceService invoice.Service

	// Auth components
	JWTManager     auth.JWTManager
	PasswordHasher auth.PasswordHasher

	// Database
	DB *sql.DB
}

// NewContainer creates and configures all dependencies
func NewContainer(db *sql.DB, dbDriver string) *Container {
	// Create password hasher
	passwordHasher := &auth.BcryptPasswordHasher{}

	// Create JWT manager
	jwtManager := &auth.SimpleJWTManager{}

	// Create repositories based on database driver
	var userRepo user.Repository
	var orderRepo order.Repository
	var productRepo product.Repository
	var invoiceRepo invoice.Repository

	switch dbDriver {
	case "clickhouse":
		userRepo = user.NewClickHouseRepository(db)
		orderRepo = order.NewClickHouseRepository(db)
		productRepo = product.NewClickHouseRepository(db)
		invoiceRepo = invoice.NewClickHouseRepository(db)
	case "postgres":
		userRepo = user.NewPostgresRepository(db)
		orderRepo = order.NewPostgresRepository(db)
		productRepo = product.NewPostgresRepository(db)
		invoiceRepo = invoice.NewPostgresRepository(db)
	default:
		log.Fatalf("Unsupported DB_DRIVER: %s", dbDriver)
	}

	// Create services
	userService := user.NewService(userRepo, passwordHasher)
	productService := product.NewService(productRepo)
	orderService := order.NewService(orderRepo, productService)
	invoiceService := invoice.NewService(invoiceRepo, orderService, productService, userService)

	return &Container{
		UserService:    userService,
		OrderService:   orderService,
		ProductService: productService,
		InvoiceService: invoiceService,
		JWTManager:     jwtManager,
		PasswordHasher: passwordHasher,
		DB:             db,
	}
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}