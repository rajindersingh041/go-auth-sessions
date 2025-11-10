package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/lib/pq"
)

// ...User struct moved to models.go...
// ...UserRepository struct and methods moved to repository.go...

// InitDB initializes the database connection (ClickHouse or Postgres) with proper configuration
func InitDB() (*sql.DB, error) {
	driver := getEnvOrDefault("DB_DRIVER", "clickhouse")
	var db *sql.DB
	var err error
	var dsn string
	if driver == "postgres" {
		dsn = getPostgresDSN()
		log.Print("Postgres DSN is ", dsn)
		db, err = sql.Open("postgres", dsn)
	} else {
		dsn = getDSN()
		log.Print("ClickHouse DSN is ", dsn)
		db, err = sql.Open("clickhouse", dsn)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Ping database with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if driver == "postgres" {
		log.Println("Connected to Postgres successfully")
	} else {
		log.Println("Connected to ClickHouse successfully")
	}

	// Create users table (ClickHouse only, Postgres handled in repo)
	if driver == "clickhouse" {
		if err := createUserTable(db); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

// getPostgresDSN constructs the Postgres connection string from environment variables
func getPostgresDSN() string {
	host := getEnvOrDefault("POSTGRES_HOST", "localhost")
	port := getEnvOrDefault("POSTGRES_PORT", "5432")
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	pass := getEnvOrDefault("POSTGRES_PASSWORD", "mysecretpassword")
	dbname := getEnvOrDefault("POSTGRES_DB", "authdb")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbname, sslmode)
}

// getDSN constructs the database connection string from environment variables
func getDSN() string {
	host := getEnvOrDefault("CLICKHOUSE_HOST", "localhost")
	port := getEnvOrDefault("CLICKHOUSE_PORT", "9000")
	user := getEnvOrDefault("CLICKHOUSE_USER", "default")
	pass := getEnvOrDefault("CLICKHOUSE_PASSWORD", "MyPassword2025")

	return fmt.Sprintf("tcp://%s:%s?username=%s&password=%s",
		host, port, user, pass)
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// createUserTable creates the users table if it doesn't exist
func createUserTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		CREATE TABLE IF NOT EXISTS users (
			user_id UInt64 DEFAULT toUInt64(rand()),
			username String,
			password_hash String
		) ENGINE = MergeTree() 
		ORDER BY user_id
	`

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Users table ensured")
	return nil
}

