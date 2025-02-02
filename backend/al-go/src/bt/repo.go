package bt

import (
	"algo/src/core/db"
	"algo/src/shared"

	"github.com/google/uuid"
)

// --------- Run Repo --------- //

type RunRepo struct {
	Name string
	db   db.Database
}

func NewRunRepoProv(
	factory *db.DBFactory,
) *RunRepo {
	return &RunRepo{
		Name: "sutram.bt_run",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *RunRepo) Create(run *shared.Run) error {
	_, err := repo.db.Insert(`
		INSERT INTO sutram.bt_run (
			id, strategy_id, state, qty, 
			period_start_time, period_end_time
		)
		VALUES (
			:id,
    	:strategy_id,
    	:state,
    	:qty,
    	:period_start_time,
    	:period_end_time
		)
		`,
		run,
	)
	return err
}

func (repo *RunRepo) GetByStrategyID(id string) (shared.Runs, error) {
	var runs shared.Runs
	err := repo.db.Select(&runs, `
		SELECT * 
		FROM $1
		WHERE strategy_id = $2;
		`,
		repo.Name,
		id,
	)
	return runs, err
}

func (repo *RunRepo) Get(id string) (*shared.Run, error) {
	var runs shared.Runs
	err := repo.db.Select(&runs, `
		SELECT * 
		FROM sutram.bt_run
		WHERE id = $1;
		`,
		id,
	)
	return &runs[0], err
}

func (repo *RunRepo) UpdateState(run *shared.Run) error {
	_, err := repo.db.Update(`
		UPDATE sutram.bt_run
		SET status = :state
		WHERE id = :id
		`,
		run,
	)
	return err
}

// --------- Run Position Repo --------- //

type PositionRepo struct {
	Name string
	db   db.Database
}

func NewPositionRepoProv(
	factory *db.DBFactory,
) *PositionRepo {
	return &PositionRepo{
		Name: "sutram.bt_position",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *PositionRepo) Create(position *shared.Position) error {
	_, err := repo.db.Insert(`
		INSERT INTO sutram.bt_position (
			id, run_id, state
		) VALUES (
			:id, :run_id, :state
		)
		`,
		position,
	)
	return err
}

func (repo *PositionRepo) GetLatest(runID uuid.UUID) (*shared.Position, error) {
	var positions shared.Positions
	err := repo.db.Select(&positions, `
		SELECT *
		FROM sutram.bt_position
		WHERE run_id = $1;
		`,
		runID,
	)
	if len(positions) == 0 {
		return nil, err
	}
	return &positions[0], err
}

func (repo *PositionRepo) GetByRunID(runID string) (*shared.Positions, error) {
	var positions shared.Positions
	err := repo.db.Select(&positions, `
		SELECT * 
		FROM $1
		WHERE run_id = $2
		`,
		repo.Name,
		runID,
	)
	return &positions, err
}

// --------- Run Position Log Repo --------- //

type PositionLogRepo struct {
	Name string
	db   db.Database
}

func NewPositionLogRepoProv(
	factory *db.DBFactory,
) *PositionLogRepo {
	return &PositionLogRepo{
		Name: "sutram.bt_position_log",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *PositionLogRepo) Create(log *shared.PositionLog) {
	repo.db.Insert(`
		INSERT INTO sutram.bt_position_log (
			id, position_id, state
		) VALUES (
			:id, :position_id, :state
		)
		`,
		log,
	)
}

func (repo *PositionLogRepo) GetByPositionID(id string) (*shared.PositionLogs, error) {
	var logs *shared.PositionLogs
	err := repo.db.Select(logs,
		`
		SELECT *
		FROM $1
		WHERE position_id = $2;
		`,
		repo.Name,
		id,
	)
	return logs, err
}

// -------- Run Positions -------- //

type RunPositionRepo struct {
	RunName      string
	PositionName string
	db           db.Database
}

func NewRunPositionRepoProv(
	factory *db.DBFactory,
) *RunPositionRepo {
	return &RunPositionRepo{
		RunName:      `sutram.bt_run`,
		PositionName: `sutram.bt_position`,
	}
}

// func (repo *RunPositionRepo) Get(id string) (Runs, error) {
// 	var runs RunPosition
//
// 	err := repo.db.Select(&runs, `
// 		SELECT *
// 		FROM $1
// 		JOIN $2
// 		ON $1.id = $2.run_id
// 		ORDER BY timestamp DESC;
// 		`,
// 		repo.RunName,
// 		repo.PositionName,
// 	)
// 	return runs, err
// }
