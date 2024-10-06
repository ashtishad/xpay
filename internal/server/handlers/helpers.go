package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

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
			messages = append(messages, fmt.Sprintf("%s must be a valid email", e.Field()))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param()))
		case "gt":
			messages = append(messages, fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param()))
		case "gte":
			messages = append(messages, fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param()))
		case "lt":
			messages = append(messages, fmt.Sprintf("%s must be less than %s", e.Field(), e.Param()))
		case "lte":
			messages = append(messages, fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param()))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s must be one of [%s]", e.Field(), e.Param()))
		case "numeric":
			messages = append(messages, fmt.Sprintf("%s must be numeric", e.Field()))
		case "alphanum":
			messages = append(messages, fmt.Sprintf("%s must be alphanumeric", e.Field()))
		case "credit_card":
			messages = append(messages, fmt.Sprintf("%s must be a valid credit card number", e.Field()))
		case "uuid":
			messages = append(messages, fmt.Sprintf("%s must be a valid UUID", e.Field()))
		case "len":
			messages = append(messages, fmt.Sprintf("%s must be exactly %s characters long", e.Field(), e.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s failed on tag %s", e.Field(), e.Tag()))
		}
	}

	return strings.Join(messages, ". ")
}
