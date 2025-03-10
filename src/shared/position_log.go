package shared

import (
	"time"

	"algo/src/util"

	"github.com/google/uuid"
)

type PositionLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	PositionID uuid.UUID `json:"position_id" db:"position_id"`
	State      PositionState  `json:"state" db:"state"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func NewPositionLog(positionID uuid.UUID, state PositionState) *PositionLog {
	return &PositionLog{
		ID:         util.NewUUID(),
		PositionID: positionID,
		State:      state,
	}
}
