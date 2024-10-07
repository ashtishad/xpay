package routes

import (
	"database/sql"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/gin-gonic/gin"
)

func InitRoutes(rg *gin.RouterGroup, db *sql.DB, config *common.AppConfig) {
	registerPingRoutes(rg)
}
