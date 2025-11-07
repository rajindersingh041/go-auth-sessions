package main

import (
	"log"
	"net/http"
)

func main() {
	// Connect to ClickHouse and set up the table
	initDB()

	// mux added
	mux := http.NewServeMux()

	// Use Go 1.22+ method-based routing.
	// This is much cleaner and more secure.
	mux.HandleFunc("POST /register", registerHandler)
	mux.HandleFunc("POST /login", loginHandler)
	
	// You can also apply middleware the same way
	protectedHandler := authMiddleware(http.HandlerFunc(protectedHandler))
	mux.Handle("GET /protected", protectedHandler) // Only accepts GET

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}}

