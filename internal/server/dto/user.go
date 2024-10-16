package dto

import (
	"time"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FullName string `json:"fullName" binding:"required,min=3,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Role     string `json:"role" binding:"required,oneof=admin user agent merchant"`
}

func (r *CreateUserRequest) ToUser(passwordHash string) *domain.User {
	now := time.Now().UTC()
	return &domain.User{
		UUID:         uuid.New(),
		FullName:     r.FullName,
		Email:        r.Email,
		PasswordHash: passwordHash,
		Status:       domain.UserStatusActive,
		Role:         r.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type CreateUserResponse struct {
	User domain.User `json:"user"`
}
