package auth

// PasswordHasher interface for password hashing operations
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) error
}

// JWTManager interface for JWT operations
type JWTManager interface {
	GenerateToken(username string) (string, error)
	ValidateToken(tokenString string) (string, error)
}
