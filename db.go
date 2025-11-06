package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

const (
	clickhouseDSN = "tcp://localhost:9000?username=default&password=MyPassword2025"
)


var dbConn *sql.DB

// initDB connects to ClickHouse and creates the users table.
func initDB() {
	var err error
	dbConn, err = sql.Open("clickhouse", clickhouseDSN)
	if err != nil {
		log.Fatalf("Error connecting to ClickHouse: %v", err)
	}

	// Ping to ensure connection is live
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Error pinging ClickHouse: %v", err)
	}

	log.Println("Connected to ClickHouse successfully.")

	createUserTable()
}

func createUserTable() {
	_, err := dbConn.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id UInt64 DEFAULT toUInt64(rand()),
		username String,
		password_hash String
	) ENGINE = MergeTree() ORDER BY user_id`)

	if err != nil {
		log.Printf("Issues in creating the table: %v", err)
		return
	}
	log.Println("Users table ensured.")
}

// createUser inserts a new user into ClickHouse.
func createUser(username string, passwordHash string) error {
	_, err := dbConn.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)",
		username, passwordHash)
	return err
}

// findUserByUsername retrieves a user's password hash from ClickHouse.
func findUserByUsername(username string) (string, error) {
	var passwordHash string
	ctx := context.Background()

	// This query (SELECT...WHERE...LIMIT 1) is very inefficient in ClickHouse
	// compared to a B-tree indexed table in an OLTP database.
	row := dbConn.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE username = ? LIMIT 1", username)
	if err := row.Scan(&passwordHash); err != nil {
		// clickhouse.ErrNoRows is not a specific error, so we check the text
		if err.Error() == "sql: no rows in result set" {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return passwordHash, nil
}

func findUserID(username string) (string, error) {
	var userID string
	ctx := context.Background()

	// This query (SELECT...WHERE...LIMIT 1) is very inefficient in ClickHouse
	// compared to a B-tree indexed table in an OLTP database.
	row := dbConn.QueryRowContext(ctx, "SELECT user_id FROM users WHERE username = ? LIMIT 1", username)
	if err := row.Scan(&userID); err != nil {
		// clickhouse.ErrNoRows is not a specific error, so we check the text
		if err.Error() == "sql: no rows in result set" {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return userID, nil
}

