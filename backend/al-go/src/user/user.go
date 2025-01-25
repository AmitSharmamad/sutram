package user

import "time"

type User struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Password    string    `json:"password" db:"password"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}
