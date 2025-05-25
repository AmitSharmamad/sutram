package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PgDeploymentPositionLogRepository implements the deployment position log repository interface
type PgDeploymentPositionLogRepository struct {
	db *sqlx.DB
}

// NewPgDeploymentPositionLogRepository creates a new Postgres deployment position log repository
func NewPgDeploymentPositionLogRepository(db *sqlx.DB) repository.DeploymentPositionLogRepository {
	return &PgDeploymentPositionLogRepository{
		db: db,
	}
}

// Create saves a new position log in the database
func (r *PgDeploymentPositionLogRepository) Create(ctx context.Context, positionLog *model.DeploymentPositionLog) error {
	query := `
		INSERT INTO deployment_position_log (
			deployment_position_id, results, timestamp, created_at
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(
		ctx,
		query,
		positionLog.DeploymentPositionID,
		positionLog.Results,
		positionLog.Timestamp,
		positionLog.CreatedAt,
	).Scan(&id)

	if err != nil {
		return err
	}

	positionLog.ID = id
	return nil
}

// GetByID retrieves a position log by its ID
func (r *PgDeploymentPositionLogRepository) GetByID(ctx context.Context, id string) (*model.DeploymentPositionLog, error) {
	query := `
		SELECT id, deployment_position_id, results, timestamp, created_at
		FROM deployment_position_log
		WHERE id = $1
	`

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var positionLog model.DeploymentPositionLog
	err = r.db.GetContext(ctx, &positionLog, query, parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &positionLog, nil
}

// GetByPositionID retrieves all logs for a specific position
func (r *PgDeploymentPositionLogRepository) GetByPositionID(ctx context.Context, positionID int64) ([]*model.DeploymentPositionLog, error) {
	query := `
		SELECT id, deployment_position_id, results, timestamp, created_at
		FROM deployment_position_log
		WHERE deployment_position_id = $1
		ORDER BY timestamp DESC
	`

	logs := []*model.DeploymentPositionLog{}
	err := r.db.SelectContext(ctx, &logs, query, positionID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByPositionIDAndTimeRange retrieves logs for a position within a specific time range
func (r *PgDeploymentPositionLogRepository) GetByPositionIDAndTimeRange(ctx context.Context, positionID int64, startTime, endTime time.Time) ([]*model.DeploymentPositionLog, error) {
	query := `
		SELECT id, deployment_position_id, results, timestamp, created_at
		FROM deployment_position_log
		WHERE deployment_position_id = $1
		AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC
	`

	logs := []*model.DeploymentPositionLog{}
	err := r.db.SelectContext(ctx, &logs, query, positionID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByDeploymentID retrieves all logs for all positions in a deployment
func (r *PgDeploymentPositionLogRepository) GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentPositionLog, error) {
	query := `
		SELECT l.id, l.deployment_position_id, l.results, l.timestamp, l.created_at
		FROM deployment_position_log l
		JOIN deployment_position p ON l.deployment_position_id = p.id
		WHERE p.deployment_id = $1
		ORDER BY l.timestamp DESC
	`

	logs := []*model.DeploymentPositionLog{}
	err := r.db.SelectContext(ctx, &logs, query, deploymentID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
