package dto

import (
	"errors"
	"time"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/google/uuid"
)

// AddCardRequest represents the request body for adding a new card.
// @Description AddCardRequest validates input for adding a new card.
// @Description CardNumber must be a valid credit card number between 13 and 19 digits.
// @Description Provider must be one of: visa, mastercard, or amex.
// @Description Type must be either credit or debit.
// @Description ExpiryDate must be a future date and "MM/YY" format.
// @Description CVV must be minimum 3 and max 4 four digits.
type AddCardRequest struct {
	CardNumber string `json:"cardNumber" binding:"required,credit_card,min=13,max=19"`
	Provider   string `json:"provider" binding:"required,oneof=visa mastercard amex"`
	Type       string `json:"type" binding:"required,oneof=credit debit"`
	ExpiryDate string `json:"expiryDate" binding:"required,len=5"`
	CVV        string `json:"cvv" binding:"required,min=3,max=4"`
}

// ToCard converts AddCardRequest to domain.Card
func (r *AddCardRequest) ToCard(userID, walletID int64, encryptedCardNumber []byte) (*domain.Card, error) {
	expiryDate, err := parseExpiryDate(r.ExpiryDate)
	if err != nil {
		return nil, err
	}

	return &domain.Card{
		UUID:                uuid.New(),
		UserID:              userID,
		WalletID:            walletID,
		EncryptedCardNumber: encryptedCardNumber,
		Provider:            r.Provider,
		Type:                r.Type,
		LastFour:            r.CardNumber[len(r.CardNumber)-4:],
		ExpiryDate:          expiryDate,
		Status:              domain.CardStatusActive,
	}, nil
}

// AddCardResponse contains the card data returned after successfully adding a card.
// @Description AddCardResponse includes the created card's details.
type AddCardResponse struct {
	Card CardResponse `json:"card"`
}

// UpdateCardRequest represents the request body for updating an existing card.
// @Description UpdateCardRequest validates input for updating a card.
// @Description ExpiryDate, if provided, must be a future date.
// @Description Status, if provided, must be either active or inactive.
type UpdateCardRequest struct {
	ExpiryDate *string `json:"expiryDate,omitempty" binding:"omitempty,len=5"`
	Status     *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
}

// UpdateCard applies the update request to an existing card
func (r *UpdateCardRequest) UpdateCard(card *domain.Card) (*domain.Card, error) {
	if r.ExpiryDate != nil {
		expiryDate, err := parseExpiryDate(*r.ExpiryDate)
		if err != nil {
			return nil, err
		}

		card.ExpiryDate = expiryDate
	}

	if r.Status != nil {
		card.Status = *r.Status
	}

	return card, nil
}

// CardResponse represents the response body for card operations.
// @Description CardResponse includes the card's details, excluding sensitive information.
type CardResponse struct {
	UUID       uuid.UUID `json:"uuid"`
	Provider   string    `json:"provider"`
	Type       string    `json:"type"`
	LastFour   string    `json:"lastFour"`
	ExpiryDate string    `json:"expiryDate"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// NewCardResponse creates a new CardResponse from a domain.Card
func NewCardResponse(card *domain.Card) CardResponse {
	return CardResponse{
		UUID:       card.UUID,
		Provider:   card.Provider,
		Type:       card.Type,
		LastFour:   card.LastFour,
		ExpiryDate: formatExpiryDate(card.ExpiryDate),
		Status:     card.Status,
		CreatedAt:  card.CreatedAt,
		UpdatedAt:  card.UpdatedAt,
	}
}

// CardListResponse represents the response body for listing cards.
// @Description CardListResponse includes a list of cards.
type CardListResponse struct {
	Cards []CardResponse `json:"cards"`
}

// NewCardListResponse creates a new CardListResponse from a slice of domain.Card
func NewCardListResponse(cards []*domain.Card) CardListResponse {
	response := CardListResponse{
		Cards: make([]CardResponse, len(cards)),
	}

	for i, card := range cards {
		response.Cards[i] = NewCardResponse(card)
	}

	return response
}

// parseExpiryDate is a helper for checking if the date is in the future and sets the day to the last day of the month
func parseExpiryDate(expiryDate string) (time.Time, error) {
	t, err := time.Parse(common.CardExpiryLayout, expiryDate)
	if err != nil {
		return time.Time{}, err
	}

	lastDay := time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	if lastDay.Before(time.Now()) {
		return time.Time{}, errors.New("expiry date must be in the future")
	}

	return lastDay, nil
}

// formatExpiryDate is ahelper to convert expiry date time.Time to "MM/YY" format
func formatExpiryDate(date time.Time) string {
	return date.Format(common.CardExpiryLayout)
}
