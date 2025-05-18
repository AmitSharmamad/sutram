// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/like_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// LikeRepository defines the operations for like persistence
type LikeRepository interface {
	// Create a new like
	Create(ctx context.Context, like *model.Like) error

	// Delete a like
	Delete(ctx context.Context, userID uuid.UUID, likeableType string, likeableID int) error

	// Check if a like exists
	Exists(ctx context.Context, userID uuid.UUID, likeableType string, likeableID int) (bool, error)

	// Get like count for a likeable object
	GetLikeCount(ctx context.Context, likeableType string, likeableID int) (int, error)

	// Get likes by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]*model.Like, int, error)
}
