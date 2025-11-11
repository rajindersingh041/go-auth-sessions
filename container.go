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
// It acts as a dependency injection container
// for easy management and testing
// Any new service or dependency should be added here
// TODO: Consider using a DI framework for larger projects
// TODO: Add logging, config, and other cross-cutting concerns as needed
// TODO: Implement proper cleanup methods if needed
// TODO: Add unit tests for container initialization

// What is Container?
// A Container is a struct that holds all the services and dependencies
// required by the application. It helps in managing dependencies
// and makes it easier to pass them around
// It is interfaced to build loosely coupled components
// that can be easily tested and maintained.
type Container struct {
	// Services
	// Services for different domains
	// Services encapsulate business logic that uses repositories and other components
	// to perform operations
	// Each service should ideally depend on interfaces for easier testing and flexibility
	// Example:
	// user.Service depends on user.Repository and auth.PasswordHasher interfaces
	// order.Service depends on order.Repository and product.Service interfaces
	// product.Service depends on product.Repository interface
	// invoice.Service depends on invoice.Repository, order.Service, product.Service, and user.Service interfaces
	// What is Service?
	// A Service is a struct that contains business logic methods
	// It uses repositories to interact with the database
	// It only purposes to perform operations related to a specific domain
	// For example, user.Service has methods for user registration, authentication, etc.
	// It uses user.Repository to perform database operations related to users
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
// and returns a Container instance
// It sets up repositories, services, and other components
// Why dbDriver is needed?
// The dbDriver parameter is needed to determine which database
// repositories to initialize (ClickHouse or Postgres).
// This allows the container to create the appropriate implementations
// based on the configured database driver.
func NewContainer(db *sql.DB, dbDriver string) *Container {
	// Create password hasher
	passwordHasher := &auth.BcryptPasswordHasher{}

	// Create JWT manager
	jwtManager := &auth.SimpleJWTManager{}

	// Create repositories based on database driver
	// What is repository?
	// A Repository is a struct that provides methods to interact with the database
	// It abstracts the database operations and provides a clean interface
	// for the services to use
	// Each repository should implement an interface to allow for easy swapping
	// of implementations (e.g., ClickHouse vs Postgres)
	var userRepo user.Repository
	var orderRepo order.Repository
	var productRepo product.Repository
	var invoiceRepo invoice.Repository


	// Initialize repositories based on dbDriver
	// Tables and schemas should be created as needed in the respective repository implementations
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
	// Services use repositories and other components to perform business logic
	// userService depends on userRepo and passwordHasher
	// productService depends on productRepo
	// orderService depends on orderRepo and productService
	// invoiceService depends on invoiceRepo, orderService, productService, and userService	
	userService := user.NewService(userRepo, passwordHasher)
	productService := product.NewService(productRepo)
	orderService := order.NewService(orderRepo, productService)
	invoiceService := invoice.NewService(invoiceRepo, orderService, productService, userService)

	// what is the purpose of newservice?
	// NewService functions create and return service instances
	// They take the required dependencies as parameters
	// and return a service that is ready to use
	// This helps in keeping the service initialization logic
	// separate and makes it easier to manage dependencies


	// Return container with all services and dependencies
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


// Example usage:
// func main() {
//     db, err := InitDB()
//     if err != nil {
//         log.Fatalf("Failed to initialize database: %v", err)
//     }
//     defer db.Close()
//
//     container := NewContainer(db, "clickhouse")
//     defer container.Close()
//
// Use container.UserService, container.OrderService, etc.
// 
// }


// Expected Output:
// A Container struct that holds all services and dependencies
// properly initialized based on the provided database connection
// and driver type.
