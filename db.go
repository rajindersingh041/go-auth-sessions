package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

// User represents a user in the database
type User struct {
	UserID       uint64
	Username     string
	PasswordHash string
}

// UserRepository handles all database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// InitDB initializes the ClickHouse connection with proper configuration
func InitDB() (*sql.DB, error) {
	dsn := getDSN()
	log.Print("DSN is ", dsn)
	db, err := sql.Open("clickhouse", dsn)
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

	log.Println("Connected to ClickHouse successfully")

	// Create users table
	if err := createUserTable(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
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

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) error {
	query := "INSERT INTO users (username, password_hash) VALUES (?, ?)"

	_, err := r.db.ExecContext(ctx, query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByUsername retrieves a user by their username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	query := "SELECT user_id, username, password_hash FROM users WHERE username = ? LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// FindUserID retrieves only the user ID by username
func (r *UserRepository) FindUserID(ctx context.Context, username string) (uint64, error) {
	var userID uint64

	query := "SELECT user_id FROM users WHERE username = ? LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to query user ID: %w", err)
	}

	return userID, nil
}

// UserExists checks if a user with the given username exists
func (r *UserRepository) UserExists(ctx context.Context, username string) (bool, error) {
	var count int

	query := "SELECT count() FROM users WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}