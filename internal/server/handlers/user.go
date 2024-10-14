package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/dto"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/ashtishad/xpay/internal/secure/rbac"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo domain.UserRepository
}

func NewUserHandler(userRepo domain.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// CreateUserWithRole godoc
// @Summary Create a new user with a specific role
// @Description Creates a new user with admin, user, agent, or merchant role. Only admins can perform this action.
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param input body dto.CreateUserRequest true "User creation details"
// @Success 201 {object} dto.CreateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUserWithRole(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	authorizedUser, exists := c.Get(common.ContextKeyAuthorizedUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "User not authenticated"})
		return
	}

	user, ok := authorizedUser.(*domain.User)
	if !ok {
		slog.Error("failed to cast authorized user")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Unexpected server error"})
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request body", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	if !rbac.CanCreateUser(user.Role, req.Role) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "You don't have permission to create a user with this role"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.User.Write)
	defer cancel()

	passwordHash, err := secure.GeneratePasswordHash(req.Password)
	if err != nil {
		slog.Error("failed to generate password hash", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: common.ErrUnexpectedServer})
		return
	}

	newUser := req.ToUser(passwordHash)
	createdUser, appErr := h.userRepo.Create(ctx, newUser)
	if appErr != nil {
		slog.Error("failed to create user", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateUserResponse{
		User: *createdUser,
	})
}
