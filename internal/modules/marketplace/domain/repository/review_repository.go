// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/repository/review_repository.go
package repository

import (
	"context"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// ReviewRepository defines the interface for marketplace review repository
type ReviewRepository interface {
	// Create creates a new review
	Create(ctx context.Context, review *model.Review) error

	// FindByID returns a review by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Review, error)

	// FindByListingID returns all reviews for a listing
	FindByListingID(ctx context.Context, listingID uuid.UUID) ([]*model.Review, error)

	// FindByUserAndListingID returns a review by user and listing
	FindByUserAndListingID(ctx context.Context, userID uuid.UUID, listingID uuid.UUID) (*model.Review, error)

	// Update updates a review
	Update(ctx context.Context, review *model.Review) error

	// Delete deletes a review
	Delete(ctx context.Context, id uuid.UUID) error

	// CalculateAverageRating calculates the average rating for a listing
	CalculateAverageRating(ctx context.Context, listingID uuid.UUID) (float32, error)
}
