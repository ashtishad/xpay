package middlewares

import (
	"crypto/ecdsa"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/dto"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token with ecdsa public key from the Authorization header,
// extracts the user uuid from claims, and fetches the corresponding user from the database.
// The authorizedUser is then stored in the request context for use in subsequent handlers.
func AuthMiddleware(userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Missing authorization header"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := secure.ValidateToken(tokenString, jwtPublicKey)
		if err != nil {
			slog.Error("failed to velidate token", "err", err)
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

		// Set the authorized user in the context
		c.Set(common.ContextKeyAuthorizedUser, user)
		c.Next()
	}
}

// AdminAuthMiddleware ensures that only users with admin role can access the protected routes.
// It builds upon the AuthMiddleware by performing these additional steps:
// 1. Runs the standard AuthMiddleware to authenticate the user.
// 2. Retrieves the authenticated user from the gin context.
// 3. Verifies if the authenticated user has the admin role.
// 4. Aborts the request if the user is not an admin, otherwise allows it to proceed.
func AdminAuthMiddleware(userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthMiddleware(userRepo, jwtPublicKey)(c)

		if c.IsAborted() {
			return
		}

		user, exists := c.Get(common.ContextKeyAuthorizedUser)
		if !exists {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "User not found in context"})
			c.Abort()
			return
		}

		authorizedUser, ok := user.(*domain.User)
		if !ok || authorizedUser.Role != domain.UserRoleAdmin {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
