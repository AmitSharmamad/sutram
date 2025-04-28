// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/service/review_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// Additional errors
var (
	ErrReviewNotFound     = errors.New("review not found")
	ErrAlreadyReviewed    = errors.New("user has already reviewed this listing")
	ErrReviewNotPermitted = errors.New("user has not purchased this listing")
)

// ReviewService handles marketplace review business logic
type ReviewService struct {
	reviewRepo  repository.ReviewRepository
	listingRepo repository.ListingRepository
}

// NewReviewService creates a new review service
func NewReviewService(
	reviewRepo repository.ReviewRepository,
	listingRepo repository.ListingRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		listingRepo: listingRepo,
	}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(
	ctx context.Context,
	userID uuid.UUID,
	listingID uuid.UUID,
	rating int,
	comment string,
) (*model.Review, error) {
	// Check if listing exists
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return nil, ErrListingNotFound
	}

	// Check if user has already reviewed this listing
	existingReview, err := s.reviewRepo.FindByUserAndListingID(ctx, userID, listingID)
	if err != nil && !errors.Is(err, ErrReviewNotFound) {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if existingReview != nil {
		return nil, ErrAlreadyReviewed
	}

	// Create the review
	now := time.Now()
	review := &model.Review{
		ID:        uuid.New(),
		ListingID: listingID,
		UserID:    userID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the review
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Update the listing's rating
	avgRating, err := s.reviewRepo.CalculateAverageRating(ctx, listingID)
	if err == nil {
		listing.Ratings = avgRating
		if err := s.listingRepo.Update(ctx, listing); err != nil {
			// Log error but continue, as the review was created successfully
			fmt.Printf("failed to update listing rating: %v\n", err)
		}
	}

	return review, nil
}

// GetReviewsByListingID gets all reviews for a listing
func (s *ReviewService) GetReviewsByListingID(ctx context.Context, listingID uuid.UUID) ([]*model.Review, error) {
	// Check if listing exists
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return nil, ErrListingNotFound
	}

	// Get reviews
	reviews, err := s.reviewRepo.FindByListingID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

// UpdateReview updates an existing review
func (s *ReviewService) UpdateReview(
	ctx context.Context,
	reviewID uuid.UUID,
	userID uuid.UUID,
	rating int,
	comment string,
) (*model.Review, error) {
	// Get the review
	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}
	if review == nil {
		return nil, ErrReviewNotFound
	}

	// Check if review belongs to the user
	if review.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields
	review.Rating = rating
	review.Comment = comment
	review.UpdatedAt = time.Now()

	// Save changes
	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Update the listing's rating
	avgRating, err := s.reviewRepo.CalculateAverageRating(ctx, review.ListingID)
	if err == nil {
		listing, err := s.listingRepo.FindByID(ctx, review.ListingID)
		if err == nil && listing != nil {
			listing.Ratings = avgRating
			if err := s.listingRepo.Update(ctx, listing); err != nil {
				// Log error but continue, as the review was updated successfully
				fmt.Printf("failed to update listing rating: %v\n", err)
			}
		}
	}

	return review, nil
}

// DeleteReview deletes a review
func (s *ReviewService) DeleteReview(ctx context.Context, reviewID uuid.UUID, userID uuid.UUID) error {
	// Get the review
	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("failed to find review: %w", err)
	}
	if review == nil {
		return ErrReviewNotFound
	}

	// Check if review belongs to the user
	if review.UserID != userID {
		return ErrUnauthorized
	}

	// Delete the review
	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	// Update the listing's rating
	listingID := review.ListingID
	avgRating, err := s.reviewRepo.CalculateAverageRating(ctx, listingID)
	if err == nil {
		listing, err := s.listingRepo.FindByID(ctx, listingID)
		if err == nil && listing != nil {
			listing.Ratings = avgRating
			if err := s.listingRepo.Update(ctx, listing); err != nil {
				// Log error but continue, as the review was deleted successfully
				fmt.Printf("failed to update listing rating: %v\n", err)
			}
		}
	}

	return nil
}
