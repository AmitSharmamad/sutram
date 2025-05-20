// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/category_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"
)

// CategoryRepository defines the operations for category persistence
type CategoryRepository interface {
	// Get all categories
	GetAll(ctx context.Context) ([]*model.Category, error)

	// Get a category by ID
	GetByID(ctx context.Context, id int) (*model.Category, error)

	// Get a category by slug
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)

	// Create a new category
	Create(ctx context.Context, category *model.Category) error

	// Update a category
	Update(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a category
	Delete(ctx context.Context, id int) error

	// Get topic count for a category
	GetTopicCount(ctx context.Context, categoryID int) (int, error)
}
