package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware configures and returns a CORS middleware.
// It allows specific origins, methods, and headers for cross-origin requests.
// TODO: Adjust AllowOrigins for production environment.
func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	// if s.Config.Server.AppEnv == common.AppEnvProduction {
	// 	// For production, I will set more restrictive origins
	// 	config.AllowOrigins = []string{"https://api.yourdomain.com", "https://www.yourdomain.com"}
	// }

	return cors.New(config)
}
