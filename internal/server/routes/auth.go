package routes

import (
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup, userRepo domain.UserRepository, jm *secure.JWTManager) {
	authHandler := handlers.NewAuthHandler(userRepo, jm)

	rg.POST("/register", authHandler.Register)
	rg.POST("/login", authHandler.Login)
}
