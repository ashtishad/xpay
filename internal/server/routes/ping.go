package routes

import (
	"log/slog"

	"github.com/ashtishad/xpay/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func registerPingRoutes(rg *gin.RouterGroup, l *slog.Logger) {
	pingHandler := handlers.NewPingHandler(l)

	rg.GET("/ping", pingHandler.PingHandler)
}
