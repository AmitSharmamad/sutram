package bt

import (
	"algo/src/core/db"
	"algo/src/shared"

	"github.com/google/uuid"
)

type RunMetricsRepo struct {
	Name string
	db   db.Database
}

func NewRunMetricsRepoProv(
	factory *db.DBFactory,
) *RunMetricsRepo {
	return &RunMetricsRepo{
		Name: "run_metrics",
		db:   factory.Get("postgres"),
	}
}

func (repo *RunMetricsRepo) Create(runMetrics *shared.RunMetrics) error {
	_, err := repo.db.Insert(
		`
		INSERT INTO run_metrics (
			id,
			run_id,
			total_return_pct,
			annualized_return,
			risk_adjusted_return,
			sharpe_ratio,
			sortino_ratio,
			max_drawdown,
			max_drawdown_duration,
			value_at_risk,
			std_deviation,
			downside_deviation,
			total_trades,
			win_rate,
			avg_win_size,
			avg_loss_size,
			avg_trade_duration,
			profit_factor,
			avg_position_size,
			max_position_size,
			position_turnover,
			avg_exposure,
			total_commission,
			total_slippage,
			transaction_costs_pct,
			metric_json,
			created_at,
			updated_at
		) VALUES (
			:id,
			:run_id,
			:total_return_pct,
			:annualized_return,
			:risk_adjusted_return,
			:sharpe_ratio,
			:sortino_ratio,
			:max_drawdown,
			:max_drawdown_duration,
			:value_at_risk,
			:std_deviation,
			:downside_deviation,
			:total_trades,
			:win_rate,
			:avg_win_size,
			:avg_loss_size,
			:avg_trade_duration,
			:profit_factor,
			:avg_position_size,
			:max_position_size,
			:position_turnover,
			:avg_exposure,
			:total_commission,
			:total_slippage,
			:transaction_costs_pct,
			:metric_json,
			NOW(),
			NOW()
		)
		ON CONFLICT (run_id) DO UPDATE SET
			total_return_pct = :total_return_pct,
			annualized_return = :annualized_return,
			risk_adjusted_return = :risk_adjusted_return,
			sharpe_ratio = :sharpe_ratio,
			sortino_ratio = :sortino_ratio,
			max_drawdown = :max_drawdown,
			max_drawdown_duration = :max_drawdown_duration,
			value_at_risk = :value_at_risk,
			std_deviation = :std_deviation,
			downside_deviation = :downside_deviation,
			total_trades = :total_trades,
			win_rate = :win_rate,
			avg_win_size = :avg_win_size,
			avg_loss_size = :avg_loss_size,
			avg_trade_duration = :avg_trade_duration,
			profit_factor = :profit_factor,
			avg_position_size = :avg_position_size,
			max_position_size = :max_position_size,
			position_turnover = :position_turnover,
			avg_exposure = :avg_exposure,
			total_commission = :total_commission,
			total_slippage = :total_slippage,
			transaction_costs_pct = :transaction_costs_pct,
			metric_json = :metric_json,
			updated_at = NOW();
		`,
		runMetrics,
	)
	return err
}

func (repo *RunMetricsRepo) Get(runId uuid.UUID) (*shared.RunMetrics, error) {
	var metrics []*shared.RunMetrics

	err := repo.db.Select(
		&metrics,
		"SELECT * FROM run_metrics WHERE run_id = :run_id",
		map[string]interface{}{
			"run_id": runId,
		},
	)

	if err != nil {
		panic(err)
	}

	if len(metrics) == 0 {
		panic("No metrics found")
	}

	return metrics[0], nil
}
