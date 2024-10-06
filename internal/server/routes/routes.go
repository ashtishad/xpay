package routes

import (
	"database/sql"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/gin-gonic/gin"
)

func InitRoutes(rg *gin.RouterGroup, logger *slog.Logger, db *sql.DB, config *common.AppConfig) {
	registerPingRoutes(rg, logger)
}
