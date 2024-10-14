package routes

import (
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func registerWalletRoutes(rg *gin.RouterGroup, walletRepo domain.WalletRepository, userRepo domain.UserRepository) {
	walletHandler := handlers.NewWalletHandler(walletRepo, userRepo)

	rg.POST("/:user_uuid/wallets", walletHandler.CreateWallet)
	rg.GET("/:user_uuid/wallets/:wallet_uuid/balance", walletHandler.GetWalletBalance)
	rg.PATCH("/:user_uuid/wallets/:wallet_uuid/status", walletHandler.UpdateWalletStatus)
}
