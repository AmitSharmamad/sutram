package bt

import (
	"algo/src/core/db"
	"algo/src/shared"

	"github.com/google/uuid"
)

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
		INSERT INTO position (
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
		FROM position
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
