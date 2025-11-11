package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implements PasswordHasher using bcrypt
type BcryptPasswordHasher struct{}

func (b BcryptPasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b BcryptPasswordHasher) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


// What is bcrypt?
// Bcrypt is a password hashing function designed to be computationally intensive
// to protect against brute-force attacks. It incorporates a salt to protect
// against rainbow table attacks and is adaptive, meaning the cost factor can
// be increased over time as computational power increases, enhancing security.

// BcryptPasswordHasher provides methods to hash passwords and verify them
// using the bcrypt algorithm. It implements the PasswordHasher interface
// defined in interfaces.go.

// Example usage:
//	hasher := &auth.BcryptPasswordHasher{}
//	hashedPassword, err := hasher.HashPassword("mysecretpassword")
//	err = hasher.CheckPassword("mysecretpassword", hashedPassword)
// In this example, we create an instance of BcryptPasswordHasher,
// hash a password, and then verify it against the hashed value.
