// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/user_reputation_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// UserReputationRepository implements the repository.UserReputationRepository interface with PostgreSQL
type UserReputationRepository struct {
	db *sql.DB
}

// GetByUserID retrieves user reputation by user ID
func (r *UserReputationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error) {
	query := `
		SELECT up.user_id, up.reputation_score, up.contribution_count, up.joined_date, up.last_active,
			u.name
		FROM user_profiles up
		JOIN "user" u ON up.user_id = u.id
		WHERE up.user_id = $1
	`

	var userProfile model.UserProfile
	var name string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&userProfile.ID,
		&userProfile.ReputationScore,
		&userProfile.ContributionCount,
		&userProfile.JoinedDate,
		&userProfile.LastActive,
		&name,
	)

	if err == sql.ErrNoRows {
		// Create a new user profile
		now := time.Now()
		userProfile = model.UserProfile{
			ID:                userID,
			ReputationScore:   0,
			ContributionCount: 0,
			JoinedDate:        now,
			LastActive:        now,
			TopicCount:        0,
			PostCount:         0,
		}

		// Get user name
		err = r.db.QueryRowContext(ctx, "SELECT name FROM \"user\" WHERE id = $1", userID).Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Create user profile
		if err := r.CreateOrUpdate(ctx, &userProfile); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	userProfile.Name = name

	return &userProfile, nil
}

// CreateOrUpdate creates or updates a user reputation profile
func (r *UserReputationRepository) CreateOrUpdate(ctx context.Context, userProfile *model.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, reputation_score, contribution_count, joined_date, last_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			reputation_score = $2,
			contribution_count = $3,
			last_active = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		userProfile.ID,
		userProfile.ReputationScore,
		userProfile.ContributionCount,
		userProfile.JoinedDate,
		userProfile.LastActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create or update user profile: %w", err)
	}

	return nil
}

// IncrementReputationScore increments the reputation score for a user
func (r *UserReputationRepository) IncrementReputationScore(ctx context.Context, userID uuid.UUID, amount int) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user_profiles (user_id, reputation_score, contribution_count, joined_date, last_active)
		VALUES ($1, $2, 0, $3, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			reputation_score = user_profiles.reputation_score + $2,
			last_active = $3`,
		userID,
		amount,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to increment reputation score: %w", err)
	}

	return nil
}

// IncrementContributionCount increments the contribution count for a user
func (r *UserReputationRepository) IncrementContributionCount(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user_profiles (user_id, reputation_score, contribution_count, joined_date, last_active)
		VALUES ($1, 0, 1, $2, $2)
		ON CONFLICT (user_id) DO UPDATE SET
			contribution_count = user_profiles.contribution_count + 1,
			last_active = $2`,
		userID,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to increment contribution count: %w", err)
	}

	return nil
}

// UpdateLastActive updates the last active timestamp for a user
func (r *UserReputationRepository) UpdateLastActive(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user_profiles (user_id, reputation_score, contribution_count, joined_date, last_active)
		VALUES ($1, 0, 0, $2, $2)
		ON CONFLICT (user_id) DO UPDATE SET
			last_active = $2`,
		userID,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to update last active: %w", err)
	}

	return nil
}

// GetTopicCount gets the number of topics created by a user
func (r *UserReputationRepository) GetTopicCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM topics WHERE author_id = $1 AND is_hidden = false",
		userID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get topic count: %w", err)
	}

	return count, nil
}

// GetPostCount gets the number of posts created by a user
func (r *UserReputationRepository) GetPostCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM posts WHERE author_id = $1",
		userID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get post count: %w", err)
	}

	return count, nil
}
