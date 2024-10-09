package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// validateUserAccess checks if the authenticated user has permission to access the requested resource.
// It verifies:
// 1. The presence of the user_uuid route parameter
// 2. The existence of an authenticated user in the context
// 3. The authenticated user's UUID matches the requested user_uuid
// Returns the authorized user on success, or an appropriate AppError on failure.
func validateUserAccess(c *gin.Context) (*domain.User, common.AppError) {
	userUUID := c.Param("user_uuid")
	if userUUID == "" {
		return nil, common.NewBadRequestError("User UUID route param is required")
	}

	authUser, exists := c.Get(common.ContextKeyAuthorizedUser)
	if !exists {
		return nil, common.NewUnauthorizedError("User not authenticated")
	}

	authorizedUser, ok := authUser.(*domain.User)
	if !ok {
		slog.Error("failed to cast authorized user")
		return nil, common.NewInternalServerError("Unexpected server error", nil)
	}

	if authorizedUser.UUID.String() != userUUID {
		return nil, common.NewForbiddenError("You can only access your own resources")
	}

	return authorizedUser, nil
}

// formatValidationError formats validation errors into a single string
func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err.Error()
	}

	messages := make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		switch e.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", e.Field()))
		case "email":
			messages = append(messages, fmt.Sprintf("%s must be a valid email address", e.Field()))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param()))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s must be one of [%s]", e.Field(), e.Param()))
		case "creditcard":
			messages = append(messages, fmt.Sprintf("%s must be a valid credit card number", e.Field()))
		case "gtfield":
			messages = append(messages, fmt.Sprintf("%s must be after %s", e.Field(), e.Param()))
		case "datetime":
			messages = append(messages, fmt.Sprintf("%s must be a valid date and time", e.Field()))
		case "credit_card":
			messages = append(messages, fmt.Sprintf("%s must be a valid credit card number", e.Field()))
		default:
			messages = append(messages, fmt.Sprintf("%s failed validation on tag %s", e.Field(), e.Tag()))
		}
	}

	return strings.Join(messages, ". ")
}
