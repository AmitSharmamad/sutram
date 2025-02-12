package shared

import (
	"time"

	"algo/src/util"

	"github.com/google/uuid"
)

type Positions []Position

type PositionLogs []PositionLog

type RunPositions []RunPosition

type PositionType string

const (
	LATEST PositionType = "LATEST"
	NEW    PositionType = "NEW"
)

type Position struct {
	ID         uuid.UUID `json:"id" db:"id"`
	RunID      uuid.UUID `json:"run_id" db:"run_id"`
	State      RunState  `json:"state" db:"state"`
	EntryPrice float64   `json:"entry_price" db:"entry_price"`
	EntryTime  time.Time `json:"entry_time" db:"entry_time"`
	ExitPrice  float64   `json:"exit_price" db:"exit_price"`
	ExitTime   time.Time `json:"exit_time" db:"exit_time"`
	ExitReason string    `json:"exit_reason" db:"exit_reason"`
	Audit
}

func NewPosition(runID uuid.UUID) *Position {
	return &Position{
		ID:    util.NewUUID(),
		RunID: runID,
		State: "NOT_ENTERED",
	}
}

func (position *Position) IsDone() bool {
	return position.State == Exited
}

type RunPosition struct {
	Run       *Run       `json:"run"`
	Positions *Positions `json:"positions"`
}

func NewRunPosition() *RunPosition {
	return new(RunPosition)
}
