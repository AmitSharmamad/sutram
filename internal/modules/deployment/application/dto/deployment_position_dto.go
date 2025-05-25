package dto

import (
	"encoding/json"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// CreatePositionRequest is the request for creating a deployment position
type CreatePositionRequest struct {
	DeploymentID uuid.UUID              `json:"deployment_id" validate:"required"`
	InstrumentID int                    `json:"instrument_id" validate:"required"`
	Quantity     float64                `json:"quantity" validate:"required,gt=0"`
	EntryPrice   float64                `json:"entry_price" validate:"required,gt=0"`
	EntryTime    time.Time              `json:"entry_time" validate:"required"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ClosePositionRequest is the request for closing a position
type ClosePositionRequest struct {
	ExitPrice float64   `json:"exit_price" validate:"required,gt=0"`
	ExitTime  time.Time `json:"exit_time" validate:"required"`
}

// PositionResponse is the response for a deployment position
type PositionResponse struct {
	ID           uuid.UUID                `json:"id"`
	DeploymentID uuid.UUID            `json:"deployment_id"`
	InstrumentID int                  `json:"instrument_id"`
	Quantity     float64              `json:"quantity"`
	EntryPrice   float64              `json:"entry_price"`
	EntryTime    time.Time            `json:"entry_time"`
	ExitPrice    *float64             `json:"exit_price,omitempty"`
	ExitTime     *time.Time           `json:"exit_time,omitempty"`
	PNL          *float64             `json:"pnl,omitempty"`
	Status       model.PositionStatus `json:"status"`
	Metadata     json.RawMessage      `json:"metadata,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

// ToPositionResponse converts a model.DeploymentPosition to dto.PositionResponse
func ToPositionResponse(position *model.DeploymentPosition) *PositionResponse {
	return &PositionResponse{
		ID:           position.ID,
		DeploymentID: position.DeploymentID,
		InstrumentID: position.InstrumentID,
		Quantity:     position.Quantity,
		EntryPrice:   position.EntryPrice,
		EntryTime:    position.EntryTime,
		ExitPrice:    position.ExitPrice,
		ExitTime:     position.ExitTime,
		PNL:          position.PNL,
		Status:       position.Status,
		Metadata:     position.Metadata,
		CreatedAt:    position.CreatedAt,
		UpdatedAt:    position.UpdatedAt,
	}
}

// ToPositionResponses converts a slice of model.DeploymentPosition to dto.PositionResponse
func ToPositionResponses(positions []*model.DeploymentPosition) []*PositionResponse {
	responses := make([]*PositionResponse, len(positions))
	for i, position := range positions {
		responses[i] = ToPositionResponse(position)
	}
	return responses
}
