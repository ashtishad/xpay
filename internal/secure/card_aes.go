package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
)

// CardEncryptor provides methods for encrypting and decrypting card numbers using AES-GCM.
type CardEncryptor struct {
	gcm cipher.AEAD
}

// NewCardEncryptor creates a new CardEncryptor instance with the provided AES key.
// It initializes the AES cipher and GCM mode for encryption and decryption.
func NewCardEncryptor(key string) (*CardEncryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		err := errors.New("invalid AES key size: must be 16, 24, or 32 bytes")
		slog.Error("Failed to create CardEncryptor", "error", err)
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		slog.Error("Failed to create AES cipher", "error", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("Failed to create GCM", "error", err)
		return nil, err
	}

	return &CardEncryptor{gcm: gcm}, nil
}

// Encrypt takes a plaintext card number, validates it, and encrypts it using AES-GCM.
// It returns the encrypted card number as a byte slice, with the nonce prepended.
func (ce *CardEncryptor) Encrypt(plaintext string) ([]byte, error) {
	if err := ce.validateCardNumber(plaintext); err != nil {
		slog.Error("Invalid card number during encryption", "error", err)
		return nil, fmt.Errorf("invalid card number: %w", err)
	}

	nonce := make([]byte, ce.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		slog.Error("Failed to generate nonce for encryption", "error", err)
		return nil, err
	}

	return ce.gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

// Decrypt takes an encrypted card number (with prepended nonce), decrypts it using AES-GCM,
// and validates the decrypted card number. It returns the decrypted card number as a string.
func (ce *CardEncryptor) Decrypt(ciphertext []byte) (string, error) {
	if len(ciphertext) < ce.gcm.NonceSize() {
		err := errors.New("ciphertext too short")
		slog.Error("Decryption failed", "error", err)
		return "", err
	}

	nonce, ciphertext := ciphertext[:ce.gcm.NonceSize()], ciphertext[ce.gcm.NonceSize():]
	plaintext, err := ce.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		slog.Error("Failed to decrypt ciphertext", "error", err)
		return "", err
	}

	decryptedNumber := string(plaintext)
	if err := ce.validateCardNumber(decryptedNumber); err != nil {
		slog.Error("Decrypted data is not a valid card number", "error", err)
		return "", fmt.Errorf("decrypted data is not a valid card number: %w", err)
	}

	return decryptedNumber, nil
}

// validateCardNumber checks if the provided card number is valid for Visa, Mastercard, or American Express
// using regex patterns. It returns an error if the card number is empty or doesn't match any valid pattern.
func (ce *CardEncryptor) validateCardNumber(cardNumber string) error {
	if cardNumber == "" {
		err := errors.New("card number cannot be empty")
		slog.Error("Card number validation failed", "error", err)
		return err
	}

	patterns := map[string]string{
		"Visa":             "^4[0-9]{12}(?:[0-9]{3})?$",
		"Mastercard":       "^5[1-5][0-9]{14}$",
		"American Express": "^3[47][0-9]{13}$",
	}

	for _, pattern := range patterns {
		ok, _ := regexp.MatchString(pattern, cardNumber)
		if ok {
			return nil
		}
	}

	err := errors.New("invalid card number format")
	return err
}
