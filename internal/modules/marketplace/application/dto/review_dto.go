// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/review_dto.go
package dto

import (
	"time"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// CreateReviewRequest represents a request to create a review
type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty"`
}

// UpdateReviewRequest represents a request to update a review
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty"`
}

// ReviewResponse represents a review response
type ReviewResponse struct {
	ID        uuid.UUID `json:"id"`
	ListingID uuid.UUID `json:"listing_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToReviewResponse converts a model.Review to a ReviewResponse
func ToReviewResponse(review *model.Review) ReviewResponse {
	return ReviewResponse{
		ID:        review.ID,
		ListingID: review.ListingID,
		UserID:    review.UserID,
		UserName:  review.UserName,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}
