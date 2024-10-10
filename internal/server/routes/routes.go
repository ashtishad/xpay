package routes

import (
	"database/sql"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/ashtishad/xpay/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(rg *gin.RouterGroup, db *sql.DB, config *common.AppConfig, jm *secure.JWTManager, cardEncryptor *secure.CardEncryptor) {
	userRepo := domain.NewUserRepository(db)
	walletRepo := domain.NewWalletRepository(db)
	cardRepo := domain.NewCardRepository(db)

	// Register public routes
	registerAuthRoutes(rg, userRepo, jm)

	// Create authenticated user gin router group
	authGroup := rg.Group("/users")
	authGroup.Use(middlewares.AuthMiddleware(userRepo, jm.GetPublicKey()))

	// Register authenticated routes
	registerUserManagementRoutes(authGroup, userRepo)
	registerWalletRoutes(authGroup, walletRepo, userRepo)
	registerCardRoutes(authGroup, cardRepo, walletRepo, cardEncryptor)
}
