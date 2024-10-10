package routes

import (
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func registerUserManagementRoutes(rg *gin.RouterGroup, userRepo domain.UserRepository) {
    userHandler := handlers.NewUserHandler(userRepo)
    rg.POST("", userHandler.CreateUserWithRole)
}
