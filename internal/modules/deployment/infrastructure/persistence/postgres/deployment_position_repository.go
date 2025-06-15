package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
			exit_price, exit_time, pnl, status, direction, stop_loss_price, 
			initial_stop_loss_price, target_profit_price, metadata, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		RETURNING id
	`

	var id uuid.UUID
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
		position.Direction,
		position.StopLossPrice,
		position.InitialStopLossPrice,
		position.TargetProfitPrice,
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
			exit_price, exit_time, pnl, status, direction, stop_loss_price,
			initial_stop_loss_price, target_profit_price, metadata, created_at, updated_at
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
			exit_price, exit_time, pnl, status, direction, stop_loss_price,
			initial_stop_loss_price, target_profit_price, metadata, created_at, updated_at
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
			exit_price, exit_time, pnl, status, direction, stop_loss_price,
			initial_stop_loss_price, target_profit_price, metadata, created_at, updated_at
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

// GetOpenPositionsByDeploymentIDs retrieves all open positions for multiple deployments
// Returns a map of deployment ID to the first open position (assuming one open position per deployment)
func (r *PgDeploymentPositionRepository) GetOpenPositionsByDeploymentIDs(ctx context.Context, deploymentIDs []uuid.UUID) (map[uuid.UUID]*model.DeploymentPosition, error) {
	if len(deploymentIDs) == 0 {
		return make(map[uuid.UUID]*model.DeploymentPosition), nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]interface{}, len(deploymentIDs)+1)
	placeholderStrings := make([]string, len(deploymentIDs))

	for i, id := range deploymentIDs {
		placeholders[i] = id
		placeholderStrings[i] = fmt.Sprintf("$%d", i+1)
	}
	placeholders[len(deploymentIDs)] = model.PositionOpen

	query := fmt.Sprintf(`
		SELECT 
			id, deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, direction, stop_loss_price,
			initial_stop_loss_price, target_profit_price, metadata, created_at, updated_at
		FROM deployment_position
		WHERE deployment_id IN (%s) AND status = $%d
		ORDER BY deployment_id, entry_time DESC
	`,
		strings.Join(placeholderStrings, ","),
		len(deploymentIDs)+1)

	var positions []*model.DeploymentPosition
	err := r.db.SelectContext(ctx, &positions, query, placeholders...)
	if err != nil {
		return nil, fmt.Errorf("failed to get open deployment positions by deployment IDs: %w", err)
	}

	// Create map with first position per deployment (assuming one open position per deployment)
	positionMap := make(map[uuid.UUID]*model.DeploymentPosition)
	for _, position := range positions {
		if _, exists := positionMap[position.DeploymentID]; !exists {
			positionMap[position.DeploymentID] = position
		}
	}

	return positionMap, nil
}

// GetByDeploymentAndInstrument retrieves all positions for a specific deployment and instrument
func (r *PgDeploymentPositionRepository) GetByDeploymentAndInstrument(ctx context.Context, deploymentID uuid.UUID, instrumentID int) ([]*model.DeploymentPosition, error) {
	query := `
		SELECT 
			id, deployment_id, instrument_id, quantity, entry_price, entry_time,
			exit_price, exit_time, pnl, status, direction, stop_loss_price,
			initial_stop_loss_price, target_profit_price, metadata, created_at, updated_at
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
			direction = $8,
			stop_loss_price = $9,
			initial_stop_loss_price = $10,
			target_profit_price = $11,
			metadata = $12,
			updated_at = $13
		WHERE id = $14
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
		position.Direction,
		position.StopLossPrice,
		position.InitialStopLossPrice,
		position.TargetProfitPrice,
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
