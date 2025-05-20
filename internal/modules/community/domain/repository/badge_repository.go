// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/badge_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// BadgeRepository defines the operations for badge persistence
type BadgeRepository interface {
	// Get all badges
	GetAll(ctx context.Context) ([]*model.Badge, error)

	// Get a badge by ID
	GetByID(ctx context.Context, id int) (*model.Badge, error)

	// Create a new badge
	Create(ctx context.Context, badge *model.Badge) error

	// Update a badge
	Update(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a badge
	Delete(ctx context.Context, id int) error

	// Award a badge to a user
	AwardBadge(ctx context.Context, userID uuid.UUID, badgeID int) error

	// Check if a user has a badge
	HasBadge(ctx context.Context, userID uuid.UUID, badgeID int) (bool, error)

	// Get badges for a user
	GetUserBadges(ctx context.Context, userID uuid.UUID) ([]*model.Badge, error)

	// Revoke a badge from a user
	RevokeBadge(ctx context.Context, userID uuid.UUID, badgeID int) error
}
