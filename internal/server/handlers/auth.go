package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/ashtishad/xpay/internal/dto"
	"github.com/ashtishad/xpay/internal/secure"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo   domain.UserRepository
	jwtManager *secure.JWTManager
}

func NewAuthHandler(userRepo domain.UserRepository, jm *secure.JWTManager) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jm,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Hashes password using bcrypt before storage.
// @Description Generates JWT access token using ECDSA encryption.
// @Description Sets HTTP-only cookie with access token and X-Request-Id header.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterUserRequest true "User registration details"
// @Success 201 {object} RegisterUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request body", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Auth.Write)
	defer cancel()

	passwordHash, err := secure.GeneratePasswordHash(req.Password)
	if err != nil {
		slog.Error("failed to generate password hash", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: common.ErrUnexpectedServer})
		return
	}

	createdUser, appErr := h.userRepo.Create(ctx, req.ToUser(passwordHash))
	if appErr != nil {
		slog.Error("failed to create user", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(createdUser.UUID.String())
	if err != nil {
		slog.Error("failed to generate access token", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: common.ErrUnexpectedServer})
		return
	}

	c.SetCookie("accessToken", accessToken, int(h.jwtManager.AccessExpiration.Seconds()), "/", "", true, true)

	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		User: *createdUser,
	})
}

// Login godoc
// @Summary Authenticate a user and provide access tokens
// @Description Verifies password using bcrypt comparison.
// @Description Generates new JWT access token using ECDSA encryption.
// @Description Sets HTTP-only cookie with new access token and and X-Request-Id header.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	requestID := c.GetString(common.ContextKeyRequestID)
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request body", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: formatValidationError(err)})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), common.Timeouts.Auth.Read)
	defer cancel()

	user, appErr := h.userRepo.FindBy(ctx, common.DBColumnEmail, req.Email)
	if appErr != nil {
		slog.Error("failed to find user", "requestID", requestID, "error", appErr.Error())
		c.JSON(appErr.Code(), dto.ErrorResponse{Error: appErr.Error()})
		return
	}

	if err := secure.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		slog.Error("invalid credentials", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(user.UUID.String())
	if err != nil {
		slog.Error("failed to generate access token", "requestID", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: common.ErrUnexpectedServer})
		return
	}

	c.SetCookie("accessToken", accessToken, int(h.jwtManager.AccessExpiration.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, dto.LoginResponse{
		User: *user,
	})
}
