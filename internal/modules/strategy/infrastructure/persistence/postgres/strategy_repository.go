package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"algo/internal/modules/strategy/domain/model"
	"algo/internal/modules/strategy/domain/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Common errors
var (
	ErrStrategyNotFound = errors.New("strategy not found")
	ErrDatabaseError    = errors.New("database error")
)

// StrategyRepository implements the repository.StrategyRepository interface
type StrategyRepository struct {
	db *sql.DB
}

// NewStrategyRepository creates a new PostgreSQL repository for strategies
func NewStrategyRepository(db *sql.DB) repository.StrategyRepository {
	return &StrategyRepository{
		db: db,
	}
}

// Create saves a new strategy in the database
func (r *StrategyRepository) Create(ctx context.Context, strategy *model.Strategy) error {
	query := `
		INSERT INTO strategy (
			id, user_id, name, description, position_type, instruments, 
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, brokerage, slippage, entry_start_time, entry_end_time, 
			start_date, end_date, segment, strategy_type, 
			entry_time, exit_time, created_at, updated_at, is_active, max_allocation
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, 
			$27, $28, $29, $30, $31, $32, $33, $34, $35, $36
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		strategy.ID,
		strategy.UserID,
		strategy.Name,
		strategy.Description,
		strategy.PositionType,
		pq.Array(strategy.Instruments),
		strategy.ChartType,
		strategy.ChartTimeframe,
		strategy.TradingDays,
		strategy.TradeType,
		strategy.Qty,
		strategy.EntryExpression,
		strategy.ExitExpression,
		strategy.TargetProfit,
		strategy.TargetProfitType,
		strategy.StopLoss,
		strategy.StopLossType,
		strategy.TrailingStopLoss,
		strategy.TrailingStopLossType,
		strategy.ExitType,
		strategy.Capital,
		strategy.MaxTrades,
		strategy.Brokerage,
		strategy.Slippage,
		strategy.EntryStartTime,
		strategy.EntryEndTime,
		strategy.StartDate,
		strategy.EndDate,
		strategy.Segment,
		strategy.StrategyType,
		strategy.EntryTime,
		strategy.ExitTime,
		strategy.CreatedAt,
		strategy.UpdatedAt,
		strategy.IsActive,
		strategy.MaxAllocation,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID retrieves a strategy by its ID
func (r *StrategyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
	query := `
		SELECT 
			id, user_id, name, description, position_type, instruments,
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, brokerage, slippage, entry_start_time, entry_end_time, 
			start_date, end_date, segment, strategy_type, 
			entry_time, exit_time, created_at, updated_at, is_active, max_allocation
		FROM strategy
		WHERE id = $1
	`

	var strategy model.Strategy
	var instruments pq.Int32Array
	var positionType sql.NullString
	var tradeType sql.NullString
	var qty sql.NullInt32
	var capital sql.NullFloat64
	var maxAllocation sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&strategy.ID,
		&strategy.UserID,
		&strategy.Name,
		&strategy.Description,
		&positionType,
		&instruments,
		&strategy.ChartType,
		&strategy.ChartTimeframe,
		&strategy.TradingDays,
		&tradeType,
		&qty,
		&strategy.EntryExpression,
		&strategy.ExitExpression,
		&strategy.TargetProfit,
		&strategy.TargetProfitType,
		&strategy.StopLoss,
		&strategy.StopLossType,
		&strategy.TrailingStopLoss,
		&strategy.TrailingStopLossType,
		&strategy.ExitType,
		&capital,
		&strategy.MaxTrades,
		&strategy.Brokerage,
		&strategy.Slippage,
		&strategy.EntryStartTime,
		&strategy.EntryEndTime,
		&strategy.StartDate,
		&strategy.EndDate,
		&strategy.Segment,
		&strategy.StrategyType,
		&strategy.EntryTime,
		&strategy.ExitTime,
		&strategy.CreatedAt,
		&strategy.UpdatedAt,
		&strategy.IsActive,
		&maxAllocation,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStrategyNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Handle nullable fields
	if positionType.Valid {
		pt := model.Direction(positionType.String)
		strategy.PositionType = &pt
	}
	if tradeType.Valid {
		tt := model.TradeType(tradeType.String)
		strategy.TradeType = &tt
	}
	if qty.Valid {
		q := int(qty.Int32)
		strategy.Qty = &q
	}
	if capital.Valid {
		strategy.Capital = &capital.Float64
	}
	if maxAllocation.Valid {
		strategy.MaxAllocation = &maxAllocation.Float64
	}

	strategy.Instruments = instruments

	return &strategy, nil
}

// GetByIDs retrieves multiple strategies by their IDs
func (r *StrategyRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Strategy, error) {
	if len(ids) == 0 {
		return []*model.Strategy{}, nil
	}

	// Build placeholder string for the query
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, name, description, position_type, instruments,
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, entry_start_time, entry_end_time, 
			start_date, end_date, created_at, updated_at, is_active, max_allocation
		FROM strategy
		WHERE id IN (%s)
		ORDER BY created_at DESC;
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var strategies []*model.Strategy

	for rows.Next() {
		var strategy model.Strategy
		var instruments pq.Int32Array

		err := rows.Scan(
			&strategy.ID,
			&strategy.UserID,
			&strategy.Name,
			&strategy.Description,
			&strategy.PositionType,
			&instruments,
			&strategy.ChartType,
			&strategy.ChartTimeframe,
			&strategy.TradingDays,
			&strategy.TradeType,
			&strategy.Qty,
			&strategy.EntryExpression,
			&strategy.ExitExpression,
			&strategy.TargetProfit,
			&strategy.TargetProfitType,
			&strategy.StopLoss,
			&strategy.StopLossType,
			&strategy.TrailingStopLoss,
			&strategy.TrailingStopLossType,
			&strategy.ExitType,
			&strategy.Capital,
			&strategy.MaxTrades,
			&strategy.EntryStartTime,
			&strategy.EntryEndTime,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.CreatedAt,
			&strategy.UpdatedAt,
			&strategy.IsActive,
			&strategy.MaxAllocation,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		strategy.Instruments = instruments
		strategies = append(strategies, &strategy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return strategies, nil
}

// GetByUserID retrieves all strategies for a specific user
func (r *StrategyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error) {
	query := `
		SELECT 
			id, user_id, name, description, position_type, instruments, 
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, brokerage, slippage, entry_start_time, entry_end_time, 
			start_date, end_date, segment, strategy_type, 
			entry_time, exit_time, created_at, updated_at, is_active, max_allocation
		FROM strategy
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var strategies []*model.Strategy

	for rows.Next() {
		var strategy model.Strategy
		var instruments pq.Int32Array

		err := rows.Scan(
			&strategy.ID,
			&strategy.UserID,
			&strategy.Name,
			&strategy.Description,
			&strategy.PositionType,
			&instruments,
			&strategy.ChartType,
			&strategy.ChartTimeframe,
			&strategy.TradingDays,
			&strategy.TradeType,
			&strategy.Qty,
			&strategy.EntryExpression,
			&strategy.ExitExpression,
			&strategy.TargetProfit,
			&strategy.TargetProfitType,
			&strategy.StopLoss,
			&strategy.StopLossType,
			&strategy.TrailingStopLoss,
			&strategy.TrailingStopLossType,
			&strategy.ExitType,
			&strategy.Capital,
			&strategy.MaxTrades,
			&strategy.Brokerage,
			&strategy.Slippage,
			&strategy.EntryStartTime,
			&strategy.EntryEndTime,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Segment,
			&strategy.StrategyType,
			&strategy.EntryTime,
			&strategy.ExitTime,
			&strategy.CreatedAt,
			&strategy.UpdatedAt,
			&strategy.IsActive,
			&strategy.MaxAllocation,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		strategy.Instruments = instruments
		strategies = append(strategies, &strategy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return strategies, nil
}

// Update updates an existing strategy in the database
func (r *StrategyRepository) Update(ctx context.Context, strategy *model.Strategy) error {
	query := `
		UPDATE strategy
		SET 
			name = $1, 
			description = $2, 
			position_type = $3, 
			instruments = $4, 
			chart_type = $5, 
			chart_timeframe = $6, 
			trading_days = $7,
			trade_type = $8, 
			qty = $9, 
			entry_expression = $10, 
			exit_expression = $11, 
			target_profit = $12, 
			target_profit_type = $13, 
			stop_loss = $14, 
			stop_loss_type = $15, 
			trailing_stop_loss = $16, 
			trailing_stop_loss_type = $17, 
			exit_type = $18, 
			capital = $19, 
			max_trades = $20, 
			brokerage = $21,
			slippage = $22,
			entry_start_time = $23, 
			entry_end_time = $24, 
			start_date = $25, 
			end_date = $26, 
			segment = $27,
			strategy_type = $28,
			entry_time = $29,
			exit_time = $30,
			updated_at = $31, 
			is_active = $32,
			max_allocation = $33
		WHERE id = $34
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		strategy.Name,
		strategy.Description,
		strategy.PositionType,
		pq.Array(strategy.Instruments),
		strategy.ChartType,
		strategy.ChartTimeframe,
		strategy.TradingDays,
		strategy.TradeType,
		strategy.Qty,
		strategy.EntryExpression,
		strategy.ExitExpression,
		strategy.TargetProfit,
		strategy.TargetProfitType,
		strategy.StopLoss,
		strategy.StopLossType,
		strategy.TrailingStopLoss,
		strategy.TrailingStopLossType,
		strategy.ExitType,
		strategy.Capital,
		strategy.MaxTrades,
		strategy.Brokerage,
		strategy.Slippage,
		strategy.EntryStartTime,
		strategy.EntryEndTime,
		strategy.StartDate,
		strategy.EndDate,
		strategy.Segment,
		strategy.StrategyType,
		strategy.EntryTime,
		strategy.ExitTime,
		time.Now(),
		strategy.IsActive,
		strategy.MaxAllocation,
		strategy.ID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrStrategyNotFound
	}

	return nil
}

// Delete removes a strategy from the database
func (r *StrategyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM strategy WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrStrategyNotFound
	}

	return nil
}

// ActivateStrategy sets the strategy as active
func (r *StrategyRepository) ActivateStrategy(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE strategy SET is_active = true, updated_at = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrStrategyNotFound
	}

	return nil
}

// DeactivateStrategy sets the strategy as inactive
func (r *StrategyRepository) DeactivateStrategy(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE strategy SET is_active = false, updated_at = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrStrategyNotFound
	}

	return nil
} // GetActiveStrategies retrieves all active strategies
func (r *StrategyRepository) GetActiveStrategies(ctx context.Context) ([]*model.Strategy, error) {
	query := `
		SELECT 
			id, user_id, name, description, position_type, instruments, 
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, brokerage, slippage, entry_start_time, entry_end_time, 
			start_date, end_date, segment, strategy_type, 
			entry_time, exit_time, created_at, updated_at, is_active, max_allocation
		FROM strategy
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var strategies []*model.Strategy

	for rows.Next() {
		var strategy model.Strategy
		var instruments pq.Int32Array

		err := rows.Scan(
			&strategy.ID,
			&strategy.UserID,
			&strategy.Name,
			&strategy.Description,
			&strategy.PositionType,
			pq.Array(&instruments),
			&strategy.ChartType,
			&strategy.ChartTimeframe,
			&strategy.TradingDays,
			&strategy.TradeType,
			&strategy.Qty,
			&strategy.EntryExpression,
			&strategy.ExitExpression,
			&strategy.TargetProfit,
			&strategy.TargetProfitType,
			&strategy.StopLoss,
			&strategy.StopLossType,
			&strategy.TrailingStopLoss,
			&strategy.TrailingStopLossType,
			&strategy.ExitType,
			&strategy.Capital,
			&strategy.MaxTrades,
			&strategy.Brokerage,
			&strategy.Slippage,
			&strategy.EntryStartTime,
			&strategy.EntryEndTime,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Segment,
			&strategy.StrategyType,
			&strategy.EntryTime,
			&strategy.ExitTime,
			&strategy.CreatedAt,
			&strategy.UpdatedAt,
			&strategy.IsActive,
			&strategy.MaxAllocation,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		strategy.Instruments = instruments
		strategies = append(strategies, &strategy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return strategies, nil
} // GetActiveStrategiesByUserID retrieves all active strategies for a specific user
func (r *StrategyRepository) GetActiveStrategiesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error) {
	query := `
		SELECT 
			id, user_id, name, description, position_type, instruments, 
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, brokerage, slippage, entry_start_time, entry_end_time, 
			start_date, end_date, segment, strategy_type, 
			entry_time, exit_time, created_at, updated_at, is_active, max_allocation
		FROM strategy
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var strategies []*model.Strategy

	for rows.Next() {
		var strategy model.Strategy
		var instruments pq.Int32Array

		err := rows.Scan(
			&strategy.ID,
			&strategy.UserID,
			&strategy.Name,
			&strategy.Description,
			&strategy.PositionType,
			pq.Array(&instruments),
			&strategy.ChartType,
			&strategy.ChartTimeframe,
			&strategy.TradingDays,
			&strategy.TradeType,
			&strategy.Qty,
			&strategy.EntryExpression,
			&strategy.ExitExpression,
			&strategy.TargetProfit,
			&strategy.TargetProfitType,
			&strategy.StopLoss,
			&strategy.StopLossType,
			&strategy.TrailingStopLoss,
			&strategy.TrailingStopLossType,
			&strategy.ExitType,
			&strategy.Capital,
			&strategy.MaxTrades,
			&strategy.Brokerage,
			&strategy.Slippage,
			&strategy.EntryStartTime,
			&strategy.EntryEndTime,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Segment,
			&strategy.StrategyType,
			&strategy.EntryTime,
			&strategy.ExitTime,
			&strategy.CreatedAt,
			&strategy.UpdatedAt,
			&strategy.IsActive,
			&strategy.MaxAllocation,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		strategy.Instruments = instruments
		strategies = append(strategies, &strategy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return strategies, nil
}

// SearchStrategies searches strategies with filter criteria
func (r *StrategyRepository) SearchStrategies(ctx context.Context, filters map[string]interface{}) ([]*model.Strategy, error) {
	// Start with base query
	query := `
		SELECT 
			id, user_id, name, description, position_type, instruments, 
			chart_type, chart_timeframe, trading_days, trade_type, qty, entry_expression, 
			exit_expression, target_profit, target_profit_type, stop_loss, 
			stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, 
			exit_type, capital, max_trades, brokerage, slippage, entry_start_time, entry_end_time, 
			start_date, end_date, created_at, updated_at, is_active, max_allocation
		FROM strategy
		WHERE 1=1
	`

	// Dynamically build the where clause based on filters
	args := []interface{}{}
	argPosition := 1
	var whereClauses []string

	// Process filters
	for key, value := range filters {
		switch key {
		case "user_id":
			whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argPosition))
			args = append(args, value)
			argPosition++
		case "name":
			whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argPosition))
			args = append(args, "%"+value.(string)+"%")
			argPosition++
		case "position_type":
			whereClauses = append(whereClauses, fmt.Sprintf("position_type = $%d", argPosition))
			args = append(args, value)
			argPosition++
		case "instrument":
			whereClauses = append(whereClauses, fmt.Sprintf("$%d = ANY(instruments)", argPosition))
			args = append(args, value)
			argPosition++
		case "chart_type":
			whereClauses = append(whereClauses, fmt.Sprintf("chart_type = $%d", argPosition))
			args = append(args, value)
			argPosition++
		case "chart_timeframe":
			whereClauses = append(whereClauses, fmt.Sprintf("chart_timeframe = $%d", argPosition))
			args = append(args, value)
			argPosition++
		case "is_active":
			whereClauses = append(whereClauses, fmt.Sprintf("is_active = $%d", argPosition))
			args = append(args, value)
			argPosition++
		case "created_after":
			whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argPosition))
			args = append(args, value)
			argPosition++
		case "created_before":
			whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argPosition))
			args = append(args, value)
			argPosition++
		}
	}

	// Append the where clauses to the query if any
	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var strategies []*model.Strategy

	for rows.Next() {
		var strategy model.Strategy
		var instruments pq.Int32Array

		err := rows.Scan(
			&strategy.ID,
			&strategy.UserID,
			&strategy.Name,
			&strategy.Description,
			&strategy.PositionType,
			pq.Array(&instruments),
			&strategy.ChartType,
			&strategy.ChartTimeframe,
			&strategy.TradingDays,
			&strategy.TradeType,
			&strategy.Qty,
			&strategy.EntryExpression,
			&strategy.ExitExpression,
			&strategy.TargetProfit,
			&strategy.TargetProfitType,
			&strategy.StopLoss,
			&strategy.StopLossType,
			&strategy.TrailingStopLoss,
			&strategy.TrailingStopLossType,
			&strategy.ExitType,
			&strategy.Capital,
			&strategy.MaxTrades,
			&strategy.Brokerage,
			&strategy.Slippage,
			&strategy.EntryStartTime,
			&strategy.EntryEndTime,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.CreatedAt,
			&strategy.UpdatedAt,
			&strategy.IsActive,
			&strategy.MaxAllocation,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		strategy.Instruments = instruments
		strategies = append(strategies, &strategy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return strategies, nil
}
