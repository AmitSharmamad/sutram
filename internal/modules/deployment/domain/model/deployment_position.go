package model

import (
	"encoding/json"
	"time"

	strategyModel "algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
)

// PositionStatus represents the possible states of a position
type PositionStatus string

const (
	PositionOpen     PositionStatus = "OPEN"
	PositionClosed   PositionStatus = "CLOSED"
	PositionCanceled PositionStatus = "CANCELED"
)

// DeploymentPosition represents a position for a deployment
type DeploymentPosition struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	DeploymentID uuid.UUID       `json:"deployment_id" db:"deployment_id"`
	InstrumentID int             `json:"instrument_id" db:"instrument_id"`
	Quantity     float64         `json:"quantity" db:"quantity"`
	EntryPrice   float64         `json:"entry_price" db:"entry_price"`
	EntryTime    time.Time       `json:"entry_time" db:"entry_time"`
	ExitPrice    *float64        `json:"exit_price" db:"exit_price"`
	ExitTime     *time.Time      `json:"exit_time" db:"exit_time"`
	PNL          *float64        `json:"pnl" db:"pnl"`
	Status       PositionStatus  `json:"status" db:"status"`
	Metadata     json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
	// Stop loss and target profit fields for live trading (matching backtest position model)
	Direction            strategyModel.Direction `json:"direction" db:"direction"`
	StopLossPrice        *float64                `json:"stop_loss_price" db:"stop_loss_price"`
	InitialStopLossPrice *float64                `json:"initial_stop_loss_price" db:"initial_stop_loss_price"`
	TargetProfitPrice    *float64                `json:"target_profit_price" db:"target_profit_price"`
}

// NewDeploymentPosition creates a new position for a deployment
func NewDeploymentPosition(
	deploymentID uuid.UUID,
	instrumentID int,
	quantity float64,
	entryPrice float64,
	entryTime time.Time,
	direction strategyModel.Direction,
	metadata map[string]interface{},
) (*DeploymentPosition, error) {
	var metadataJSON json.RawMessage
	var err error

	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
	}

	return &DeploymentPosition{
		ID:           uuid.New(),
		DeploymentID: deploymentID,
		InstrumentID: instrumentID,
		Quantity:     quantity,
		EntryPrice:   entryPrice,
		EntryTime:    entryTime,
		Status:       PositionOpen,
		Metadata:     metadataJSON,
		Direction:    direction,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// ClosePosition closes the position with the exit details
func (p *DeploymentPosition) ClosePosition(exitPrice float64, exitTime time.Time) {
	p.ExitPrice = &exitPrice
	p.ExitTime = &exitTime
	pnl := (exitPrice - p.EntryPrice) * p.Quantity
	p.PNL = &pnl
	p.Status = PositionClosed
	p.UpdatedAt = time.Now()
}

// CancelPosition cancels the position
func (p *DeploymentPosition) CancelPosition() {
	p.Status = PositionCanceled
	p.UpdatedAt = time.Now()
}

// SetStopLossPrice sets the stop loss price for the position
func (p *DeploymentPosition) SetStopLossPrice(stopLossPrice float64) {
	p.StopLossPrice = &stopLossPrice
	if p.InitialStopLossPrice == nil {
		p.InitialStopLossPrice = &stopLossPrice
	}
	p.UpdatedAt = time.Now()
}

// SetTargetProfitPrice sets the target profit price for the position
func (p *DeploymentPosition) SetTargetProfitPrice(targetProfitPrice float64) {
	p.TargetProfitPrice = &targetProfitPrice
	p.UpdatedAt = time.Now()
}

// UpdateStopLossPrice updates the stop loss price (for trailing stops)
func (p *DeploymentPosition) UpdateStopLossPrice(newStopLossPrice float64) {
	p.StopLossPrice = &newStopLossPrice
	p.UpdatedAt = time.Now()
}
