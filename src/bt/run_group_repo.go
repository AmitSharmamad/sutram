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
			id,
			strategy_id
		)
		VALUES (
			:id,
			:strategy_id
		)
		`,
		runGroup,
	)
	return err
}
