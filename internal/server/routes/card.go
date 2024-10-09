package routes

import (
	"crypto/ecdsa"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/ashtishad/xpay/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func registerCardRoutes(rg *gin.RouterGroup, cardRepo domain.CardRepository, walletRepo domain.WalletRepository, userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey, cardEncryptor *secure.CardEncryptor) {
	cardHandler := handlers.NewCardHandler(cardRepo, walletRepo, cardEncryptor)

	users := rg.Group("/users")
	users.Use(middlewares.AuthMiddleware(userRepo, jwtPublicKey))
	{
		wallets := users.Group("/:user_uuid/wallets")
		{
			cards := wallets.Group("/:wallet_uuid/cards")
			{
				cards.POST("", cardHandler.AddCardToWallet)
				cards.GET("/:card_uuid", cardHandler.GetCard)
				cards.PATCH("/:card_uuid", cardHandler.UpdateCard)
				cards.DELETE("/:card_uuid", cardHandler.DeleteCard)
				cards.GET("", cardHandler.ListCards)
			}
		}
	}
}
