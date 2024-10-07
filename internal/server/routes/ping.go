package routes

import (
	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func registerPingRoutes(rg *gin.RouterGroup) {
	pingHandler := handlers.NewPingHandler()

	rg.GET("/ping", pingHandler.PingHandler)
}
