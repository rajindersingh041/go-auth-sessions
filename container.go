package main

import (
	"database/sql"
	"log"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/order"
	"github.com/rajindersingh041/go-auth-sessions/user"
)

// Container holds all services and dependencies
type Container struct {
	// Services
	UserService  user.Service
	OrderService order.Service

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

	switch dbDriver {
	case "clickhouse":
		userRepo = user.NewClickHouseRepository(db)
		orderRepo = order.NewClickHouseRepository(db)
	case "postgres":
		userRepo = user.NewPostgresRepository(db)
		orderRepo = order.NewPostgresRepository(db)
	default:
		log.Fatalf("Unsupported DB_DRIVER: %s", dbDriver)
	}

	// Create services
	userService := user.NewService(userRepo, passwordHasher)
	orderService := order.NewService(orderRepo)

	return &Container{
		UserService:    userService,
		OrderService:   orderService,
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