package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// InitMiddlewares configures and returns all middleware functions.
// It includes:
// - Panic recovery
// - Custom logger (in debug mode)
// - CORS handling
// - Request ID generation and propagation
func InitMiddlewares(logger *slog.Logger) []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{
		gin.Recovery(),
	}

	if gin.IsDebugging() {
		middlewares = append(middlewares, CustomLogger(logger))
	}

	middlewares = append(middlewares,
		CorsMiddleware(),
		RequestID(),
	)

	return middlewares
}
