package secure

import (
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

var (
	errIncorrectPassword = errors.New("incorrect password")
	errEmptyPassword     = errors.New("password should not be empty")
	errPasswordTooShort  = errors.New("password should be at least 8 characters long")
	errPasswordTooLong   = errors.New("password should be at max 64 characters long")
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

// VerifyPassword verifies bcrypt hash and the given password text
func VerifyPassword(hashedPassword, password string) error {
	if err := validatePasswordText(password); err != nil {
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errIncorrectPassword
		}

		slog.Error("failed to verify password", "err", err)
		return err
	}

	return nil
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
