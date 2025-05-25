package dto

import (
	"encoding/json"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentLogRequest is the request for creating a deployment log
type DeploymentLogRequest struct {
	DeploymentID uuid.UUID              `json:"deployment_id" validate:"required"`
	Status       model.DeploymentStatus `json:"status" validate:"required"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// DeploymentLogResponse is the response for a deployment log
type DeploymentLogResponse struct {
	ID           uuid.UUID                  `json:"id"`
	DeploymentID uuid.UUID              `json:"deployment_id"`
	Status       model.DeploymentStatus `json:"status"`
	Metadata     json.RawMessage        `json:"metadata,omitempty"`
	ActionAt     time.Time              `json:"action_at"`
}

// ToDeploymentLogResponse converts a model.DeploymentLog to dto.DeploymentLogResponse
func ToDeploymentLogResponse(log *model.DeploymentLog) *DeploymentLogResponse {
	return &DeploymentLogResponse{
		ID:           log.ID,
		DeploymentID: log.DeploymentID,
		Status:       log.Status,
		Metadata:     log.Metadata,
		ActionAt:     log.ActionAt,
	}
}

// ToDeploymentLogResponses converts a slice of model.DeploymentLog to dto.DeploymentLogResponse
func ToDeploymentLogResponses(logs []*model.DeploymentLog) []*DeploymentLogResponse {
	responses := make([]*DeploymentLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = ToDeploymentLogResponse(log)
	}
	return responses
}
