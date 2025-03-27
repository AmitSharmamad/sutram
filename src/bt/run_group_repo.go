package bt

import (
	"algo/src/core/db"
	"algo/src/shared"
)

type RunGroupRepo struct {
	Name string
	db   db.Database
}

func NewRunGroupRepoProv(
	factory *db.DBFactory,
) *RunGroupRepo {
	return &RunGroupRepo{
		Name: "sutram.run_group",
		db:   factory.Get("postgres"),
	}
}

func (repo *RunGroupRepo) Create(runGroup *shared.RunGroup) error {
	_, err := repo.db.Insert(`
		INSERT INTO run_group (
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
		runGroup,
	)
	return err
}
