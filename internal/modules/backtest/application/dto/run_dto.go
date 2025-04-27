package dto

import (
	"time"

	"github.com/google/uuid"
)

// RunType defines the type of a run
type RunType string

const (
	Backtest   RunType = "BACKTEST"
	PaperTrade RunType = "PAPER_TRADE"
	LiveTrade  RunType = "LIVE_TRADE"
)

// RunState defines the state of a run
type RunState string

const (
	NotStarted RunState = "NOT_STARTED"
	Paused     RunState = "PAUSED"
	Started    RunState = "STARTED"
	Ended      RunState = "ENDED"
)

// RunGroupDTO represents a data transfer object for run groups
type RunGroupDTO struct {
	ID         uuid.UUID `json:"id"`
	StrategyID uuid.UUID `json:"strategy_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RunDTO represents a data transfer object for runs
type RunDTO struct {
	ID              uuid.UUID `json:"id"`
	RunGroupID      uuid.UUID `json:"run_group_id"`
	StrategyID      uuid.UUID `json:"strategy_id"`
	Exchange        string    `json:"exchange"`
	InstrumentID    int       `json:"instrument_id"`
	Type            RunType   `json:"type"`
	State           RunState  `json:"state"`
	Qty             int       `json:"qty"`
	PeriodStartTime time.Time `json:"period_start_time"`
	PeriodEndTime   time.Time `json:"period_end_time"`
	StartedTime     time.Time `json:"started_time"`
	EndedTime       time.Time `json:"ended_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateRunGroupRequest represents a request to create a new run group
type CreateRunGroupRequest struct {
	StrategyID uuid.UUID `json:"strategy_id" binding:"required"`
}

// CreateRunRequest represents a request to create a new run
type CreateRunRequest struct {
	RunGroupID      uuid.UUID `json:"run_group_id" binding:"required"`
	StrategyID      uuid.UUID `json:"strategy_id" binding:"required"`
	Exchange        string    `json:"exchange" binding:"required"`
	InstrumentID    int       `json:"instrument_id" binding:"required"`
	Type            RunType   `json:"type" binding:"required"`
	Qty             int       `json:"qty" binding:"required"`
	PeriodStartTime time.Time `json:"period_start_time" binding:"required"`
	PeriodEndTime   time.Time `json:"period_end_time" binding:"required"`
}

// UpdateRunStateRequest represents a request to update a run's state
type UpdateRunStateRequest struct {
	State RunState `json:"state" binding:"required"`
}

// RunFilter represents a filter for searching runs
type RunFilter struct {
	StrategyID   *uuid.UUID `json:"strategy_id"`
	Type         *RunType   `json:"type"`
	State        *RunState  `json:"state"`
	Exchange     *string    `json:"exchange"`
	InstrumentID *int       `json:"instrument_id"`
}
