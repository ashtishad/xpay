package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	WalletStatusActive   = "active"
	WalletStatusInactive = "inactive"
	WalletStatusBlocked  = "blocked"

	WalletCurrencyUSD = "USD"
)

type Wallet struct {
	ID             int64     `json:"-"`
	UUID           uuid.UUID `json:"uuid"`
	UserID         int64     `json:"-"`
	BalanceInCents int64     `json:"balanceInCents"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
