package main

// User represents a user in the database
type User struct {
	UserID       uint64
	Username     string
	PasswordHash string
}
