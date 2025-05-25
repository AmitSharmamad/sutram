package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PgDeploymentLogRepository implements the deployment log repository interface
type PgDeploymentLogRepository struct {
	db *sqlx.DB
}

// NewPgDeploymentLogRepository creates a new Postgres deployment log repository
func NewPgDeploymentLogRepository(db *sqlx.DB) repository.DeploymentLogRepository {
	return &PgDeploymentLogRepository{
		db: db,
	}
}

// Create saves a new deployment log entry in the database
func (r *PgDeploymentLogRepository) Create(ctx context.Context, log *model.DeploymentLog) error {
	query := `
		INSERT INTO deployment_log (deployment_id, status, metadata, action_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(
		ctx,
		query,
		log.DeploymentID,
		log.Status,
		log.Metadata,
		log.ActionAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create deployment log: %w", err)
	}

	log.ID = id
	return nil
}

// GetByID retrieves a deployment log entry by its ID
func (r *PgDeploymentLogRepository) GetByID(ctx context.Context, id int64) (*model.DeploymentLog, error) {
	query := `
		SELECT id, deployment_id, status, metadata, action_at
		FROM deployment_log
		WHERE id = $1
	`

	var log model.DeploymentLog
	err := r.db.GetContext(ctx, &log, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get deployment log by ID: %w", err)
	}

	return &log, nil
}

// GetByDeploymentID retrieves all log entries for a specific deployment
func (r *PgDeploymentLogRepository) GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentLog, error) {
	query := `
		SELECT id, deployment_id, status, metadata, action_at
		FROM deployment_log
		WHERE deployment_id = $1
		ORDER BY action_at DESC
	`

	var logs []*model.DeploymentLog
	err := r.db.SelectContext(ctx, &logs, query, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment logs by deployment ID: %w", err)
	}

	return logs, nil
}

// GetByDeploymentIDAndStatus retrieves all log entries for a specific deployment and status
func (r *PgDeploymentLogRepository) GetByDeploymentIDAndStatus(ctx context.Context, deploymentID uuid.UUID, status model.DeploymentStatus) ([]*model.DeploymentLog, error) {
	query := `
		SELECT id, deployment_id, status, metadata, action_at
		FROM deployment_log
		WHERE deployment_id = $1 AND status = $2
		ORDER BY action_at DESC
	`

	var logs []*model.DeploymentLog
	err := r.db.SelectContext(ctx, &logs, query, deploymentID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment logs by deployment ID and status: %w", err)
	}

	return logs, nil
}

// GetLatestByDeploymentID retrieves the latest log entry for a specific deployment
func (r *PgDeploymentLogRepository) GetLatestByDeploymentID(ctx context.Context, deploymentID uuid.UUID) (*model.DeploymentLog, error) {
	query := `
		SELECT id, deployment_id, status, metadata, action_at
		FROM deployment_log
		WHERE deployment_id = $1
		ORDER BY action_at DESC
		LIMIT 1
	`

	var log model.DeploymentLog
	err := r.db.GetContext(ctx, &log, query, deploymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest deployment log by deployment ID: %w", err)
	}

	return &log, nil
}

// ListLogs retrieves logs with pagination
func (r *PgDeploymentLogRepository) ListLogs(ctx context.Context, offset, limit int) ([]*model.DeploymentLog, error) {
	query := `
		SELECT id, deployment_id, status, metadata, action_at
		FROM deployment_log
		ORDER BY action_at DESC
		LIMIT $1 OFFSET $2
	`

	var logs []*model.DeploymentLog
	err := r.db.SelectContext(ctx, &logs, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployment logs: %w", err)
	}

	return logs, nil
}

// CountLogs counts total logs by deployment ID
func (r *PgDeploymentLogRepository) CountLogs(ctx context.Context, deploymentID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM deployment_log
		WHERE deployment_id = $1
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, deploymentID)
	if err != nil {
		return 0, fmt.Errorf("failed to count deployment logs: %w", err)
	}

	return count, nil
}
