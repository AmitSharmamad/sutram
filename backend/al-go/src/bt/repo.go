package bt

import (
	"algo/src/core/db"
	"algo/src/shared"
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

func (repo *RunRepo) Create(strategyID string) error {
	_, err := repo.db.Insert(`
		INSERT INTO $1
		VALUES (
			$2
		)
		`, repo.Name,
		strategyID,
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
		FROM $1
		WHERE strategy_id = $2;
		`,
		repo.Name,
		id,
	)
	return &runs[0], err
}

func (repo *RunRepo) UpdateStatus(id string, status RunState) error {
	_, err := repo.db.Update(`
		UPDATE $1
		SET status = $2
		WHERE id = $3;
		`,
		repo.Name,
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

func (repo *PositionRepo) GetLatest(id string) (*shared.Position, error) {
	var positions shared.Positions
	err := repo.db.Select(&positions, `
		SELECT *
		FROM $1
		ORDER BY created_date DESC;
		`,
		repo.Name,
	)
	return &positions[0], err
}

func (repo *PositionRepo) GetByRunID(id string) (*shared.Positions, error) {
	var positions shared.Positions
	err := repo.db.Select(&positions, `
		SELECT * 
		FROM $1
		WHERE run_id = $2
		`,
		repo.Name,
		id,
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
