package strategy

import (
	"context"

	"algo/src/core/db"
	"algo/src/shared"
	"algo/src/util"

	"github.com/google/uuid"
)

type Repo struct {
	Name string
	db   db.Database
}

func NewRepoProv(
	factory *db.DBFactory,
) *Repo {
	return &Repo{
		Name: "strategy",
		db:   factory.Get("postgres"),
	}
}

func (repo *Repo) Create(ctx context.Context, strategy *shared.Strategy) error {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	strategy.UserID = util.Parse(userID)

	_, err := repo.db.Insert(
		`
		INSERT INTO strategy (
			id, 
			user_id,

			name,
			description,
			position_type,
			instruments,
			chart_type,
			chart_timeframe,
			qty,

			entry_expression,
			exit_expression,
			target_profit,
			target_profit_type,
			stop_loss,
			stop_loss_type,
			trailing_stop_loss,
			trailing_stop_loss_type,
			exit_type,

			capital,
			max_trades,
			entry_start_time,
			entry_end_time,
			start_date,
			end_date,
			is_active
		) VALUES (
			:id, 
			:user_id, 

			:name, 
			:description, 
			:position_type,
			:instruments,
			:chart_type,
			:chart_timeframe,
			:qty,

			:entry_expression,
			:exit_expression,
			:target_profit,
			:target_profit_type,
			:stop_loss,
			:stop_loss_type,
			:trailing_stop_loss,
			:trailing_stop_loss_type,
			:exit_type,

			:capital,
			:max_trades,
			:entry_start_time,
			:entry_end_time,
			:start_date,
			:end_date,
			:is_active
		);
		`,
		strategy,
	)
	return err
}

func (repo *Repo) Get(
	ctx context.Context,
	id uuid.UUID,
) (*shared.Strategy, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategy shared.Strategy

	err := repo.db.Get(
		&strategy,
		`
		SELECT *
		FROM strategy
		WHERE user_id = $1
		AND id = $2;
		`,
		userID,
		id,
	)

	return &strategy, err
}
func (repo *Repo) GetLatest(
	ctx context.Context,
) (*shared.Strategy, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies shared.Strategies

	err := repo.db.Select(
		&strategies,
		`
			SELECT *
			FROM $1
			WHERE user_id = $2
			ORDER BY timestamp DESC
			LIMIT 1;
		`,
		repo.Name,
		userID,
	)

	return &strategies[0], err
}

func (repo *Repo) GetAll(
	ctx context.Context,
) (*shared.Strategies, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies shared.Strategies

	err := repo.db.Select(
		&strategies,
		"SELECT * FROM strategy WHERE user_id = $1;",
		userID,
	)

	if err != nil {
		return nil, err
	}

	return &strategies, nil
}

func (repo *Repo) Update(ctx context.Context, strategy *shared.Strategy) error {
	_, err := repo.db.Update(`
		UPDATE strategy 
		SET
			name = :name,
			description = :description,
			position_type = :position_type,
			instruments = :instruments,
			chart_type = :chart_type,
			chart_timeframe = :chart_timeframe,
			qty = :qty,
			entry_expression = :entry_expression,
			exit_expression = :exit_expression,
			target_profit = :target_profit,
			target_profit_type = :target_profit_type,
			stop_loss = :stop_loss,
			stop_loss_type = :stop_loss_type,
			trailing_stop_loss = :trailing_stop_loss,
			trailing_stop_loss_type = :trailing_stop_loss_type,
			exit_type = :exit_type,
			capital = :capital,
			max_trades = :max_trades,
			entry_start_time = :entry_start_time,
			entry_end_time = :entry_end_time,
			start_date = :start_date,
			end_date = :end_date,
			is_active = :is_active
		WHERE id = :id;
	`, strategy)
	return err
}

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Delete(`
		DELETE FROM strategy 
		WHERE id = :id;
	`, map[string]interface{}{
		"id": id,
	})
	return err
}
