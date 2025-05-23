package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeploymentStatus represents the possible states of a deployment
type DeploymentStatus string

const (
	StatusCreated   DeploymentStatus = "CREATED"
	StatusStarted   DeploymentStatus = "STARTED"
	StatusStopped   DeploymentStatus = "STOPPED"
	StatusFailed    DeploymentStatus = "FAILED"
	StatusCompleted DeploymentStatus = "COMPLETED"
)

// DeploymentLog represents a log entry for a deployment
type DeploymentLog struct {
	ID           uuid.UUID            `json:"id" db:"id"`
	DeploymentID uuid.UUID        `json:"deployment_id" db:"deployment_id"`
	Status       DeploymentStatus `json:"status" db:"status"`
	Metadata     json.RawMessage  `json:"metadata" db:"metadata"`
	ActionAt     time.Time        `json:"action_at" db:"action_at"`
}

// NewDeploymentLog creates a new deployment log entry
func NewDeploymentLog(
	deploymentID uuid.UUID,
	status DeploymentStatus,
	metadata map[string]interface{},
) (*DeploymentLog, error) {
	var metadataJSON json.RawMessage
	var err error

	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
	}

	return &DeploymentLog{
		DeploymentID: deploymentID,
		Status:       status,
		Metadata:     metadataJSON,
		ActionAt:     time.Now(),
	}, nil
}
