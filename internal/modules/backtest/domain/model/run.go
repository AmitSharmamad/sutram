package model

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

// RunGroup represents a group of related runs for a strategy
type RunGroup struct {
	ID         uuid.UUID `json:"id"`
	StrategyID uuid.UUID `json:"strategy_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Run represents a backtest, paper trade, or live trade run
type Run struct {
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

// NewRunGroup creates a new run group
func NewRunGroup(strategyID uuid.UUID) *RunGroup {
	now := time.Now()
	return &RunGroup{
		ID:         uuid.New(),
		StrategyID: strategyID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// NewRun creates a new run
func NewRun(
	runGroupID uuid.UUID,
	strategyID uuid.UUID,
	exchange string,
	instrumentID int,
	runType RunType,
	qty int,
	periodStartTime time.Time,
	periodEndTime time.Time,
) *Run {
	now := time.Now()
	return &Run{
		ID:              uuid.New(),
		RunGroupID:      runGroupID,
		StrategyID:      strategyID,
		Exchange:        exchange,
		InstrumentID:    instrumentID,
		Type:            runType,
		State:           NotStarted,
		Qty:             qty,
		PeriodStartTime: periodStartTime,
		PeriodEndTime:   periodEndTime,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
