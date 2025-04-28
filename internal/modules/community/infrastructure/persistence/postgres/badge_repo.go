// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/badge_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// BadgeRepository implements the repository.BadgeRepository interface with PostgreSQL
type BadgeRepository struct {
	db *sql.DB
}

// GetAll retrieves all badges
func (r *BadgeRepository) GetAll(ctx context.Context) ([]*model.Badge, error) {
	query := `
		SELECT id, name, slug, description, icon, color, points, badge_type, criteria, created_at
		FROM badges
		ORDER BY points ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query badges: %w", err)
	}
	defer rows.Close()

	var badges []*model.Badge
	for rows.Next() {
		var badge model.Badge
		var icon, color, criteria sql.NullString

		err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Slug,
			&badge.Description,
			&icon,
			&color,
			&badge.Points,
			&badge.BadgeType,
			&criteria,
			&badge.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan badge row: %w", err)
		}

		// Handle nullable fields
		if icon.Valid {
			badge.Icon = &icon.String
		}

		if color.Valid {
			badge.Color = &color.String
		}

		if criteria.Valid {
			badge.Criteria = &criteria.String
		}

		badges = append(badges, &badge)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during badge row iteration: %w", err)
	}

	return badges, nil
}

// GetByID retrieves a badge by its ID
func (r *BadgeRepository) GetByID(ctx context.Context, id int) (*model.Badge, error) {
	query := `
		SELECT id, name, slug, description, icon, color, points, badge_type, criteria, created_at
		FROM badges
		WHERE id = $1
	`

	var badge model.Badge
	var icon, color, criteria sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&badge.ID,
		&badge.Name,
		&badge.Slug,
		&badge.Description,
		&icon,
		&color,
		&badge.Points,
		&badge.BadgeType,
		&criteria,
		&badge.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get badge: %w", err)
	}

	// Handle nullable fields
	if icon.Valid {
		badge.Icon = &icon.String
	}

	if color.Valid {
		badge.Color = &color.String
	}

	if criteria.Valid {
		badge.Criteria = &criteria.String
	}

	return &badge, nil
}

// Create creates a new badge
func (r *BadgeRepository) Create(ctx context.Context, badge *model.Badge) error {
	query := `
		INSERT INTO badges (name, slug, description, icon, color, points, badge_type, criteria, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		badge.Name,
		badge.Slug,
		badge.Description,
		badge.Icon,
		badge.Color,
		badge.Points,
		badge.BadgeType,
		badge.Criteria,
		badge.CreatedAt,
	).Scan(&badge.ID)

	if err != nil {
		return fmt.Errorf("failed to create badge: %w", err)
	}

	return nil
}

// Update updates a badge
func (r *BadgeRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	// Build dynamic query
	setClause := ""
	args := []interface{}{id} // First arg is the ID
	argCount := 2

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += key + " = $" + fmt.Sprintf("%d", argCount)
		args = append(args, value)
		argCount++
	}

	if setClause == "" {
		return nil // Nothing to update
	}

	query := "UPDATE badges SET " + setClause + " WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update badge: %w", err)
	}

	return nil
}

// Delete deletes a badge
func (r *BadgeRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM badges WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete badge: %w", err)
	}

	return nil
}

// AwardBadge awards a badge to a user
func (r *BadgeRepository) AwardBadge(ctx context.Context, userID uuid.UUID, badgeID int) error {
	query := `
		INSERT INTO user_badges (user_id, badge_id, awarded_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		badgeID,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to award badge: %w", err)
	}

	return nil
}

// HasBadge checks if a user has a badge
func (r *BadgeRepository) HasBadge(ctx context.Context, userID uuid.UUID, badgeID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM user_badges WHERE user_id = $1 AND badge_id = $2)",
		userID,
		badgeID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check badge award: %w", err)
	}

	return exists, nil
}

// GetUserBadges gets badges for a user
func (r *BadgeRepository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]*model.Badge, error) {
	query := `
		SELECT b.id, b.name, b.slug, b.description, b.icon, b.color, b.points, b.badge_type, b.criteria,
			b.created_at, ub.awarded_at
		FROM badges b
		JOIN user_badges ub ON b.id = ub.badge_id
		WHERE ub.user_id = $1
		ORDER BY ub.awarded_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user badges: %w", err)
	}
	defer rows.Close()

	var badges []*model.Badge
	for rows.Next() {
		var badge model.Badge
		var icon, color, criteria sql.NullString
		var awardedAt time.Time

		err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Slug,
			&badge.Description,
			&icon,
			&color,
			&badge.Points,
			&badge.BadgeType,
			&criteria,
			&badge.CreatedAt,
			&awardedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan badge row: %w", err)
		}

		// Handle nullable fields
		if icon.Valid {
			badge.Icon = &icon.String
		}

		if color.Valid {
			badge.Color = &color.String
		}

		if criteria.Valid {
			badge.Criteria = &criteria.String
		}

		badge.AwardedAt = &awardedAt

		badges = append(badges, &badge)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during badge row iteration: %w", err)
	}

	return badges, nil
}

// RevokeBadge revokes a badge from a user
func (r *BadgeRepository) RevokeBadge(ctx context.Context, userID uuid.UUID, badgeID int) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM user_badges WHERE user_id = $1 AND badge_id = $2",
		userID,
		badgeID,
	)

	if err != nil {
		return fmt.Errorf("failed to revoke badge: %w", err)
	}

	return nil
}
