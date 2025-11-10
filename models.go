package main

// User represents a user in the database
type User struct {
	UserID       uint64
	Username     string
	PasswordHash string
}

// Order represents an order placed by a user
type Order struct {
	OrderID   uint64
	UserID    uint64
	Item      string
	Quantity  int
	CreatedAt string // or time.Time if you want to use time
}
