package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"

	"github.com/google/uuid"
)

// RunGroupRepository implements the repository.RunGroupRepository interface

// RunGroupRepository implements the repository.RunGroupRepository interface
type RunGroupRepository struct {
	db *sql.DB
}

// NewRunGroupRepository creates a new PostgreSQL repository for run groups
func NewRunGroupRepository(db *sql.DB) repository.RunGroupRepository {
	return &RunGroupRepository{
		db: db,
	}
}

// Create saves a new run group in the database
func (r *RunGroupRepository) Create(ctx context.Context, runGroup *model.RunGroup) error {
	query := `
		INSERT INTO run_group (
			id, strategy_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		runGroup.ID,
		runGroup.StrategyID,
		runGroup.CreatedAt,
		runGroup.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID retrieves a run group by its ID
func (r *RunGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.RunGroup, error) {
	query := `
		SELECT 
			id, strategy_id, created_at, updated_at
		FROM run_group
		WHERE id = $1
	`

	var runGroup model.RunGroup
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&runGroup.ID,
		&runGroup.StrategyID,
		&runGroup.CreatedAt,
		&runGroup.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRunGroupNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &runGroup, nil
}

// GetByStrategyID retrieves all run groups for a specific strategy
func (r *RunGroupRepository) GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.RunGroup, error) {
	query := `
		SELECT 
			id, strategy_id, created_at, updated_at
		FROM run_group
		WHERE strategy_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, strategyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var runGroups []*model.RunGroup
	for rows.Next() {
		var runGroup model.RunGroup
		if err := rows.Scan(
			&runGroup.ID,
			&runGroup.StrategyID,
			&runGroup.CreatedAt,
			&runGroup.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		runGroups = append(runGroups, &runGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return runGroups, nil
}

// Delete removes a run group from the database
func (r *RunGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM run_group
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrRunGroupNotFound
	}

	return nil
}
