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

func (repo *RunRepo) Create(run *shared.Run) error {
	_, err := repo.db.Insert(`
		INSERT INTO run (
			id, strategy_id, instrument, state, qty,
			period_start_time, period_end_time
		)
		VALUES (
			:id,
			:strategy_id,
			:instrument,
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
		SELECT id, strategy_id, instrument, run_group_id, type
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
		UPDATE sutram.bt_run
		SET status = :state
		WHERE id = :id
		`,
		run,
	)
	return err
}
