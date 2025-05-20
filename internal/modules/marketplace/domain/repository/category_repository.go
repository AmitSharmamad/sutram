// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/repository/category_repository.go
package repository

import (
	"context"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// CategoryRepository defines the interface for marketplace category repository
type CategoryRepository interface {
	// FindAll returns all categories
	FindAll(ctx context.Context) ([]model.Category, error)

	// FindByID returns a category by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error)

	// FindByName returns a category by its name
	FindByName(ctx context.Context, name string) (*model.Category, error)

	// FindByIDs returns categories by their IDs
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Category, error)
}
