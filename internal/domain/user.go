package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserStatusActive   string = "active"
	UserStatusInactive string = "inactive"
	UserStatusDeleted  string = "deleted"

	UserRoleAdmin    string = "admin"
	UserRoleUser     string = "user"
	UserRoleAgent    string = "agent"
	UserRoleMerchant string = "merchant"
)

type User struct {
	ID           int64     `json:"-"`
	UUID         uuid.UUID `json:"uuid"`
	FullName     string    `json:"fullName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Status       string    `json:"status"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
