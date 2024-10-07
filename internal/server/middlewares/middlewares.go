package middlewares

import (
	"github.com/gin-gonic/gin"
)

// InitMiddlewares configures and returns all middleware functions.
// It includes:
// - Panic recovery
// - Custom logger (in debug mode)
// - CORS handling
// - Request ID generation and propagation
func InitMiddlewares() []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{
		gin.Recovery(),
	}

	if gin.IsDebugging() {
		middlewares = append(middlewares, CustomLogger())
	}

	middlewares = append(middlewares,
		CorsMiddleware(),
		RequestID(),
	)

	return middlewares
}
