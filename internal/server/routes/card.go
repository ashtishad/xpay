package routes

import (
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func registerCardRoutes(rg *gin.RouterGroup, cardRepo domain.CardRepository, walletRepo domain.WalletRepository, cardEncryptor *secure.CardEncryptor) {
    cardHandler := handlers.NewCardHandler(cardRepo, walletRepo, cardEncryptor)

    cards := rg.Group("/:user_uuid/wallets/:wallet_uuid/cards")
    {
        cards.POST("", cardHandler.AddCardToWallet)
        cards.GET("/:card_uuid", cardHandler.GetCard)
        cards.PATCH("/:card_uuid", cardHandler.UpdateCard)
        cards.DELETE("/:card_uuid", cardHandler.DeleteCard)
        cards.GET("", cardHandler.ListCards)
    }
}
