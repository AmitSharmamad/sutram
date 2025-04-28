// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/post_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// PostRepository defines the operations for post persistence
type PostRepository interface {
	// Create a new post
	Create(ctx context.Context, post *model.Post) error

	// Get a post by ID
	GetByID(ctx context.Context, id int) (*model.Post, error)

	// Get posts for a topic
	GetByTopicID(ctx context.Context, topicID int, page, perPage int) ([]*model.Post, int, error)

	// Get posts by a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]*model.Post, int, error)

	// Update a post
	Update(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a post
	Delete(ctx context.Context, id int) error

	// Get like count for a post
	GetLikeCount(ctx context.Context, postID int) (int, error)

	// Check if a user is the author of a post
	IsAuthor(ctx context.Context, postID int, userID uuid.UUID) (bool, error)
}
