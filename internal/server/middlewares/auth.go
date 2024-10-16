package middlewares

import (
	"crypto/ecdsa"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/ashtishad/xpay/internal/secure/rbac"
	"github.com/ashtishad/xpay/internal/server/dto"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens, authenticates users, enforces RBAC policies,
// and sets the authenticated user in the request context for subsequent handlers.
func AuthMiddleware(userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey, rbac *rbac.RBAC) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(common.AuthorizationHeaderKey)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Missing authorization header"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != common.TokenTypeBearer {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := secure.ValidateToken(tokenString, jwtPublicKey)
		if err != nil {
			slog.Error("failed to validate token", "err", err)
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid or expired token"})
			c.Abort()
			return
		}

		user, appErr := userRepo.FindBy(c.Request.Context(), common.DBColumnUUID, claims.UserUUID)
		if appErr != nil {
			slog.Error("failed to find user from jwt claims user uuid", "err", err)
			c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
			c.Abort()
			return
		}

		if !rbac.HasPermission(user.Role, c.FullPath(), c.Request.Method) {
			slog.Warn("access denied", "role", user.Role, "method", c.Request.Method, "path", c.FullPath())
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Access denied"})
			c.Abort()
			return
		}

		c.Set(common.ContextKeyAuthorizedUser, user)
		c.Next()
	}
}
