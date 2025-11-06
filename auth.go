// auth.go
package main

import "golang.org/x/crypto/bcrypt"

// User struct for parsing registration/login requests
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// hashPassword creates a bcrypt hash of the password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is a good cost
	return string(bytes), err
}

// checkPasswordHash compares a plain-text password with a hash.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}