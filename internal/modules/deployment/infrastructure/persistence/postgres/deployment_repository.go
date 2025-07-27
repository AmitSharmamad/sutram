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
)

// DeploymentRepository implements repository.DeploymentRepository
type DeploymentRepository struct {
	db *sql.DB
}

// NewDeploymentRepository creates a new deployment repository
func NewDeploymentRepository(db *sql.DB) repository.DeploymentRepository {
	return &DeploymentRepository{
		db: db,
	}
}

// Create creates a new deployment
func (r *DeploymentRepository) Create(ctx context.Context, deployment *model.Deployment) error {
	query := `
		INSERT INTO deployment (
			id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		deployment.ID,
		deployment.UserID,
		deployment.StrategyID,
		deployment.InstrumentID,
		deployment.Name,
		deployment.Description,
		deployment.IsActive,
		deployment.CreatedAt,
		deployment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

// GetByID retrieves a deployment by its ID
func (r *DeploymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE id = $1
	`

	var deployment model.Deployment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&deployment.ID,
		&deployment.UserID,
		&deployment.StrategyID,
		&deployment.InstrumentID,
		&deployment.Name,
		&deployment.Description,
		&deployment.IsActive,
		&deployment.CreatedAt,
		&deployment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return &deployment, nil
}

// GetByUserID retrieves all deployments for a specific user
func (r *DeploymentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}

// GetByStrategyID retrieves all deployments for a specific strategy
func (r *DeploymentRepository) GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE strategy_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}

// GetByInstrumentID retrieves all deployments for a specific instrument
func (r *DeploymentRepository) GetByInstrumentID(ctx context.Context, instrumentID int) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE instrument_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, instrumentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}

// GetByUserAndStrategy retrieves all deployments for a specific user and strategy
func (r *DeploymentRepository) GetByUserAndStrategy(ctx context.Context, userID, strategyID uuid.UUID) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE user_id = $1 AND strategy_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}

// GetByUserStrategyInstrument retrieves a deployment by user, strategy, and instrument
func (r *DeploymentRepository) GetByUserStrategyInstrument(ctx context.Context, userID, strategyID uuid.UUID, instrumentID int) (*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE user_id = $1 AND strategy_id = $2 AND instrument_id = $3
	`

	var deployment model.Deployment
	err := r.db.QueryRowContext(ctx, query, userID, strategyID, instrumentID).Scan(
		&deployment.ID,
		&deployment.UserID,
		&deployment.StrategyID,
		&deployment.InstrumentID,
		&deployment.Name,
		&deployment.Description,
		&deployment.IsActive,
		&deployment.CreatedAt,
		&deployment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return &deployment, nil
}

// GetActiveByUserStrategyInstrument retrieves an active deployment by user, strategy, and instrument
func (r *DeploymentRepository) GetActiveByUserStrategyInstrument(ctx context.Context, userID, strategyID uuid.UUID, instrumentID int) (*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE user_id = $1 AND strategy_id = $2 AND instrument_id = $3 AND is_active = true
	`

	var deployment model.Deployment
	err := r.db.QueryRowContext(ctx, query, userID, strategyID, instrumentID).Scan(
		&deployment.ID,
		&deployment.UserID,
		&deployment.StrategyID,
		&deployment.InstrumentID,
		&deployment.Name,
		&deployment.Description,
		&deployment.IsActive,
		&deployment.CreatedAt,
		&deployment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil without error if no active deployment found
		}
		return nil, fmt.Errorf("failed to get active deployment: %w", err)
	}

	return &deployment, nil
}

// Update updates an existing deployment
func (r *DeploymentRepository) Update(ctx context.Context, deployment *model.Deployment) error {
	query := `
		UPDATE deployment
		SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		deployment.Name,
		deployment.Description,
		deployment.IsActive,
		time.Now(),
		deployment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

// Delete removes a deployment
func (r *DeploymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM deployment WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	return nil
}

// ActivateDeployment sets the deployment as active
func (r *DeploymentRepository) ActivateDeployment(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE deployment
		SET is_active = true, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to activate deployment: %w", err)
	}

	return nil
}

// DeactivateDeployment sets the deployment as inactive
func (r *DeploymentRepository) DeactivateDeployment(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE deployment
		SET is_active = false, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to deactivate deployment: %w", err)
	}

	return nil
}

// GetActiveDeployments retrieves all active deployments
func (r *DeploymentRepository) GetActiveDeployments(ctx context.Context) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}

// GetActiveDeploymentsByUserID retrieves all active deployments for a specific user
func (r *DeploymentRepository) GetActiveDeploymentsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}

// GetInactiveDeploymentsByUserID retrieves all inactive deployments for a specific user
func (r *DeploymentRepository) GetInactiveDeploymentsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Deployment, error) {
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE user_id = $1 AND is_active = false
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inactive deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployment rows: %w", err)
	}

	return deployments, nil
}

// SearchDeployments searches deployments with filter criteria
func (r *DeploymentRepository) SearchDeployments(ctx context.Context, filters map[string]interface{}) ([]*model.Deployment, error) {
	// Start with base query
	query := `
		SELECT id, user_id, strategy_id, instrument_id, name, description, is_active, created_at, updated_at
		FROM deployment
		WHERE 1=1
	`

	// Build dynamic WHERE clause and params based on filters
	var params []any
	paramIndex := 1

	if userID, ok := filters["user_id"]; ok {
		query += fmt.Sprintf(" AND user_id = $%d", paramIndex)
		params = append(params, userID)
		paramIndex++
	}

	if strategyID, ok := filters["strategy_id"]; ok {
		query += fmt.Sprintf(" AND strategy_id = $%d", paramIndex)
		params = append(params, strategyID)
		paramIndex++
	}

	if instrumentID, ok := filters["instrument_id"]; ok {
		query += fmt.Sprintf(" AND instrument_id = $%d", paramIndex)
		params = append(params, instrumentID)
		paramIndex++
	}

	if isActive, ok := filters["is_active"]; ok {
		query += fmt.Sprintf(" AND is_active = $%d", paramIndex)
		params = append(params, isActive)
		paramIndex++
	}

	if name, ok := filters["name"]; ok {
		query += fmt.Sprintf(" AND name ILIKE $%d", paramIndex)
		params = append(params, "%"+name.(string)+"%")
		paramIndex++
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to search deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*model.Deployment
	for rows.Next() {
		var deployment model.Deployment
		err := rows.Scan(
			&deployment.ID,
			&deployment.UserID,
			&deployment.StrategyID,
			&deployment.InstrumentID,
			&deployment.Name,
			&deployment.Description,
			&deployment.IsActive,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, &deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployments: %w", err)
	}

	return deployments, nil
}
