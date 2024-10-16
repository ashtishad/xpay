package dto

import (
	"time"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/google/uuid"
)

// RegisterUserRequest holds the data for creating a new user.
// @Description RegisterUserRequest validates input for user registration.
// @Description FullName must be at least 3 and at max 255 characters long.
// @Description Email must be a valid email address.
// @Description Password must be at least 8 and at max 64 characters long.
type RegisterUserRequest struct {
	FullName string `json:"fullName" binding:"required,min=3,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// ToUser converts RegisterUserRequest to domain.User
func (r *RegisterUserRequest) ToUser(passwordHash string) *domain.User {
	now := time.Now().UTC()
	return &domain.User{
		UUID:         uuid.New(),
		FullName:     r.FullName,
		Email:        r.Email,
		PasswordHash: passwordHash,
		Status:       domain.UserStatusActive,
		Role:         domain.UserRoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// RegisterUserResponse contains the user data returned after successful registration.
// @Description RegisterUserResponse includes the created user's details.
type RegisterUserResponse struct {
	User domain.User `json:"user"`
}

// LoginRequest represents the credentials required for user authentication.
// @Description LoginRequest validates input for user login.
// @Description Email must be a valid email address.
// @Description Password is required.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse contains the user data returned after successful login.
// @Description LoginResponse includes the authenticated user's details.
type LoginResponse struct {
	User domain.User `json:"user"`
}
