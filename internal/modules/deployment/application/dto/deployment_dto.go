package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateDeploymentRequest represents the request to create a new deployment
type CreateDeploymentRequest struct {
	StrategyID   uuid.UUID `json:"strategy_id" binding:"required"`
	InstrumentID int       `json:"instrument_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
}

// UpdateDeploymentRequest represents the request to update an existing deployment
type UpdateDeploymentRequest struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
}

// DeploymentResponse represents the response for a deployment
type DeploymentResponse struct {
	ID            uuid.UUID        `json:"id"`
	UserID        uuid.UUID        `json:"user_id"`
	StrategyID    uuid.UUID        `json:"strategy_id"`
	InstrumentID  int              `json:"instrument_id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	IsActive      bool             `json:"is_active"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	LatestMetrics *MetricsResponse `json:"metrics,omitempty"`
}

// DeploymentFilter represents the filter criteria for searching deployments
type DeploymentFilter struct {
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	StrategyID   *uuid.UUID `json:"strategy_id,omitempty"`
	InstrumentID *int       `json:"instrument_id,omitempty"`
	IsActive     *bool      `json:"is_active,omitempty"`
	Name         *string    `json:"name,omitempty"`
	Status       *string    `json:"status,omitempty"` // "waiting" or "entered"
}

// DeploymentStatusResponse represents the status of deployments for specific instruments against a strategy
type DeploymentStatusResponse struct {
	StrategyID  uuid.UUID                    `json:"strategy_id"`
	Instruments []DeploymentInstrumentStatus `json:"instruments"`
}

// DeploymentInstrumentStatus represents the deployment status for a specific instrument
type DeploymentInstrumentStatus struct {
	InstrumentID    int                 `json:"instrument_id"`
	HasDeployment   bool                `json:"has_deployment"`
	Deployment      *DeploymentResponse `json:"deployment,omitempty"`
	PositionSummary *PositionSummary    `json:"position_summary,omitempty"`
	LatestMetrics   *MetricsResponse    `json:"metrics,omitempty"`
}

// PositionSummary represents a summary of positions for a deployment
type PositionSummary struct {
	TotalPositions  int     `json:"total_positions"`
	OpenPositions   int     `json:"open_positions"`
	ClosedPositions int     `json:"closed_positions"`
	TotalPNL        float64 `json:"total_pnl"`
}
