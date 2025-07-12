package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"algo/internal/domain/entity"
	"algo/internal/domain/repository"

	"github.com/google/uuid"
)

// UsageLimitRepository implements repository.UsageLimitRepository for PostgreSQL
type UsageLimitRepository struct {
	db *sql.DB
}

// NewUsageLimitRepository creates a new PostgreSQL usage limit repository
func NewUsageLimitRepository(db *sql.DB) repository.UsageLimitRepository {
	return &UsageLimitRepository{db: db}
}

// Create creates a new usage limit
func (r *UsageLimitRepository) Create(ctx context.Context, limit *entity.UsageLimit) error {
	query := `
		INSERT INTO user_usage_limits (user_id, limit_type, limit_value, current_usage, reset_period, last_reset_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	now := time.Now()
	limit.CreatedAt = now
	limit.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		limit.UserID,
		string(limit.LimitType),
		limit.LimitValue,
		limit.CurrentUsage,
		string(limit.ResetPeriod),
		limit.LastResetAt,
		limit.CreatedAt,
		limit.UpdatedAt,
	).Scan(&limit.ID)

	if err != nil {
		return fmt.Errorf("failed to create usage limit: %w", err)
	}

	return nil
}

// GetByUserIDAndType gets a usage limit by user ID and limit type
func (r *UsageLimitRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) (*entity.UsageLimit, error) {
	query := `
		SELECT id, user_id, limit_type, limit_value, current_usage, reset_period, last_reset_at, created_at, updated_at
		FROM user_usage_limits
		WHERE user_id = $1 AND limit_type = $2`

	limit := &entity.UsageLimit{}
	var limitTypeStr, resetPeriodStr string

	err := r.db.QueryRowContext(ctx, query, userID, string(limitType)).Scan(
		&limit.ID,
		&limit.UserID,
		&limitTypeStr,
		&limit.LimitValue,
		&limit.CurrentUsage,
		&resetPeriodStr,
		&limit.LastResetAt,
		&limit.CreatedAt,
		&limit.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get usage limit: %w", err)
	}

	limit.LimitType = entity.UsageLimitType(limitTypeStr)
	limit.ResetPeriod = entity.ResetPeriodType(resetPeriodStr)

	return limit, nil
}

// GetByUserID gets all usage limits for a user
func (r *UsageLimitRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UsageLimit, error) {
	query := `
		SELECT id, user_id, limit_type, limit_value, current_usage, reset_period, last_reset_at, created_at, updated_at
		FROM user_usage_limits
		WHERE user_id = $1
		ORDER BY limit_type`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query usage limits: %w", err)
	}
	defer rows.Close()

	var limits []*entity.UsageLimit
	for rows.Next() {
		limit := &entity.UsageLimit{}
		var limitTypeStr, resetPeriodStr string

		err := rows.Scan(
			&limit.ID,
			&limit.UserID,
			&limitTypeStr,
			&limit.LimitValue,
			&limit.CurrentUsage,
			&resetPeriodStr,
			&limit.LastResetAt,
			&limit.CreatedAt,
			&limit.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan usage limit: %w", err)
		}

		limit.LimitType = entity.UsageLimitType(limitTypeStr)
		limit.ResetPeriod = entity.ResetPeriodType(resetPeriodStr)

		limits = append(limits, limit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating usage limit rows: %w", err)
	}

	return limits, nil
}

// Update updates a usage limit
func (r *UsageLimitRepository) Update(ctx context.Context, limit *entity.UsageLimit) error {
	query := `
		UPDATE user_usage_limits
		SET limit_value = $1, current_usage = $2, reset_period = $3, last_reset_at = $4, updated_at = $5
		WHERE id = $6`

	limit.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		limit.LimitValue,
		limit.CurrentUsage,
		string(limit.ResetPeriod),
		limit.LastResetAt,
		limit.UpdatedAt,
		limit.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update usage limit: %w", err)
	}

	return nil
}

// IncrementUsage increments the current usage for a user and limit type
func (r *UsageLimitRepository) IncrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error {
	query := `
		UPDATE user_usage_limits
		SET current_usage = current_usage + 1, updated_at = now()
		WHERE user_id = $1 AND limit_type = $2`

	result, err := r.db.ExecContext(ctx, query, userID, string(limitType))
	if err != nil {
		return fmt.Errorf("failed to increment usage: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no usage limit found for user %s and type %s", userID, limitType)
	}

	return nil
}

// DecrementUsage decrements the current usage for a user and limit type
func (r *UsageLimitRepository) DecrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error {
	query := `
		UPDATE user_usage_limits
		SET current_usage = GREATEST(current_usage - 1, 0), updated_at = now()
		WHERE user_id = $1 AND limit_type = $2`

	result, err := r.db.ExecContext(ctx, query, userID, string(limitType))
	if err != nil {
		return fmt.Errorf("failed to decrement usage: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no usage limit found for user %s and type %s", userID, limitType)
	}

	return nil
}

// ResetUsage resets the current usage for a user and limit type
func (r *UsageLimitRepository) ResetUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error {
	query := `
		UPDATE user_usage_limits
		SET current_usage = 0, last_reset_at = now(), updated_at = now()
		WHERE user_id = $1 AND limit_type = $2`

	result, err := r.db.ExecContext(ctx, query, userID, string(limitType))
	if err != nil {
		return fmt.Errorf("failed to reset usage: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no usage limit found for user %s and type %s", userID, limitType)
	}

	return nil
}

// ResetAllUsageByType resets all usage for a specific limit type
func (r *UsageLimitRepository) ResetAllUsageByType(ctx context.Context, limitType entity.UsageLimitType) error {
	query := `
		UPDATE user_usage_limits
		SET current_usage = 0, last_reset_at = now(), updated_at = now()
		WHERE limit_type = $1`

	_, err := r.db.ExecContext(ctx, query, string(limitType))
	if err != nil {
		return fmt.Errorf("failed to reset all usage for type %s: %w", limitType, err)
	}

	return nil
}

// SetUserLimit sets a custom limit for a user
func (r *UsageLimitRepository) SetUserLimit(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limitValue int) error {
	query := `
		UPDATE user_usage_limits
		SET limit_value = $1, updated_at = now()
		WHERE user_id = $2 AND limit_type = $3`

	result, err := r.db.ExecContext(ctx, query, limitValue, userID, string(limitType))
	if err != nil {
		return fmt.Errorf("failed to set user limit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no usage limit found for user %s and type %s", userID, limitType)
	}

	return nil
}

// GetUsersRequiringReset gets users whose usage should be reset for a specific limit type
func (r *UsageLimitRepository) GetUsersRequiringReset(ctx context.Context, limitType entity.UsageLimitType) ([]*entity.UsageLimit, error) {
	// This is a simplified implementation - in practice, you'd want more sophisticated reset logic
	query := `
		SELECT id, user_id, limit_type, limit_value, current_usage, reset_period, last_reset_at, created_at, updated_at
		FROM user_usage_limits
		WHERE limit_type = $1 AND current_usage > 0`

	rows, err := r.db.QueryContext(ctx, query, string(limitType))
	if err != nil {
		return nil, fmt.Errorf("failed to query users requiring reset: %w", err)
	}
	defer rows.Close()

	var limits []*entity.UsageLimit
	for rows.Next() {
		limit := &entity.UsageLimit{}
		var limitTypeStr, resetPeriodStr string

		err := rows.Scan(
			&limit.ID,
			&limit.UserID,
			&limitTypeStr,
			&limit.LimitValue,
			&limit.CurrentUsage,
			&resetPeriodStr,
			&limit.LastResetAt,
			&limit.CreatedAt,
			&limit.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan usage limit: %w", err)
		}

		limit.LimitType = entity.UsageLimitType(limitTypeStr)
		limit.ResetPeriod = entity.ResetPeriodType(resetPeriodStr)

		limits = append(limits, limit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating usage limit rows: %w", err)
	}

	return limits, nil
}

// BulkResetUsage resets usage for all users of a specific limit type
func (r *UsageLimitRepository) BulkResetUsage(ctx context.Context, limitType entity.UsageLimitType) error {
	return r.ResetAllUsageByType(ctx, limitType)
}
