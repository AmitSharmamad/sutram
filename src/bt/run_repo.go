package bt

import (
	"algo/src/core/db"
	"algo/src/shared"
)

type RunRepo struct {
	Name string
	db   db.Database
}

func NewRunRepoProv(
	factory *db.DBFactory,
) *RunRepo {
	return &RunRepo{
		Name: "sutram.bt_run",
		db:   factory.Get("postgres"),
	}
}

func (repo *RunRepo) Create(runs []*shared.Run) error {
	_, err := repo.db.Insert(`
		INSERT INTO run (
			id, strategy_id, instrument, "type", exchange, state, qty,
			period_start_time, period_end_time
		)
		VALUES (
			:id,
			:strategy_id,
			:instrument,
			:type,
			:exchange,
			:state,
			:qty,
			:period_start_time,
			:period_end_time
		)
		`,
		runs,
	)
	return err
}

func (repo *RunRepo) GetByStrategyID(id string) (shared.Runs, error) {
	var runs shared.Runs
	err := repo.db.Select(&runs, `
		SELECT id, run_group_id, strategy_id, exchange, instrument,
		"type", "state", qty,
		period_start_time, period_end_time, started_time, ended_time,
		created_at, updated_at
		FROM run
		WHERE strategy_id = $1;
		`,
		id,
	)
	return runs, err
}

func (repo *RunRepo) Get(id string) (*shared.Run, error) {
	var runs shared.Runs
	err := repo.db.Select(&runs, `
		SELECT * 
		FROM run
		WHERE id = $1;
		`,
		id,
	)
	return &runs[0], err
}

func (repo *RunRepo) UpdateState(run *shared.Run) error {
	_, err := repo.db.Update(`
		UPDATE run
		SET
			state = :state,
			started_time = :started_time,
			ended_time = :ended_time
		WHERE id = :id
		`,
		run,
	)
	return err
}
