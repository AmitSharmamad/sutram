package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/strategy/domain/model"
	"algo/internal/modules/strategy/domain/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Common errors
var (
	ErrOptionLegNotFound      = errors.New("option leg not found")
	ErrOptionLegDatabaseError = errors.New("option leg database error")
)

// OptionLegRepository implements the repository.OptionLegRepository interface
type OptionLegRepository struct {
	db *sql.DB
}

// NewOptionLegRepository creates a new PostgreSQL repository for option legs
func NewOptionLegRepository(db *sql.DB) repository.OptionLegRepository {
	return &OptionLegRepository{db: db}
}

// Create saves a new option leg in the database
func (r *OptionLegRepository) Create(ctx context.Context, optionLeg *model.OptionLeg) error {
	query := `
		INSERT INTO option_leg (
			id, strategy_id, leg_sequence, select_index, segment, option_type,
			action_type, strike_type, strike_price, total_lot,
			lot_size, order_type, tick_offset, expiry_type,
			expiry_cycle, expiry_date, exchange, asset, 
			exchange_security_id, instrument_type, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)`

	_, err := r.db.ExecContext(ctx, query,
		optionLeg.ID, optionLeg.StrategyID, optionLeg.LegSequence, optionLeg.SelectIndex,
		optionLeg.Segment, optionLeg.OptionType, optionLeg.ActionType, optionLeg.StrikeType,
		optionLeg.StrikePrice, optionLeg.TotalLot, optionLeg.LotSize,
		optionLeg.OrderType, optionLeg.TickOffset, optionLeg.ExpiryType,
		optionLeg.ExpiryCycle, optionLeg.ExpiryDate, optionLeg.Exchange, optionLeg.Asset,
		optionLeg.ExchangeSecurityID, optionLeg.InstrumentType, optionLeg.CreatedAt,
		optionLeg.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("option leg with this leg sequence already exists for the strategy")
			case "23503": // foreign_key_violation
				return fmt.Errorf("strategy with ID %s does not exist", optionLeg.StrategyID)
			}
		}
		return fmt.Errorf("failed to create option leg: %w", err)
	}

	return nil
}

// GetByID retrieves an option leg by its ID
func (r *OptionLegRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OptionLeg, error) {
	query := `
		SELECT id, strategy_id, leg_sequence, select_index, segment, option_type,
			   action_type, strike_type, strike_price, total_lot,
			   lot_size, order_type, tick_offset, expiry_type,
			   expiry_cycle, expiry_date, exchange, asset, 
			   exchange_security_id, instrument_type, created_at, updated_at
		FROM option_leg 
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	optionLeg, err := r.scanOptionLeg(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get option leg by ID: %w", err)
	}

	return optionLeg, nil
}

// GetByStrategyID retrieves all option legs for a specific strategy
func (r *OptionLegRepository) GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.OptionLeg, error) {
	query := `
		SELECT id, strategy_id, leg_sequence, select_index, segment, option_type,
			   action_type, strike_type, strike_price, total_lot,
			   lot_size, order_type, tick_offset, expiry_type,
			   expiry_cycle, expiry_date, exchange, asset, 
			   exchange_security_id, instrument_type, created_at, updated_at
		FROM option_leg 
		WHERE strategy_id = $1 ORDER BY leg_sequence`

	rows, err := r.db.QueryContext(ctx, query, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get option legs by strategy ID: %w", err)
	}
	defer rows.Close()

	return r.scanOptionLegs(rows)
}

// GetByStrategyIDAndSequence retrieves an option leg by strategy ID and leg sequence
func (r *OptionLegRepository) GetByStrategyIDAndSequence(ctx context.Context, strategyID uuid.UUID, legSequence int) (*model.OptionLeg, error) {
	query := `
		SELECT id, strategy_id, leg_sequence, select_index, segment, option_type,
			   action_type, strike_type, strike_price, total_lot,
			   lot_size, order_type, tick_offset, expiry_type,
			   expiry_cycle, expiry_date, exchange, asset, 
			   exchange_security_id, instrument_type, created_at, updated_at
		FROM option_leg 
		WHERE strategy_id = $1 AND leg_sequence = $2`

	row := r.db.QueryRowContext(ctx, query, strategyID, legSequence)
	optionLeg, err := r.scanOptionLeg(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get option leg by strategy ID and sequence: %w", err)
	}

	return optionLeg, nil
}

// Update updates an existing option leg in the database
func (r *OptionLegRepository) Update(ctx context.Context, optionLeg *model.OptionLeg) error {
	query := `
		UPDATE option_leg SET
			leg_sequence = $2, select_index = $3, segment = $4, option_type = $5,
			action_type = $6, strike_type = $7, strike_price = $8,
			total_lot = $9, lot_size = $10, order_type = $11,
			tick_offset = $12, expiry_type = $13, expiry_cycle = $14, expiry_date = $15,
			exchange = $16, asset = $17, exchange_security_id = $18, 
			instrument_type = $19, updated_at = $20
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		optionLeg.ID, optionLeg.LegSequence, optionLeg.SelectIndex, optionLeg.Segment,
		optionLeg.OptionType, optionLeg.ActionType, optionLeg.StrikeType, optionLeg.StrikePrice,
		optionLeg.TotalLot, optionLeg.LotSize,
		optionLeg.OrderType, optionLeg.TickOffset, optionLeg.ExpiryType, optionLeg.ExpiryCycle, optionLeg.ExpiryDate,
		optionLeg.Exchange, optionLeg.Asset,
		optionLeg.ExchangeSecurityID, optionLeg.InstrumentType, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update option leg: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrOptionLegNotFound
	}

	return nil
}

// Delete removes an option leg from the database
func (r *OptionLegRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM option_leg WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete option leg: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrOptionLegNotFound
	}

	return nil
}

// DeleteByStrategyID removes all option legs for a specific strategy
func (r *OptionLegRepository) DeleteByStrategyID(ctx context.Context, strategyID uuid.UUID) error {
	query := `DELETE FROM option_leg WHERE strategy_id = $1`

	_, err := r.db.ExecContext(ctx, query, strategyID)
	if err != nil {
		return fmt.Errorf("failed to delete option legs by strategy ID: %w", err)
	}

	return nil
}

// GetByIndex retrieves all option legs for a specific index/underlying
func (r *OptionLegRepository) GetByIndex(ctx context.Context, selectIndex string) ([]*model.OptionLeg, error) {
	query := `
		SELECT id, strategy_id, leg_sequence, select_index, segment, option_type,
			   action_type, strike_type, strike_price, total_lot,
			   lot_size, order_type, tick_offset, expiry_type,
			   expiry_cycle, expiry_date, exchange, asset, 
			   exchange_security_id, instrument_type, created_at, updated_at
		FROM option_leg 
		WHERE select_index = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, selectIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get option legs by index: %w", err)
	}
	defer rows.Close()

	return r.scanOptionLegs(rows)
}

// GetByExpiryDate retrieves all option legs with a specific expiry date
func (r *OptionLegRepository) GetByExpiryDate(ctx context.Context, expiryDate string) ([]*model.OptionLeg, error) {
	query := `
		SELECT id, strategy_id, leg_sequence, select_index, segment, option_type,
			   action_type, strike_type, strike_price, total_lot,
			   lot_size, order_type, tick_offset, expiry_type,
			   expiry_cycle, expiry_date, exchange, asset, 
			   exchange_security_id, instrument_type, created_at, updated_at
		FROM option_leg 
		WHERE expiry_date::text = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, expiryDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get option legs by expiry date: %w", err)
	}
	defer rows.Close()

	return r.scanOptionLegs(rows)
}

// List retrieves all option legs with pagination
func (r *OptionLegRepository) List(ctx context.Context, offset, limit int) ([]*model.OptionLeg, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM option_leg`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, strategy_id, leg_sequence, select_index, segment, option_type,
			   action_type, strike_type, strike_price, total_lot,
			   lot_size, order_type, tick_offset, expiry_type,
			   expiry_cycle, expiry_date, exchange, asset, 
			   exchange_security_id, instrument_type, created_at, updated_at
		FROM option_leg 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list option legs: %w", err)
	}
	defer rows.Close()

	optionLegs, err := r.scanOptionLegs(rows)
	if err != nil {
		return nil, 0, err
	}

	return optionLegs, total, nil
}

// Helper method to scan a single option leg from a row
func (r *OptionLegRepository) scanOptionLeg(row *sql.Row) (*model.OptionLeg, error) {
	var optionLeg model.OptionLeg
	var optionType sql.NullString
	var strikePrice sql.NullFloat64
	var lotSize sql.NullInt32
	var tickOffset sql.NullFloat64
	var expiryDate sql.NullTime
	var exchangeSecurityID sql.NullInt32

	err := row.Scan(
		&optionLeg.ID, &optionLeg.StrategyID, &optionLeg.LegSequence, &optionLeg.SelectIndex,
		&optionLeg.Segment, &optionType, &optionLeg.ActionType, &optionLeg.StrikeType,
		&strikePrice, &optionLeg.TotalLot, &lotSize,
		&optionLeg.OrderType, &tickOffset, &optionLeg.ExpiryType, &optionLeg.ExpiryCycle, &expiryDate,
		&optionLeg.Exchange, &optionLeg.Asset,
		&exchangeSecurityID, &optionLeg.InstrumentType, &optionLeg.CreatedAt, &optionLeg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if optionType.Valid {
		ot := model.OptionType(optionType.String)
		optionLeg.OptionType = &ot
	}
	if strikePrice.Valid {
		optionLeg.StrikePrice = &strikePrice.Float64
	}
	if lotSize.Valid {
		size := int(lotSize.Int32)
		optionLeg.LotSize = &size
	}
	if tickOffset.Valid {
		optionLeg.TickOffset = &tickOffset.Float64
	}
	if expiryDate.Valid {
		optionLeg.ExpiryDate = &expiryDate.Time
	}
	if exchangeSecurityID.Valid {
		id := int(exchangeSecurityID.Int32)
		optionLeg.ExchangeSecurityID = &id
	}

	return &optionLeg, nil
}

// Helper method to scan multiple option legs from rows
func (r *OptionLegRepository) scanOptionLegs(rows *sql.Rows) ([]*model.OptionLeg, error) {
	var optionLegs []*model.OptionLeg

	for rows.Next() {
		var optionLeg model.OptionLeg
		var optionType sql.NullString
		var strikePrice sql.NullFloat64
		var lotSize sql.NullInt32
		var tickOffset sql.NullFloat64
		var expiryDate sql.NullTime
		var exchangeSecurityID sql.NullInt32

		err := rows.Scan(
			&optionLeg.ID, &optionLeg.StrategyID, &optionLeg.LegSequence, &optionLeg.SelectIndex,
			&optionLeg.Segment, &optionType, &optionLeg.ActionType, &optionLeg.StrikeType,
			&strikePrice, &optionLeg.TotalLot, &lotSize,
			&optionLeg.OrderType, &tickOffset, &optionLeg.ExpiryType, &optionLeg.ExpiryCycle, &expiryDate,
			&optionLeg.Exchange, &optionLeg.Asset,
			&exchangeSecurityID, &optionLeg.InstrumentType, &optionLeg.CreatedAt, &optionLeg.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan option leg: %w", err)
		}

		// Handle nullable fields
		if optionType.Valid {
			ot := model.OptionType(optionType.String)
			optionLeg.OptionType = &ot
		}
		if strikePrice.Valid {
			optionLeg.StrikePrice = &strikePrice.Float64
		}
		if lotSize.Valid {
			size := int(lotSize.Int32)
			optionLeg.LotSize = &size
		}
		if tickOffset.Valid {
			optionLeg.TickOffset = &tickOffset.Float64
		}
		if expiryDate.Valid {
			optionLeg.ExpiryDate = &expiryDate.Time
		}
		if exchangeSecurityID.Valid {
			id := int(exchangeSecurityID.Int32)
			optionLeg.ExchangeSecurityID = &id
		}

		optionLegs = append(optionLegs, &optionLeg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over option leg rows: %w", err)
	}

	return optionLegs, nil
}
