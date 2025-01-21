package bt

import "time"

type BTRun struct {
	ID              int       `json:"id" db:"id"`                               // Primary Key
	StrategyID      int       `json:"strategy_id" db:"strategy_id"`             // Foreign key to strategy(id)
	State           string    `json:"state" db:"state"`                         // possible values: "start", "in_progress", "ended"
	Qty             float64   `json:"qty" db:"qty"`                             // Quantity for the run
	PeriodStartTime time.Time `json:"period_start_time" db:"period_start_time"` // Start time of the period
	PeriodEndTime   time.Time `json:"period_end_time" db:"period_end_time"`     // End time of the period
	StartedTime     time.Time `json:"started_time" db:"started_time"`           // When the run started
	EndedTime       time.Time `json:"ended_time" db:"ended_time"`               // When the run ended
	CreatedAt       time.Time `json:"created_at" db:"created_at"`               // Timestamp when the record was created
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`               // Timestamp when the record was last updated
}

func Get(id string) *BTRun {
	return &BTRun{}
}

func GetLatest() *BTRun {
	return &BTRun{}
}

func Create(run *BTRun) *BTRun {
	return &BTRun{}
}

func Update(run *BTRun) *BTRun {
	return &BTRun{}
}

func Delete(id string) bool {
	return true
}
