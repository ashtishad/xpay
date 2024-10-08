package common

import (
	"fmt"
	"net/http"
)

// AppError is the interface for structured error handling.
// Error() Returns a user-friendly error message
// Code()  Returns the HTTP status code
// DetailedError() Returns a detailed error message for logging
// Wrap(err error) Wraps an internal error
type AppError interface {
	Error() string
	Code() int
	DetailedError() string
	Wrap(err error) AppError
}

// appErr represents the application-specific error structure.
type appErr struct {
	userMessage string
	statusCode  int
	internalErr error
}

// Error returns a user-friendly error message.
//
// Example:
//
//	err := NewBadRequestError("Invalid input")
//	fmt.Println(err.Error()) // Prints: Invalid input
func (e *appErr) Error() string {
	return e.userMessage
}

// Code returns the HTTP status code associated with the error.
//
// Example:
//
//	err := NewBadRequestError("Invalid input")
//	fmt.Println(err.Code()) // Prints: 400
func (e *appErr) Code() int {
	return e.statusCode
}

// DetailedError returns a detailed error message for logging purposes.
//
// Example:
//
//	err := NewInternalServerError("Database query failed", dbErr)
//	log.Println(err.DetailedError()) // Logs: Database query failed: <details of dbErr>
func (e *appErr) DetailedError() string {
	if e.internalErr != nil {
		return fmt.Sprintf("%s: %v", e.userMessage, e.internalErr)
	}
	return e.userMessage
}

// Wrap adds an internal error to the AppError.
//
// Example:
//
//	err := NewBadRequestError("Invalid JSON").Wrap(jsonErr)
func (e *appErr) Wrap(err error) AppError {
	e.internalErr = err
	return e
}

// newAppError creates a new appErr instance.
func newAppError(statusCode int, userMessage string) *appErr {
	return &appErr{
		userMessage: userMessage,
		statusCode:  statusCode,
	}
}

// NewBadRequestError creates a new AppError for bad request errors.
//
// Example:
//
//	err := NewBadRequestError("Invalid input parameters")
func NewBadRequestError(message string) AppError {
	return newAppError(http.StatusBadRequest, message)
}

// NewInternalServerError creates a new AppError for internal server errors.
//
// Example:
//
//	err := NewInternalServerError("Failed to query database", dbErr)
//	slog.Error(err.DetailedError()) // Log the detailed error
func NewInternalServerError(message string, err error) AppError {
	appErr := &appErr{
		userMessage: message,
		statusCode:  http.StatusInternalServerError,
	}

	if err != nil {
		appErr.internalErr = fmt.Errorf("%s: %w", message, err)
	}

	return appErr
}

// NewNotFoundError creates a new AppError for not found errors.
//
// Example:
//
//	err := NewNotFoundError("User not found")
func NewNotFoundError(message string) AppError {
	return newAppError(http.StatusNotFound, message)
}

// NewUnauthorizedError creates a new AppError for unauthorized requests.
//
// Example:
//
//	err := NewUnauthorizedError("Invalid token")
func NewUnauthorizedError(message string) AppError {
	return newAppError(http.StatusUnauthorized, message)
}

// NewConflictError creates a new AppError for conflict errors.
//
// Example:
//
//	err := NewConflictError("User already exists")
func NewConflictError(message string) AppError {
	return newAppError(http.StatusConflict, message)
}

// NewRateLimitError creates a new AppError for rati limit errors.
//
// Example:
//
//	err := NewRateLimitError("too many requests")
func NewRateLimitError(message string) AppError {
	return newAppError(http.StatusTooManyRequests, message)
}

// NewForbiddenError creates a new AppError for forbidden access errors.
//
// Example:
//
//	err := NewForbiddenError("Access denied")
func NewForbiddenError(message string) AppError {
	return newAppError(http.StatusForbidden, message)
}
