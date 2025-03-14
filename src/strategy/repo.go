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
			chart_unknown,
			qty,

			entry_expression,
			exit_expression,
			target_profit,
			stop_loss,
			trailing_stop_loss,
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
			:chart_unknown,
			:qty,

			:entry_expression,
			:exit_expression,
			:target_profit,
			:stop_loss,
			:trailing_stop_loss,
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
