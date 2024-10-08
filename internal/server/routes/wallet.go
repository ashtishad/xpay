package routes

import (
	"crypto/ecdsa"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/ashtishad/xpay/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func registerWalletRoutes(rg *gin.RouterGroup, walletRepo domain.WalletRepository, userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey) {
	walletHandler := handlers.NewWalletHandler(walletRepo, userRepo)

	users := rg.Group("/users")
	users.Use(middlewares.AuthMiddleware(userRepo, jwtPublicKey))
	{
		wallets := users.Group("/:user_uuid/wallets")
		{
			wallets.POST("", walletHandler.CreateWallet)
			wallets.GET("/:wallet_uuid/balance", walletHandler.GetWalletBalance)
			wallets.PATCH("/:wallet_uuid/status", walletHandler.UpdateWalletStatus)
		}
	}
}
