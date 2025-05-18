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

// RunRepository implements the repository.RunRepository interface
type RunRepository struct {
	db *sql.DB
}

// NewRunRepository creates a new PostgreSQL repository for runs
func NewRunRepository(db *sql.DB) repository.RunRepository {
	return &RunRepository{
		db: db,
	}
}

// Create saves a new run in the database
func (r *RunRepository) Create(ctx context.Context, run *model.Run) error {
	query := `
		INSERT INTO run (
			id, run_group_id, strategy_id, exchange, instrument_id, 
			type, state, qty, period_start_time, period_end_time, 
			started_time, ended_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		run.ID,
		run.RunGroupID,
		run.StrategyID,
		run.Exchange,
		run.InstrumentID,
		run.Type,
		run.State,
		run.Qty,
		run.PeriodStartTime,
		run.PeriodEndTime,
		run.StartedTime,
		run.EndedTime,
		run.CreatedAt,
		run.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID retrieves a run by its ID
func (r *RunRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Run, error) {
	query := `
		SELECT 
			id, run_group_id, strategy_id, exchange, instrument_id, 
			type, state, qty, period_start_time, period_end_time, 
			started_time, ended_time, created_at, updated_at
		FROM run
		WHERE id = $1
	`

	var run model.Run
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&run.ID,
		&run.RunGroupID,
		&run.StrategyID,
		&run.Exchange,
		&run.InstrumentID,
		&run.Type,
		&run.State,
		&run.Qty,
		&run.PeriodStartTime,
		&run.PeriodEndTime,
		&run.StartedTime,
		&run.EndedTime,
		&run.CreatedAt,
		&run.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRunNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &run, nil
}

// GetByRunGroupID retrieves all runs for a specific run group
func (r *RunRepository) GetByRunGroupID(ctx context.Context, runGroupID uuid.UUID) ([]*model.Run, error) {
	query := `
		SELECT 
			id, run_group_id, strategy_id, exchange, instrument_id, 
			type, state, qty, period_start_time, period_end_time, 
			started_time, ended_time, created_at, updated_at
		FROM run
		WHERE run_group_id = $1
		ORDER BY created_at DESC
	`

	return r.queryRuns(ctx, query, runGroupID)
}

// GetByStrategyID retrieves all runs for a specific strategy
func (r *RunRepository) GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.Run, error) {
	query := `
		SELECT 
			id, run_group_id, strategy_id, exchange, instrument_id, 
			type, state, qty, period_start_time, period_end_time, 
			started_time, ended_time, created_at, updated_at
		FROM run
		WHERE strategy_id = $1
		ORDER BY created_at DESC
	`

	return r.queryRuns(ctx, query, strategyID)
}

// Update updates an existing run in the database
func (r *RunRepository) Update(ctx context.Context, run *model.Run) error {
	query := `
		UPDATE run
		SET 
			run_group_id = $2, 
			strategy_id = $3, 
			exchange = $4, 
			instrument_id = $5,
			type = $6, 
			state = $7, 
			qty = $8, 
			period_start_time = $9, 
			period_end_time = $10,
			started_time = $11, 
			ended_time = $12, 
			updated_at = $13
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		run.ID,
		run.RunGroupID,
		run.StrategyID,
		run.Exchange,
		run.InstrumentID,
		run.Type,
		run.State,
		run.Qty,
		run.PeriodStartTime,
		run.PeriodEndTime,
		run.StartedTime,
		run.EndedTime,
		run.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrRunNotFound
	}

	return nil
}

// Delete removes a run from the database
func (r *RunRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM run
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
		return ErrRunNotFound
	}

	return nil
}

// UpdateState updates the state of a run
func (r *RunRepository) UpdateState(ctx context.Context, id uuid.UUID, state model.RunState) error {
	query := `
		UPDATE run
		SET 
			state = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, state)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrRunNotFound
	}

	return nil
}

// SearchRuns searches runs with filter criteria
func (r *RunRepository) SearchRuns(ctx context.Context, filters map[string]interface{}) ([]*model.Run, error) {
	// Base query
	query := `
		SELECT 
			id, run_group_id, strategy_id, exchange, instrument_id, 
			type, state, qty, period_start_time, period_end_time, 
			started_time, ended_time, created_at, updated_at
		FROM run
		WHERE 1=1
	`

	// Build where clause and args
	var args []interface{}
	argCount := 1

	if strategyID, ok := filters["strategy_id"]; ok {
		query += fmt.Sprintf(" AND strategy_id = $%d", argCount)
		args = append(args, strategyID)
		argCount++
	}

	if runType, ok := filters["type"]; ok {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, runType)
		argCount++
	}

	if state, ok := filters["state"]; ok {
		query += fmt.Sprintf(" AND state = $%d", argCount)
		args = append(args, state)
		argCount++
	}

	if exchange, ok := filters["exchange"]; ok {
		query += fmt.Sprintf(" AND exchange = $%d", argCount)
		args = append(args, exchange)
		argCount++
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var runs []*model.Run
	for rows.Next() {
		var run model.Run
		if err := rows.Scan(
			&run.ID,
			&run.RunGroupID,
			&run.StrategyID,
			&run.Exchange,
			&run.InstrumentID,
			&run.Type,
			&run.State,
			&run.Qty,
			&run.PeriodStartTime,
			&run.PeriodEndTime,
			&run.StartedTime,
			&run.EndedTime,
			&run.CreatedAt,
			&run.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		runs = append(runs, &run)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return runs, nil
}

// queryRuns executes a query and returns a slice of runs
func (r *RunRepository) queryRuns(ctx context.Context, query string, args ...interface{}) ([]*model.Run, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var runs []*model.Run
	for rows.Next() {
		var run model.Run
		if err := rows.Scan(
			&run.ID,
			&run.RunGroupID,
			&run.StrategyID,
			&run.Exchange,
			&run.InstrumentID,
			&run.Type,
			&run.State,
			&run.Qty,
			&run.PeriodStartTime,
			&run.PeriodEndTime,
			&run.StartedTime,
			&run.EndedTime,
			&run.CreatedAt,
			&run.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		runs = append(runs, &run)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return runs, nil
}

// BulkCreate saves multiple runs in the database in a single operation
func (r *RunRepository) BulkCreate(ctx context.Context, runs []*model.Run) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback()

	// Prepare the statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO run (
			id, run_group_id, strategy_id, exchange, instrument_id, 
			type, state, qty, period_start_time, period_end_time, 
			started_time, ended_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer stmt.Close()

	// Execute for each run
	for _, run := range runs {
		_, err := stmt.ExecContext(
			ctx,
			run.ID,
			run.RunGroupID,
			run.StrategyID,
			run.Exchange,
			run.InstrumentID,
			run.Type,
			run.State,
			run.Qty,
			run.PeriodStartTime,
			run.PeriodEndTime,
			run.StartedTime,
			run.EndedTime,
			run.CreatedAt,
			run.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}
