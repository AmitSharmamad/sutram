package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PgDeploymentPositionRepository implements the deployment position repository interface
type PgDeploymentPositionRepository struct {
	db *sqlx.DB
}

// NewPgDeploymentPositionRepository creates a new Postgres deployment position repository
func NewPgDeploymentPositionRepository(db *sqlx.DB) repository.DeploymentPositionRepository {
	return &PgDeploymentPositionRepository{
		db: db,
	}
}

// Create saves a new position in the database
func (r *PgDeploymentPositionRepository) Create(ctx context.Context, position *model.DeploymentPosition) error {
	query := `
		INSERT INTO deployment_position (
			deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, metadata, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		position.DeploymentID,
		position.InstrumentID,
		position.Quantity,
		position.EntryPrice,
		position.EntryTime,
		position.ExitPrice,
		position.ExitTime,
		position.PNL,
		position.Status,
		position.Metadata,
		position.CreatedAt,
		position.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create deployment position: %w", err)
	}

	position.ID = id
	return nil
}

// GetByID retrieves a position by its ID
func (r *PgDeploymentPositionRepository) GetByID(ctx context.Context, id int64) (*model.DeploymentPosition, error) {
	query := `
		SELECT 
			id, deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, metadata, created_at, updated_at
		FROM deployment_position
		WHERE id = $1
	`

	var position model.DeploymentPosition
	err := r.db.GetContext(ctx, &position, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get deployment position by ID: %w", err)
	}

	return &position, nil
}

// GetByDeploymentID retrieves all positions for a specific deployment
func (r *PgDeploymentPositionRepository) GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentPosition, error) {
	query := `
		SELECT 
			id, deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, metadata, created_at, updated_at
		FROM deployment_position
		WHERE deployment_id = $1
		ORDER BY entry_time DESC
	`

	var positions []*model.DeploymentPosition
	err := r.db.SelectContext(ctx, &positions, query, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment positions by deployment ID: %w", err)
	}

	return positions, nil
}

// GetOpenPositionsByDeploymentID retrieves all open positions for a specific deployment
func (r *PgDeploymentPositionRepository) GetOpenPositionsByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentPosition, error) {
	query := `
		SELECT 
			id, deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, metadata, created_at, updated_at
		FROM deployment_position
		WHERE deployment_id = $1 AND status = $2
		ORDER BY entry_time DESC
	`

	var positions []*model.DeploymentPosition
	err := r.db.SelectContext(ctx, &positions, query, deploymentID, model.PositionOpen)
	if err != nil {
		return nil, fmt.Errorf("failed to get open deployment positions by deployment ID: %w", err)
	}

	return positions, nil
}

// GetByDeploymentAndInstrument retrieves all positions for a specific deployment and instrument
func (r *PgDeploymentPositionRepository) GetByDeploymentAndInstrument(ctx context.Context, deploymentID uuid.UUID, instrumentID int) ([]*model.DeploymentPosition, error) {
	query := `
		SELECT 
			id, deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, metadata, created_at, updated_at
		FROM deployment_position
		WHERE deployment_id = $1 AND instrument_id = $2
		ORDER BY entry_time DESC
	`

	var positions []*model.DeploymentPosition
	err := r.db.SelectContext(ctx, &positions, query, deploymentID, instrumentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment positions by deployment ID and instrument ID: %w", err)
	}

	return positions, nil
}

// Update updates an existing position in the database
func (r *PgDeploymentPositionRepository) Update(ctx context.Context, position *model.DeploymentPosition) error {
	query := `
		UPDATE deployment_position
		SET
			quantity = $1,
			entry_price = $2,
			entry_time = $3,
			exit_price = $4,
			exit_time = $5,
			pnl = $6,
			status = $7,
			metadata = $8,
			updated_at = $9
		WHERE id = $10
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		position.Quantity,
		position.EntryPrice,
		position.EntryTime,
		position.ExitPrice,
		position.ExitTime,
		position.PNL,
		position.Status,
		position.Metadata,
		position.UpdatedAt,
		position.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update deployment position: %w", err)
	}

	return nil
}

// ClosePosition updates a position with exit details
func (r *PgDeploymentPositionRepository) ClosePosition(ctx context.Context, id int64, exitPrice float64, exitTime time.Time) error {
	// First retrieve the position to calculate PnL
	position, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if position == nil {
		return errors.New("position not found")
	}

	position.ClosePosition(exitPrice, exitTime)

	return r.Update(ctx, position)
}

// CancelPosition marks a position as canceled
func (r *PgDeploymentPositionRepository) CancelPosition(ctx context.Context, id int64) error {
	position, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if position == nil {
		return errors.New("position not found")
	}

	position.CancelPosition()

	return r.Update(ctx, position)
}

// CountPositions counts positions by deployment ID and optional status
func (r *PgDeploymentPositionRepository) CountPositions(ctx context.Context, deploymentID uuid.UUID, status *model.PositionStatus) (int, error) {
	var query string
	var args []interface{}

	if status != nil {
		query = `
			SELECT COUNT(*)
			FROM deployment_position
			WHERE deployment_id = $1 AND status = $2
		`
		args = []interface{}{deploymentID, *status}
	} else {
		query = `
			SELECT COUNT(*)
			FROM deployment_position
			WHERE deployment_id = $1
		`
		args = []interface{}{deploymentID}
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count deployment positions: %w", err)
	}

	return count, nil
}
