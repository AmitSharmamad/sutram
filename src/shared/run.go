package shared

import (
	"time"

	"github.com/google/uuid"
)

const (
	NotEntered RunState = "NOT_ENTERED"
	Entered    RunState = "ENTERED"
	Exited     RunState = "EXITED"
)

type RunState string

const (
	BackTest   RunType = "BACKTEST"
	PaperTrade RunType = "PAPER_TRADE"
	LiveTrade  RunType = "LIVE_TRADE"
)

type RunType string

type Runs []Run

type Run struct {
	ID              uuid.UUID `json:"id" db:"id"`
	StrategyID      uuid.UUID `json:"strategy_id" db:"strategy_id"`
	Instrument      string    `json:"instrument" db:"instrument"`
	RunGroupId      uuid.UUID `json:"run_group_id" db:"run_group_id"`
	State           RunState  `json:"state" db:"state"`
	RunType         RunType   `json:"run_type" db:"type"`
	Qty             float64   `json:"qty" db:"qty"`
	PeriodStartTime time.Time `json:"period_start_time" db:"period_start_time"`
	PeriodEndTime   time.Time `json:"period_end_time" db:"period_end_time"`
	StartedTime     time.Time `json:"started_time" db:"started_time"`
	EndedTime       time.Time `json:"ended_time" db:"ended_time"`
	Audit
}

func NewRun() *Run {
	return new(Run)
}

func (run *Run) GetPositionType() PositionType {
	switch run.State {
	case "START":
		run.State = "IN_PROGRESS"
		return NEW
	case "IN_PROGRESS":
		// latest := position.GetLatest()
		// if latest.State == "EXITED" {
		// 	// create a new position
		// 	return position.NewPosition()
		// }
		// return latest
		return LATEST
	case "ENDED":
		panic("Run ended, but loop going on")
	}
	panic("Error in run get position")
}
