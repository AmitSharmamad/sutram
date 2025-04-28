// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/like_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// LikeRepository implements the repository.LikeRepository interface with PostgreSQL
type LikeRepository struct {
	db *sql.DB
}

// Create creates a new like
func (r *LikeRepository) Create(ctx context.Context, like *model.Like) error {
	query := `
		INSERT INTO likes (user_id, likeable_type, likeable_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		like.UserID,
		like.LikeableType,
		like.LikeableID,
		like.CreatedAt,
	).Scan(&like.ID)

	if err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	return nil
}

// Delete deletes a like
func (r *LikeRepository) Delete(ctx context.Context, userID uuid.UUID, likeableType string, likeableID int) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM likes WHERE user_id = $1 AND likeable_type = $2 AND likeable_id = $3",
		userID,
		likeableType,
		likeableID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	return nil
}

// Exists checks if a like exists
func (r *LikeRepository) Exists(ctx context.Context, userID uuid.UUID, likeableType string, likeableID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND likeable_type = $2 AND likeable_id = $3)",
		userID,
		likeableType,
		likeableID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check like existence: %w", err)
	}

	return exists, nil
}

// GetLikeCount gets the count of likes for a likeable object
func (r *LikeRepository) GetLikeCount(ctx context.Context, likeableType string, likeableID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM likes WHERE likeable_type = $1 AND likeable_id = $2",
		likeableType,
		likeableID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	return count, nil
}

// GetByUserID gets likes created by a user
func (r *LikeRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page int, perPage int) ([]*model.Like, int, error) {
	// Count total matching rows
	countQuery := "SELECT COUNT(*) FROM likes WHERE user_id = $1"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user likes: %w", err)
	}

	// Calculate limit and offset
	offset := (page - 1) * perPage

	query := `
		SELECT id, user_id, likeable_type, likeable_id, created_at
		FROM likes
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query user likes: %w", err)
	}
	defer rows.Close()

	var likes []*model.Like
	for rows.Next() {
		var like model.Like

		err := rows.Scan(
			&like.ID,
			&like.UserID,
			&like.LikeableType,
			&like.LikeableID,
			&like.CreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan like row: %w", err)
		}

		likes = append(likes, &like)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during like row iteration: %w", err)
	}

	return likes, total, nil
}
