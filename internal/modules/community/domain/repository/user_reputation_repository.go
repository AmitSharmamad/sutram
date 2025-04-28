// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/user_reputation_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// UserReputationRepository defines the operations for user reputation persistence
type UserReputationRepository interface {
	// Get user reputation by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error)

	// Create or update user reputation
	CreateOrUpdate(ctx context.Context, userProfile *model.UserProfile) error

	// Increment reputation score
	IncrementReputationScore(ctx context.Context, userID uuid.UUID, amount int) error

	// Increment contribution count
	IncrementContributionCount(ctx context.Context, userID uuid.UUID) error

	// Update last active timestamp
	UpdateLastActive(ctx context.Context, userID uuid.UUID) error

	// Get topic count for a user
	GetTopicCount(ctx context.Context, userID uuid.UUID) (int, error)

	// Get post count for a user
	GetPostCount(ctx context.Context, userID uuid.UUID) (int, error)
}
