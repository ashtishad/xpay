package secure

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	errEmptyPassword    = errors.New("password should not be empty")
	errPasswordTooShort = errors.New("password should be at least 8 characters long")
	errPasswordTooLong  = errors.New("password should be at max 64 characters long")
)

// GeneratePasswordHash creates a bcrypt hash of the given password text.
// Requires password to be at least 8 characters long and 64 chars at max.
// Returns the hashed password as a string or an error if hashing fails.
func GeneratePasswordHash(password string) (string, error) {
	if err := validatePasswordText(password); err != nil {
		return "", err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func validatePasswordText(password string) error {
	switch {
	case password == "":
		return errEmptyPassword
	case len(password) < 8:
		return errPasswordTooShort
	case len(password) > 64:
		return errPasswordTooLong
	default:
		return nil
	}
}
