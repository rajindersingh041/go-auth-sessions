package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()
	// Initialize database
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Select repository implementation based on DB_DRIVER env variable
	dbDriver := getEnv("DB_DRIVER", "clickhouse")
	var userRepo UserRepository
	var orderRepo OrderRepository
	switch dbDriver {
	case "clickhouse":
		userRepo = NewClickHouseUserRepository(db)
		orderRepo = NewClickHouseOrderRepository(db)
	case "postgres":
		userRepo = NewPostgresUserRepository(db)
		orderRepo = NewPostgresOrderRepository(db)
	default:
		log.Fatalf("Unsupported DB_DRIVER: %s", dbDriver)
	}

	// Create password hasher and JWT manager
	passwordHasher := BcryptPasswordHasher{}
	jwtManager := SimpleJWTManager{}

	// Create and configure server
	server := NewServer(userRepo, orderRepo, passwordHasher, jwtManager)

	// Get port from environment
	port := getEnv("PORT", "8080")

	// Configure HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      server.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s...", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}