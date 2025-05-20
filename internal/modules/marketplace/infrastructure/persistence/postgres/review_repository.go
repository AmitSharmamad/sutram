// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/infrastructure/persistence/postgres/review_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// ReviewRepository implements repository.ReviewRepository
type ReviewRepository struct {
	db *sql.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *sql.DB) repository.ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

// Create creates a new review
func (r *ReviewRepository) Create(ctx context.Context, review *model.Review) error {
	query := `
		INSERT INTO marketplace_review (
			id, listing_id, user_id, rating, comment, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		review.ID,
		review.ListingID,
		review.UserID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
		review.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

// FindByID returns a review by ID
func (r *ReviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Review, error) {
	query := `
		SELECT 
			r.id, r.listing_id, r.user_id, r.rating, r.comment, r.created_at, r.updated_at,
			u.name as user_name
		FROM marketplace_review r
		LEFT JOIN "user" u ON r.user_id = u.id
		WHERE r.id = $1
	`

	var review model.Review
	var userName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.ListingID,
		&review.UserID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
		&userName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	if userName.Valid {
		review.UserName = userName.String
	}

	return &review, nil
}

// FindByListingID returns all reviews for a listing
func (r *ReviewRepository) FindByListingID(ctx context.Context, listingID uuid.UUID) ([]*model.Review, error) {
	query := `
		SELECT 
			r.id, r.listing_id, r.user_id, r.rating, r.comment, r.created_at, r.updated_at,
			u.name as user_name
		FROM marketplace_review r
		LEFT JOIN "user" u ON r.user_id = u.id
		WHERE r.listing_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		var review model.Review
		var userName sql.NullString

		err := rows.Scan(
			&review.ID,
			&review.ListingID,
			&review.UserID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
			&userName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review row: %w", err)
		}

		if userName.Valid {
			review.UserName = userName.String
		}

		reviews = append(reviews, &review)
	}

	return reviews, nil
}

// FindByUserAndListingID returns a review by user and listing
func (r *ReviewRepository) FindByUserAndListingID(ctx context.Context, userID uuid.UUID, listingID uuid.UUID) (*model.Review, error) {
	query := `
		SELECT 
			id, listing_id, user_id, rating, comment, created_at, updated_at
		FROM marketplace_review
		WHERE user_id = $1 AND listing_id = $2
	`

	var review model.Review
	err := r.db.QueryRowContext(ctx, query, userID, listingID).Scan(
		&review.ID,
		&review.ListingID,
		&review.UserID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	return &review, nil
}

// Update updates a review
func (r *ReviewRepository) Update(ctx context.Context, review *model.Review) error {
	query := `
		UPDATE marketplace_review SET
			rating = $1,
			comment = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		review.Rating,
		review.Comment,
		time.Now(),
		review.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

// Delete deletes a review
func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM marketplace_review WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

// CalculateAverageRating calculates the average rating for a listing
func (r *ReviewRepository) CalculateAverageRating(ctx context.Context, listingID uuid.UUID) (float32, error) {
	query := `
		SELECT AVG(rating)
		FROM marketplace_review
		WHERE listing_id = $1
	`

	var avgRating sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, listingID).Scan(&avgRating)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate average rating: %w", err)
	}

	if !avgRating.Valid {
		return 0, nil
	}

	return float32(avgRating.Float64), nil
}
