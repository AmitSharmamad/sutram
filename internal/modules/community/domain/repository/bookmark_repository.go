// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/bookmark_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// BookmarkRepository defines the operations for bookmark persistence
type BookmarkRepository interface {
	// Create a new bookmark
	Create(ctx context.Context, bookmark *model.Bookmark) error

	// Delete a bookmark
	Delete(ctx context.Context, userID uuid.UUID, topicID int) error

	// Check if a bookmark exists
	Exists(ctx context.Context, userID uuid.UUID, topicID int) (bool, error)

	// Get bookmarks by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]*model.Topic, int, error)

	// Get bookmark count for a topic
	GetBookmarkCount(ctx context.Context, topicID int) (int, error)
}
