package bt

import (
	"algo/src/core/db"
	"algo/src/shared"
)

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
