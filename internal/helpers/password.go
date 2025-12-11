package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Custom errors for better error handling
var (
	ErrPasswordEmpty   = errors.New("password cannot be empty")
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	MinPasswordLength = 6
	MaxPasswordLength = 72
	DefaultCost       = 12
)

// HashPassword generates a bcrypt hash using the default cost
// Optimized for dating app performance requirements
func HashPassword(password string) (string, error) {
	return HashPasswordWithCost(password, DefaultCost)
}

// HashPasswordWithCost generates a bcrypt hash with specified cost
func HashPasswordWithCost(password string, cost int) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a password with its bcrypt hash
// Returns nil if password matches, error otherwise
func ComparePassword(hashedPassword, password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	// bcrypt.CompareHashAndPassword is constant-time and handles all validation
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}

	return nil
}

// GetHashCost extracts the cost parameter from a bcrypt hash
// Useful for monitoring hash strength in production
func GetHashCost(hashedPassword string) (int, error) {
	cost, err := bcrypt.Cost([]byte(hashedPassword))
	if err != nil {
		return 0, err
	}
	return cost, nil
}

// NeedsRehash checks if a password hash needs to be updated
// Useful for upgrading hash strength over time
func NeedsRehash(hashedPassword string, targetCost int) bool {
	cost, err := GetHashCost(hashedPassword)
	if err != nil {
		return true // If we can't determine cost, assume it needs rehashing
	}
	return cost < targetCost
}
