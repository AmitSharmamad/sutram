package strategy

import (
	"database/sql"

	"algo/src/core/db"
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

func (repo *Repo) Create(strategy *Strategy) (bool, error) {
	_, err := repo.db.Insert(
		`
		INSERT INTO sutram.strategy (
            user_id, name, description, position_type, instruments, 
            chart_type, chart_timeframe, chart_unknown, qty, entry_expression, 
            exit_expression, target_profit, stop_loss, trailing_stop_loss, 
            exit_type, capital, max_trades, entry_start_time, entry_end_time, 
            start_date, end_date, created_at, is_active
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
        )
		`,
		strategy.UserID,
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

func (repo *Repo) Get(id string) (*Strategy, error) {
	var strategies Strategies
	err := repo.db.Select(
		&strategies,
		`
		SELECT *
		FROM sutram.strategy
		WHERE id = $2;
		`,
		repo.Name,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &strategies[0], nil
}

func (repo *Repo) GetLatest() (any, error) {
	return nil, nil
	// return repo.db.Select(
	// 	`
	// 		SELECT *
	// 		FROM ?
	// 		ORDER BY DESC
	// 		LIMIT 1;
	// 	`,
	// 	repo.Name,
	// )
}

func (this *Repo) GetAll() (Strategies, error) {
	var strategies Strategies
	this.db.Select(&strategies, "SELECT * FROM sutram.strategy LIMIT 1;")
	return strategies, nil
}

func (repo *Repo) Update(id string, strategy *Strategy) (sql.Result, error) {
	return repo.db.Update("UPDATE strategy SET ... WHERE id = {id}")
}

func (repo *Repo) Delete(query string) (sql.Result, error) {
	return repo.db.Delete("DELETE FROM strategy")
}
