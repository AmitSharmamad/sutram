// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/tag_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"
)

// TagRepository defines the operations for tag persistence
type TagRepository interface {
	// Get all tags
	GetAll(ctx context.Context) ([]*model.Tag, error)

	// Get a tag by ID
	GetByID(ctx context.Context, id int) (*model.Tag, error)

	// Get a tag by slug
	GetBySlug(ctx context.Context, slug string) (*model.Tag, error)

	// Create a new tag
	Create(ctx context.Context, tag *model.Tag) error

	// Update a tag
	Update(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a tag
	Delete(ctx context.Context, id int) error

	// Get tags for a topic
	GetTagsByTopicID(ctx context.Context, topicID int) ([]*model.Tag, error)
}
