package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeploymentPositionLog represents a log entry for a deployment position's data point execution results
type DeploymentPositionLog struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	DeploymentPositionID uuid.UUID           `json:"deployment_position_id" db:"deployment_position_id"`
	Results              json.RawMessage `json:"results" db:"results"`
	Timestamp            time.Time       `json:"timestamp" db:"timestamp"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
}

// NewDeploymentPositionLog creates a new deployment position log instance
func NewDeploymentPositionLog(deploymentPositionID uuid.UUID, results json.RawMessage, timestamp time.Time) *DeploymentPositionLog {
	return &DeploymentPositionLog{
		DeploymentPositionID: deploymentPositionID,
		Results:              results,
		Timestamp:            timestamp,
		CreatedAt:            time.Now(),
	}
}
