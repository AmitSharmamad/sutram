package clickhouse

import (
	"context"
	"fmt"
	"strings"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"

	"github.com/jmoiron/sqlx"
)

// InstrumentMasterRepository is a ClickHouse implementation of the InstrumentMasterRepository interface
type InstrumentMasterRepository struct {
	db *sqlx.DB
}

// NewInstrumentMasterRepository creates a new ClickHouse instrument master repository
func NewInstrumentMasterRepository(db *sqlx.DB) repository.InstrumentMasterRepository {
	return &InstrumentMasterRepository{
		db: db,
	}
}

// FindByExchangeSecurityID finds an instrument master by exchange security ID (this is now our primary lookup)
func (r *InstrumentMasterRepository) FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int64) (*model.InstrumentMaster, error) {
	query := `
		SELECT
			exchange_security_id,
			exchange,
			segment,
			symbol,
			instrument,
			start_date,
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
		WHERE exchange_security_id = ?
		LIMIT 1
	`

	instrumentMaster := &model.InstrumentMaster{}
	err := r.db.GetContext(ctx, instrumentMaster, query, exchangeSecurityID)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument master by exchange security ID: %w", err)
	}

	return instrumentMaster, nil
}

// FindBySymbol finds an instrument master by symbol
func (r *InstrumentMasterRepository) FindBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error) {
	query := `
		SELECT 
			exchange_security_id,
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
		WHERE symbol = ?
		LIMIT 1
	`

	instrumentMaster := &model.InstrumentMaster{}
	err := r.db.GetContext(ctx, instrumentMaster, query, symbol)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument master by symbol: %w", err)
	}

	return instrumentMaster, nil
}

// FindByExchange finds instrument masters by exchange
func (r *InstrumentMasterRepository) FindByExchange(ctx context.Context, exchange string) ([]*model.InstrumentMaster, error) {
	query := `
		SELECT 
			exchange_security_id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
		WHERE exchange = ?
	`

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query, exchange)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument masters by exchange: %w", err)
	}

	return instrumentMasters, nil
}

// FindBySegment finds instrument masters by segment
func (r *InstrumentMasterRepository) FindBySegment(ctx context.Context, segment string) ([]*model.InstrumentMaster, error) {
	query := `
		SELECT 
			exchange_security_id,
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
		WHERE segment = ?
	`

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query, segment)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument masters by segment: %w", err)
	}

	return instrumentMasters, nil
}

// FindAll returns all instrument masters
func (r *InstrumentMasterRepository) FindAll(ctx context.Context) ([]*model.InstrumentMaster, error) {
	query := `
		SELECT 
			exchange_security_id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
	`

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query)
	if err != nil {
		return nil, fmt.Errorf("error querying all instrument masters: %w", err)
	}

	return instrumentMasters, nil
}

// Create creates a new instrument master
func (r *InstrumentMasterRepository) Create(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	query := `
		INSERT INTO instrument_master (
			exchange_security_id,
			exchange,
			segment,
			symbol,
			instrument,
			start_date,
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrumentMaster.ExchangeSecurityID,
		instrumentMaster.Exchange,
		instrumentMaster.Segment,
		instrumentMaster.Symbol,
		instrumentMaster.Name,
		instrumentMaster.StartDate,
		instrumentMaster.Expiry,
		instrumentMaster.ExmExchID,
		instrumentMaster.SmstSecurityID,
		instrumentMaster.InstrumentName,
		instrumentMaster.TradingSymbol,
		instrumentMaster.LotUnits,
		instrumentMaster.StrikePrice,
		instrumentMaster.OptionType,
		instrumentMaster.TickSize,
		instrumentMaster.ExpiryFlag,
		instrumentMaster.ExchInstrumentType,
		instrumentMaster.SymbolName,
	)

	if err != nil {
		return fmt.Errorf("error creating instrument master: %w", err)
	}

	return nil
}

// Update updates an instrument master by exchange security ID
func (r *InstrumentMasterRepository) Update(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	// Note: ClickHouse typically doesn't support direct updates as in traditional RDBMS
	// This is a simplified version assuming you're using ALTER TABLE UPDATE syntax
	query := `
		ALTER TABLE instrument_master UPDATE
			exchange = ?,
			segment = ?,
			symbol = ?,
			instrument = ?,
			start_date = ?,
			expiry = ?,
			exm_exch_id = ?,
			smst_security_id = ?,
			instrument_name = ?,
			trading_symbol = ?,
			lot_units = ?,
			strike_price = ?,
			option_type = ?,
			tick_size = ?,
			expiry_flag = ?,
			exch_instrument_type = ?,
			symbol_name = ?
		WHERE exchange_security_id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrumentMaster.Exchange,
		instrumentMaster.Segment,
		instrumentMaster.Symbol,
		instrumentMaster.Name,
		instrumentMaster.StartDate,
		instrumentMaster.Expiry,
		instrumentMaster.ExmExchID,
		instrumentMaster.SmstSecurityID,
		instrumentMaster.InstrumentName,
		instrumentMaster.TradingSymbol,
		instrumentMaster.LotUnits,
		instrumentMaster.StrikePrice,
		instrumentMaster.OptionType,
		instrumentMaster.TickSize,
		instrumentMaster.ExpiryFlag,
		instrumentMaster.ExchInstrumentType,
		instrumentMaster.SymbolName,
		instrumentMaster.ExchangeSecurityID,
	)

	if err != nil {
		return fmt.Errorf("error updating instrument master: %w", err)
	}

	return nil
}

// Delete deletes an instrument master by exchange security ID
func (r *InstrumentMasterRepository) Delete(ctx context.Context, exchangeSecurityID int64) error {
	query := `
		ALTER TABLE instrument_master DELETE WHERE exchange_security_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, exchangeSecurityID)
	if err != nil {
		return fmt.Errorf("error deleting instrument master: %w", err)
	}

	return nil
}

// FindWithFilters returns instrument masters with filtering and pagination
func (r *InstrumentMasterRepository) FindWithFilters(ctx context.Context, exchange, segment, searchQuery string, offset, limit int) ([]*model.InstrumentMaster, int, error) {
	// Build the WHERE clause based on filters
	var whereConditions []string
	var args []interface{}

	if exchange != "" {
		whereConditions = append(whereConditions, "exchange = ?")
		args = append(args, exchange)
	}

	if segment != "" {
		whereConditions = append(whereConditions, "segment = ?")
		args = append(args, segment)
	}

	if searchQuery != "" {
		// Search in both name and symbol fields
		whereConditions = append(whereConditions, "(lowerUTF8(symbol) LIKE ? OR lowerUTF8(instrument) LIKE ? OR lowerUTF8(instrument_name) LIKE ? OR lowerUTF8(trading_symbol) LIKE ? OR lowerUTF8(symbol_name) LIKE ?)")
		searchPattern := "%" + strings.ToLower(searchQuery) + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Build the full WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Query to get filtered and paginated data
	sqlQuery := fmt.Sprintf(`
		SELECT 
			exchange_security_id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
		%s
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add limit and offset to the arguments
	args = append(args, limit, offset)

	// Execute the query
	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, sqlQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying instrument masters with filters: %w", err)
	}

	// Count query to get total records matching the filter
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sutram.instrument_master
		%s
	`, whereClause)

	// Execute the count query
	var totalCount int
	err = r.db.GetContext(ctx, &totalCount, countQuery, args[:len(args)-2]...) // Exclude limit and offset from count query
	if err != nil {
		return nil, 0, fmt.Errorf("error counting instrument masters with filters: %w", err)
	}

	return instrumentMasters, totalCount, nil
}

// FindByExchangeSecurityIDs finds multiple instrument masters by their exchange security IDs
func (r *InstrumentMasterRepository) FindByExchangeSecurityIDs(ctx context.Context, exchangeSecurityIDs []int64) ([]*model.InstrumentMaster, error) {
	if len(exchangeSecurityIDs) == 0 {
		return []*model.InstrumentMaster{}, nil
	}

	// Convert the slice of IDs to a comma-separated string for the IN clause
	placeholders := make([]string, len(exchangeSecurityIDs))
	args := make([]interface{}, len(exchangeSecurityIDs))
	for i, id := range exchangeSecurityIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT 
			exchange_security_id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date,
			expiry,
			exm_exch_id,
			smst_security_id,
			instrument_name,
			trading_symbol,
			lot_units,
			strike_price,
			option_type,
			tick_size,
			expiry_flag,
			exch_instrument_type,
			symbol_name
		FROM sutram.instrument_master
		WHERE exchange_security_id IN (%s)
	`, strings.Join(placeholders, ","))

	instrumentMasters := []*model.InstrumentMaster{}
	err := r.db.SelectContext(ctx, &instrumentMasters, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument masters by exchange security IDs: %w", err)
	}

	return instrumentMasters, nil
}
