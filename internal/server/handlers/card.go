package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/dto"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardRepo      domain.CardRepository
	walletRepo    domain.WalletRepository
	cardEncryptor *secure.CardEncryptor
}

func NewCardHandler(cardRepo domain.CardRepository, walletRepo domain.WalletRepository, cardEncryptor *secure.CardEncryptor) *CardHandler {
	return &CardHandler{
		cardRepo:      cardRepo,
		walletRepo:    walletRepo,
		cardEncryptor: cardEncryptor,
	}
}

// AddCardToWallet godoc
// @Summary Add a new card to a wallet
// @Description Adds a new card to the specified wallet, encrypting sensitive data
// @Tags card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Param input body dto.AddCardRequest true "Card details"
// @Success 201 {object} dto.AddCardResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/cards [post]
func (h *CardHandler) AddCardToWallet(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	authorizedUser, appErr := validateUserAccess(c)
	if appErr != nil {
		slog.Error("failed to validate user access", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Card.Write)
	defer cancel()

	walletID, appErr := h.getWalletID(ctx, c.Param("wallet_uuid"))
	if appErr != nil {
		slog.Error("failed to get wallet ID", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	var req dto.AddCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request body", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	encryptedCardNumber, err := h.cardEncryptor.Encrypt(req.CardNumber)
	if err != nil {
		slog.Error("failed to encrypt card number", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to process card data"})
		return
	}

	card, err := req.ToCard(authorizedUser.ID, walletID, encryptedCardNumber)
	if err != nil {
		slog.Error("failed to create card object", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	createdCard, appErr := h.cardRepo.AddCardToWallet(ctx, card)
	if appErr != nil {
		slog.Error("failed to add card to wallet", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.AddCardResponse{Card: dto.NewCardResponse(createdCard)})
}

// GetCard godoc
// @Summary Get card details
// @Description Retrieves details of a specific card
// @Tags card
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Param card_uuid path string true "Card UUID"
// @Success 200 {object} dto.CardResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid} [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	_, appErr := validateUserAccess(c)
	if appErr != nil {
		slog.Error("failed to validate user access", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	cardUUID := c.Param("card_uuid")
	if cardUUID == "" {
		appErr := common.NewBadRequestError("Card UUID is required")
		slog.Error("missing card UUID", "requestID", requestID)
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Card.Read)
	defer cancel()

	card, appErr := h.cardRepo.FindBy(ctx, common.DBColumnUUID, cardUUID)
	if appErr != nil {
		slog.Error("failed to find card", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CardResponse(dto.NewCardResponse(card)))
}

// UpdateCard godoc
// @Summary Update card details
// @Description Updates the details of a specific card
// @Tags card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Param card_uuid path string true "Card UUID"
// @Param input body dto.UpdateCardRequest true "Updated card details"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid} [patch]
func (h *CardHandler) UpdateCard(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	_, appErr := validateUserAccess(c)
	if appErr != nil {
		slog.Error("failed to validate user access", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	cardUUID := c.Param("card_uuid")
	if cardUUID == "" {
		appErr := common.NewBadRequestError("Card UUID is required")
		slog.Error("missing card UUID", "requestID", requestID)
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	var req dto.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request body", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Card.Write)
	defer cancel()

	card, appErr := h.cardRepo.FindBy(ctx, common.DBColumnUUID, cardUUID)
	if appErr != nil {
		slog.Error("failed to find card", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	updatedCard, err := req.UpdateCard(card)
	if err != nil {
		slog.Error("failed to update card", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	appErr = h.cardRepo.Update(ctx, updatedCard)
	if appErr != nil {
		slog.Error("failed to save updated card", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Card updated successfully"})
}

// DeleteCard godoc
// @Summary Delete a card
// @Description Soft deletes a specific card
// @Tags card
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Param card_uuid path string true "Card UUID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid} [delete]
func (h *CardHandler) DeleteCard(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	_, appErr := validateUserAccess(c)
	if appErr != nil {
		slog.Error("failed to validate user access", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	cardUUID := c.Param("card_uuid")
	if cardUUID == "" {
		appErr := common.NewBadRequestError("Card UUID is required")
		slog.Error("missing card UUID", "requestID", requestID)
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Card.Write)
	defer cancel()

	appErr = h.cardRepo.Delete(ctx, cardUUID)
	if appErr != nil {
		slog.Error("failed to delete card", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Card deleted successfully"})
}

// ListCards godoc
// @Summary List cards
// @Description Retrieves a list of cards for a specific wallet
// @Tags card
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Param provider query string false "Filter by card provider"
// @Param status query string false "Filter by card status"
// @Success 200 {object} dto.CardListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/cards [get]
func (h *CardHandler) ListCards(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	authorizedUser, appErr := validateUserAccess(c)
	if appErr != nil {
		slog.Error("failed to validate user access", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Card.Read)
	defer cancel()

	walletID, appErr := h.getWalletID(ctx, c.Param("wallet_uuid"))
	if appErr != nil {
		slog.Error("failed to get wallet ID", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	filters := domain.CardFilters{
		UserID:   &authorizedUser.ID,
		WalletID: &walletID,
	}

	if provider := c.Query("provider"); provider != "" {
		filters.Provider = &provider
	}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	cards, appErr := h.cardRepo.List(ctx, filters)
	if appErr != nil {
		slog.Error("failed to list cards", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.NewCardListResponse(cards))
}

// Helper function to get wallet ID from UUID
func (h *CardHandler) getWalletID(ctx context.Context, walletUUID string) (int64, common.AppError) {
	if walletUUID == "" {
		return 0, common.NewBadRequestError("Wallet UUID is required")
	}

	return h.walletRepo.FindIDFromUUID(ctx, walletUUID)
}
