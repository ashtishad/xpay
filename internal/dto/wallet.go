package dto

import (
	"strings"
	"time"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/google/uuid"
)

type CreateWalletRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD"`
}

// ToNewWallet converts CreateWalletRequest to *domain.Wallet
func (r *CreateWalletRequest) ToNewWallet(userID int64) *domain.Wallet {
	now := time.Now().UTC()
	return &domain.Wallet{
		UUID:           uuid.New(),
		UserID:         userID,
		BalanceInCents: 0,
		Currency:       strings.ToUpper(r.Currency),
		Status:         domain.WalletStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

type CreateWalletResponse struct {
	Wallet domain.Wallet `json:"wallet"`
}

type GetWalletBalanceResponse struct {
	BalanceInCents int64  `json:"balanceInCents"`
	Currency       string `json:"currency"`
}

type UpdateWalletStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive blocked"`
}
