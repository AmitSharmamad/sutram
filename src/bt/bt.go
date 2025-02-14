package bt

type RunState string

const (
	Start      RunState = "START"
	InProgress RunState = "IN_PROGRESS"
	Ended      RunState = "ENDED"
)

// type Run shared.Run
//
// type Runs []Run
//
// type Position shared.Position
//
// type Positions []Position
//
// type RunPosition struct {
// 	Run       *Run       `json:"run"`
// 	Positions *Positions `json:"positions"`
// }
//
// type RunPositions []RunPosition
//
// type PositionLog shared.PositionLog
//
// type PositionLogs []PositionLog

/*

type Run struct {
	ID              int       `json:"id" db:"id"`
	StrategyID      string    `json:"strategy_id" db:"strategy_id"`
	State           string    `json:"state" db:"state"`
	Qty             float64   `json:"qty" db:"qty"`
	PeriodStartTime time.Time `json:"period_start_time" db:"period_start_time"`
	PeriodEndTime   time.Time `json:"period_end_time" db:"period_end_time"`
	StartedTime     time.Time `json:"started_time" db:"started_time"`
	EndedTime       time.Time `json:"ended_time" db:"ended_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func (run *Run) Run() {
	/// get the data range from the strategy
	s := strategy.Get(run.StrategyID)
	s.GetEntryDataRange()

	/// fetch the data
	data := strategy.GetData()

	/// start taking positions
	/// begin from period start date and iterate till the end
	for i := 0; i < len(data); i++ {
		position := run.GetPosition()
		position.Eval(data)
		position.IsDone()
	}
}

func (run *Run) GetPosition() *position.Position {
	if run.State == "START" {
		run.State = "IN_PROGRESS"
		return position.NewPosition()
	} else if run.State == "IN_PROGRESS" {
		latest := position.GetLatest()
		if latest.State == "EXITED" {
			return position.NewPosition()
		}
		return latest
	} else if run.State == "ENDED" {
		panic("Run ended, but loop going on")
	}
	panic("Error in run get position")
}

*/
