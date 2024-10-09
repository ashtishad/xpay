package domain

import (
	"time"

	"github.com/google/uuid"
)

// Constants for card providers, types, and statuses
const (
	CardProviderVisa       = "visa"
	CardProviderMastercard = "mastercard"
	CardProviderAmex       = "amex"

	CardTypeCredit = "credit"
	CardTypeDebit  = "debit"

	CardStatusActive   = "active"
	CardStatusInactive = "inactive"
	CardStatusDeleted  = "deleted"
)

type Card struct {
	ID                  int64     `json:"-"`
	UUID                uuid.UUID `json:"cardId"`
	UserID              int64     `json:"-"`
	WalletID            int64     `json:"-"`
	EncryptedCardNumber []byte    `json:"-"`
	Provider            string    `json:"provider"`
	Type                string    `json:"type"`
	LastFour            string    `json:"lastFour"`
	ExpiryDate          time.Time `json:"expiryDate"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// IsValidCardProvider utility method to validate queryParams, request body is validated with validator/v10
func IsValidCardProvider(provider string) bool {
	if provider == CardProviderAmex || provider == CardProviderMastercard || provider == CardProviderVisa {
		return true
	}

	return false
}

// IsValidCardStatus utility method to validate queryParams, request body is validated with validator/v10
func IsValidCardStatus(status string) bool {
	if status == CardStatusActive || status == CardStatusInactive || status == CardStatusDeleted {
		return true
	}

	return false
}
