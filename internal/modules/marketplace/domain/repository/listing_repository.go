// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/repository/listing_repository.go
package repository

import (
	"context"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// ListingRepository defines the interface for marketplace listing repository
type ListingRepository interface {
	// Create creates a new marketplace listing
	Create(ctx context.Context, listing *model.Listing) error

	// FindByID returns a listing by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Listing, error)

	// FindAll returns all listings with optional filtering and pagination
	FindAll(ctx context.Context, filters map[string]interface{}, page, perPage int) ([]*model.Listing, int, error)

	// Update updates a listing
	Update(ctx context.Context, listing *model.Listing) error

	// Delete deletes a listing
	Delete(ctx context.Context, id uuid.UUID) error

	// IncrementViews increments the view count for a listing
	IncrementViews(ctx context.Context, id uuid.UUID) error

	// FindByUserID returns all listings created by a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Listing, error)

	// AddCategories adds categories to a listing
	AddCategories(ctx context.Context, listingID uuid.UUID, categoryIDs []uuid.UUID) error

	// RemoveCategories removes categories from a listing
	RemoveCategories(ctx context.Context, listingID uuid.UUID, categoryIDs []uuid.UUID) error

	// GetCategories gets all categories for a listing
	GetCategories(ctx context.Context, listingID uuid.UUID) ([]model.Category, error)

	// SetFeatures sets the feature list for a listing
	SetFeatures(ctx context.Context, listingID uuid.UUID, features []string) error

	// GetFeatures gets the feature list for a listing
	GetFeatures(ctx context.Context, listingID uuid.UUID) ([]string, error)
}
