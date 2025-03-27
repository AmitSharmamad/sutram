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
		Name: "position",
		db:   factory.Get("postgres"),
	}
}

func (repo *PositionRepo) Create(position *shared.Position) error {
	_, err := repo.db.Insert(`
		INSERT INTO position (
			id, run_id, state, direction
		) VALUES (
			:id, :run_id, :state, :direction
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
		WHERE run_id = $1
		ORDER BY created_at DESC;
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
		FROM position
		WHERE run_id = $1
		`,
		runID,
	)
	return &positions, err
}

func (repo *PositionRepo) Update(position *shared.Position) error {
	_, err := repo.db.Update(`
		UPDATE position
		SET
			state = :state,
			entry_price_signal = :entry_price_signal,
			entry_time_signal = :entry_time_signal,
			exit_price_signal = :exit_price_signal,
			exit_time_signal = :exit_time_signal,

			entry_price_order = :entry_price_order,
			entry_time_order = :entry_time_order,
			exit_price_order = :exit_price_order,
			exit_time_order = :exit_time_order,

			profit_loss_absolute = :profit_loss_absolute,
			profit_loss_percentage = :profit_loss_percentage,

			entry_fees = :entry_fees,
			exit_fees = :exit_fees,
			transaction_cost = :transaction_cost,

			net_profit_loss = :net_profit_loss,

			exit_reason = :exit_reason
		WHERE id = :id
		`,
		position,
	)
	return err
}
