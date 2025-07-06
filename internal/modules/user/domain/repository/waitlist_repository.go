package repository

import (
	"context"

	"algo/internal/modules/user/domain/model"
)

// WaitlistRepository defines the interface for waitlist operations
type WaitlistRepository interface {
	// Create a new waitlist entry
	Create(ctx context.Context, waitlist *model.Waitlist) error

	// FindByEmail finds a waitlist entry by email
	FindByEmail(ctx context.Context, email string) (*model.Waitlist, error)

	// FindAll retrieves all waitlist entries
	FindAll(ctx context.Context) ([]model.Waitlist, error)

	// GetCount returns the total number of waitlist entries
	GetCount(ctx context.Context) (int64, error)
}
