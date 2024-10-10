package routes

import (
	"crypto/ecdsa"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/ashtishad/xpay/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func registerUserManagementRoutes(rg *gin.RouterGroup, userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey) {
	userHandler := handlers.NewUserHandler(userRepo)

	users := rg.Group("/users")
	users.Use(middlewares.AuthMiddleware(userRepo, jwtPublicKey))
	{
		adminRoutes := users.Group("")
		adminRoutes.Use(middlewares.AdminAuthMiddleware(userRepo, jwtPublicKey))
		{
			adminRoutes.POST("", userHandler.CreateUserWithRole) // Only admins can create new users (admin, agent, merchant)
		}
	}
}
