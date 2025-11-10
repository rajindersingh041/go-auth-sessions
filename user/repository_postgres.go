package user

import (
	"context"
	"database/sql"
)

// PostgresRepository implements Repository for PostgreSQL database
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL user repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ensureUsersTable creates the users table if it doesn't exist
func (r *PostgresRepository) ensureUsersTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	)`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, username, passwordHash string) error {
	if err := r.ensureUsersTable(ctx); err != nil {
		return err
	}
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, username, passwordHash)
	return err
}

func (r *PostgresRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	if err := r.ensureUsersTable(ctx); err != nil {
		return nil, err
	}
	var user User
	query := "SELECT user_id, username, password_hash FROM users WHERE username = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.UserID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, userID uint64) (*User, error) {
	if err := r.ensureUsersTable(ctx); err != nil {
		return nil, err
	}
	var user User
	query := "SELECT user_id, username, password_hash FROM users WHERE user_id = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.UserID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) FindUserID(ctx context.Context, username string) (uint64, error) {
	if err := r.ensureUsersTable(ctx); err != nil {
		return 0, err
	}
	var userID uint64
	query := "SELECT user_id FROM users WHERE username = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func (r *PostgresRepository) UserExists(ctx context.Context, username string) (bool, error) {
	if err := r.ensureUsersTable(ctx); err != nil {
		return false, err
	}
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}