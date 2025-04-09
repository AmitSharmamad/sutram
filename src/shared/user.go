package shared

import (
	"github.com/google/uuid"
)

type Users []User

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Password    string    `json:"password" db:"password"`
	Verified    bool      `json:"verified" db:"verified"`
	IsActive    bool      `json:"is_active" db:"is_active"`

	Audit
}

type UserService interface {
	GetByUsername(username string) (User, error)
	Create(user *User) (bool, error)
}
