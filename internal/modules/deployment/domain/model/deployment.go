package model

import (
	"time"

	"github.com/google/uuid"
)

// Deployment represents a user strategy deployed against a specific instrument
type Deployment struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	StrategyID   uuid.UUID `json:"strategy_id" db:"strategy_id"`
	InstrumentID int       `json:"instrument_id" db:"instrument_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// NewDeployment creates a new deployment with the given parameters
func NewDeployment(
	userID uuid.UUID,
	strategyID uuid.UUID,
	instrumentID int,
	name string,
	description string,
) *Deployment {
	return &Deployment{
		ID:           uuid.New(),
		UserID:       userID,
		StrategyID:   strategyID,
		InstrumentID: instrumentID,
		Name:         name,
		Description:  description,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Update updates the deployment fields
func (d *Deployment) Update(name string, description string) {
	d.Name = name
	d.Description = description
	d.UpdatedAt = time.Now()
}

// Activate activates the deployment
func (d *Deployment) Activate() {
	d.IsActive = true
	d.UpdatedAt = time.Now()
}

// Deactivate deactivates the deployment
func (d *Deployment) Deactivate() {
	d.IsActive = false
	d.UpdatedAt = time.Now()
}
