package handlers

import (
	"context"
	"net/http"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/dto"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletRepo domain.WalletRepository
	userRepo   domain.UserRepository
}

func NewWalletHandler(walletRepo domain.WalletRepository, userRepo domain.UserRepository) *WalletHandler {
	return &WalletHandler{
		walletRepo: walletRepo,
		userRepo:   userRepo,
	}
}

// CreateWallet godoc
// @Summary Create a new wallet for a user
// @Description Creates a new wallet for the specified user
// @Tags wallet
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param input body dto.CreateWalletRequest true "Wallet creation details"
// @Success 201 {object} dto.CreateWalletResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	authorizedUser, appErr := validateUserAccess(c)
	if appErr != nil {
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	var req dto.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	wallet := req.ToNewWallet(authorizedUser.ID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Wallet.Write)
	defer cancel()

	createdWallet, appErr := h.walletRepo.Create(ctx, wallet)
	if appErr != nil {
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateWalletResponse{Wallet: *createdWallet})
}

// GetWalletBalance godoc
// @Summary Get wallet balance
// @Description Retrieves the balance of a specific wallet for a user
// @Tags wallet
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Success 200 {object} dto.GetWalletBalanceResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/balance [get]
func (h *WalletHandler) GetWalletBalance(c *gin.Context) {
	_, appErr := validateUserAccess(c)
	if appErr != nil {
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	walletUUID := c.Param("wallet_uuid")

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Wallet.Read)
	defer cancel()

	balance, appErr := h.walletRepo.GetBalance(ctx, walletUUID)
	if appErr != nil {
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.GetWalletBalanceResponse{
		BalanceInCents: balance,
		Currency:       domain.WalletCurrencyUSD,
	})
}

// UpdateWalletStatus godoc
// @Summary Update wallet status
// @Description Updates the status of a specific wallet for a user
// @Tags wallet
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_uuid path string true "User UUID"
// @Param wallet_uuid path string true "Wallet UUID"
// @Param input body dto.UpdateWalletStatusRequest true "New wallet status"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{user_uuid}/wallets/{wallet_uuid}/status [patch]
func (h *WalletHandler) UpdateWalletStatus(c *gin.Context) {
	_, appErr := validateUserAccess(c)
	if appErr != nil {
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	walletUUID := c.Param("wallet_uuid")

	var req dto.UpdateWalletStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Wallet.Write)
	defer cancel()

	appErr = h.walletRepo.UpdateStatus(ctx, walletUUID, req.Status)
	if appErr != nil {
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Wallet status updated successfully"})
}
