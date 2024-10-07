package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomLogger creates a middleware for logging HTTP requests.
// It captures request details, timing, and response status for each request.
// Logs are formatted using structured logging via slog.
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		slog.Info("HTTP Request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
			slog.String("error", errorMessage),
		)
	}
}
