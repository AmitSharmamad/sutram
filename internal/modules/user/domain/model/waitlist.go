package model

import (
	"time"
)

// Waitlist represents a waitlist entry
type Waitlist struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewWaitlist creates a new waitlist entry
func NewWaitlist(email string) *Waitlist {
	now := time.Now()
	return &Waitlist{
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
