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
	"github.com/rajindersingh041/go-auth-sessions/invoice"
	"github.com/rajindersingh041/go-auth-sessions/order"
	"github.com/rajindersingh041/go-auth-sessions/product"
	"github.com/rajindersingh041/go-auth-sessions/user"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Initialize database connection
	// connection is established based on DB_DRIVER env variable
	// db contains the sql.DB pointer
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Get database driver from environment
	dbDriver := getEnv("DB_DRIVER", "clickhouse")
	log.Printf("Using database driver: %s", dbDriver)

	// Create dependency injection container with all services
	container := NewContainer(db, dbDriver)
	defer container.Close()

	// Initialize sample products
	if err := container.ProductService.InitializeSampleProducts(context.Background()); err != nil {
		log.Printf("Warning: Failed to initialize sample products: %v", err)
	}

	// Create HTTP handlers
	userHandler := user.NewHandler(container.UserService, container.JWTManager)
	orderHandler := order.NewHandler(container.OrderService, container.UserService, container.JWTManager)
	productHandler := product.NewHandler(container.ProductService, container.JWTManager)
	invoiceHandler := invoice.NewHandler(container.InvoiceService, container.JWTManager)

	// Setup HTTP server with routes
	server := setupServer(userHandler, orderHandler, productHandler, invoiceHandler)

	// Get port from environment
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s...", port)

	// Configure HTTP server with timeouts and security settings
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      server,
		ReadTimeout:  15 * time.Second,  // Maximum duration for reading the entire request
		WriteTimeout: 15 * time.Second,  // Maximum duration before timing out writes
		IdleTimeout:  60 * time.Second,  // Maximum duration to wait for the next request
	}

	// Start server in a goroutine for graceful shutdown
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Println("Server started successfully. Press Ctrl+C to stop.")

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

// setupServer configures HTTP routes and middleware
func setupServer(userHandler *user.Handler, orderHandler *order.Handler, productHandler *product.Handler, invoiceHandler *invoice.Handler) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", handleHealth())

	// Register domain-specific routes
	userHandler.RegisterRoutes(mux)
	orderHandler.RegisterRoutes(mux)
	productHandler.RegisterRoutes(mux)
	invoiceHandler.RegisterRoutes(mux)

	// Apply global middleware: logging, recovery, CORS, etc.
	handler := globalLoggingMiddleware(globalRecoveryMiddleware(mux))
	return handler
}

// handleHealth returns a simple health check endpoint
func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// globalLoggingMiddleware logs all incoming requests
func globalLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

// globalRecoveryMiddleware recovers from panics and returns 500 status
func globalRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}