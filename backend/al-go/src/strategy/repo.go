package strategy

import (
	"context"

	"algo/src/core/db"
	"algo/src/shared"
)

type Strategies []Strategy

type Repo struct {
	Name string
	db   db.Database
}

func NewRepoProv(
	factory *db.DBFactory,
) *Repo {
	return &Repo{
		Name: "sutram.strategy",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *Repo) Create(ctx context.Context, strategy *Strategy) (bool, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	_, err := repo.db.Insert(
		`
		INSERT INTO $1 (
            user_id, name, description, position_type, instruments, 
            chart_type, chart_timeframe, chart_unknown, qty, entry_expression, 
            exit_expression, target_profit, stop_loss, trailing_stop_loss, 
            exit_type, capital, max_trades, entry_start_time, entry_end_time, 
            start_date, end_date, created_at, is_active
        ) VALUES (
             $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
        )
		`,
		repo.Name,
		userID,
		strategy.Name,
		strategy.Description,
		strategy.PositionType,
		strategy.Instruments,
		strategy.ChartType,
		strategy.ChartTimeframe,
		strategy.ChartUnknown,
		strategy.Qty,
		strategy.EntryExpression,
		strategy.ExitExpression,
		strategy.TargetProfit,
		strategy.StopLoss,
		strategy.TrailingStopLoss,
		strategy.ExitType,
		strategy.Capital,
		strategy.MaxTrades,
		strategy.EntryStartTime,
		strategy.EntryEndTime,
		strategy.StartDate,
		strategy.EndDate,
		strategy.CreatedAt,
		strategy.IsActive,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *Repo) Get(
	ctx context.Context,
	id string,
) (*Strategy, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies Strategies

	err := repo.db.Select(
		&strategies,
		`
		SELECT *
		FROM $1
		WHERE user_id = $2 
		AND id = $3;
		`,
		repo.Name,
		userID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &strategies[0], nil
}

func (repo *Repo) GetLatest(
	ctx context.Context,
) (*Strategy, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies Strategies

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
) (Strategies, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies Strategies

	repo.db.Select(
		&strategies,
		"SELECT * FROM $1 WHERE user_id = $2;",
		repo.Name,
		userID,
	)

	return strategies, nil
}
