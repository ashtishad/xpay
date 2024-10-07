package middlewares

import (
	"github.com/ashtishad/xpay/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID middleware generates or propagates a unique ID for each request.
// If a request ID is provided in the header, it uses that; otherwise, it generates a new UUID.
// The ID is set in both the request context and response header for tracing purposes.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(common.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(common.RequestIDKey, requestID)
		c.Header(common.RequestIDHeader, requestID)
		c.Next()
	}
}
