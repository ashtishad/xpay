package routes

import (
	"database/sql"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/gin-gonic/gin"
)

func InitRoutes(rg *gin.RouterGroup, db *sql.DB, config *common.AppConfig, jm *secure.JWTManager) {
	userRepo := domain.NewUserRepository(db)
	walletRepo := domain.NewWalletRepository(db)

	registerAuthRoutes(rg, userRepo, jm)
	registerWalletRoutes(rg, walletRepo, userRepo, jm.GetPublicKey())
}
