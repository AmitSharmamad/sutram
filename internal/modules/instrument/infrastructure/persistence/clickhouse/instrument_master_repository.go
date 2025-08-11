package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"

	"github.com/jmoiron/sqlx"
)

// InstrumentMasterRepository implements the repository interface
type InstrumentMasterRepository struct {
	db *sqlx.DB
}

// NewInstrumentMasterRepository creates a new instance
func NewInstrumentMasterRepository(db *sqlx.DB) repository.InstrumentMasterRepository {
	return &InstrumentMasterRepository{db: db}
}

// FindByExchangeSecurityID finds an instrument master by exchange security ID
func (r *InstrumentMasterRepository) FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) (*model.InstrumentMaster, error) {
	return nil, fmt.Errorf("FindByExchangeSecurityID not implemented")
}

// FindBySymbol finds an instrument master by symbol
func (r *InstrumentMasterRepository) FindBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error) {
	return nil, fmt.Errorf("FindBySymbol not implemented")
}

// FindAll returns all instrument masters
func (r *InstrumentMasterRepository) FindAll(ctx context.Context) ([]*model.InstrumentMaster, error) {
	return []*model.InstrumentMaster{}, nil
}

// FindWithFilters returns instrument masters with filtering and pagination
func (r *InstrumentMasterRepository) FindWithFilters(ctx context.Context, exchange, segment, query string, offset, limit int) ([]*model.InstrumentMaster, int, error) {
	// Build the base query
	whereClause := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	// Add exchange filter
	if exchange != "" {
		whereClause = append(whereClause, fmt.Sprintf("exchange = ?"))
		args = append(args, exchange)
		argIndex++
	}

	// Add segment filter
	if segment != "" {
		whereClause = append(whereClause, fmt.Sprintf("segment = ?"))
		args = append(args, segment)
		argIndex++
	}

	// Add search query filter (search in symbol, instrument_name)
	if query != "" {
		whereClause = append(whereClause, fmt.Sprintf("(symbol ILIKE ? OR instrument_name ILIKE ?)"))
		searchTerm := "%" + query + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// Build WHERE clause
	whereStr := strings.Join(whereClause, " AND ")

	// First, get the total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM sutram.instrument_master
		WHERE %s
	`, whereStr)

	var totalCount int
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting instrument masters: %w", err)
	}

	// If no records, return early
	if totalCount == 0 {
		return []*model.InstrumentMaster{}, 0, nil
	}

	// Get paginated data
	selectQuery := fmt.Sprintf(`
		SELECT 
			exchange_security_id,
			exchange,
			segment,
			symbol,
			instrument,
			instrument_name,
			isin,
			underlying_security_id,
			underlying_symbol,
			expiry,
			expiry_flag,
			lot_units,
			tick_size,
			strike_price,
			option_type,
			instrument_type,
			bracket_flag,
			cover_flag,
			asm_gsm_flag,
			created_at,
			updated_at
		FROM sutram.instrument_master
		WHERE %s
		ORDER BY exchange, segment, symbol
		LIMIT ? OFFSET ?
	`, whereStr)

	// Add limit and offset to args
	args = append(args, limit, offset)

	var instrumentMasters []*model.InstrumentMaster
	err = r.db.SelectContext(ctx, &instrumentMasters, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying instrument masters: %w", err)
	}

	return instrumentMasters, totalCount, nil
}

// Create creates a new instrument master
func (r *InstrumentMasterRepository) Create(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	return fmt.Errorf("Create not implemented")
}

// Update updates an instrument master
func (r *InstrumentMasterRepository) Update(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	return fmt.Errorf("Update not implemented")
}

// Delete deletes an instrument master
func (r *InstrumentMasterRepository) Delete(ctx context.Context, exchangeSecurityID int32) error {
	return fmt.Errorf("Delete not implemented")
}

// FindByExchangeSecurityIDs finds instruments by multiple exchange security IDs
func (r *InstrumentMasterRepository) FindByExchangeSecurityIDs(ctx context.Context, exchangeSecurityIDs []int32) ([]*model.InstrumentMaster, error) {
	return nil, fmt.Errorf("FindByExchangeSecurityIDs not implemented")
}

// FindOptionsContracts finds option contracts based on filters
func (r *InstrumentMasterRepository) FindOptionsContracts(
	ctx context.Context,
	strikePrice float64,
	optionType string,
	exchange string,
	symbol string,
	expiryDate time.Time,
	underlyingSecurityID int32,
) ([]*model.InstrumentMaster, error) {
	// Build the base query with filters
	whereClause := []string{
		"(segment = 'DERIVATIVE' OR segment = 'COMMODITY')",
		"instrument IN ('STOCK_OPTIONS', 'INDEX_OPTIONS', 'COMMODITY_OPTIONS')",
	}
	args := []interface{}{}

	// Add exchange filter
	if exchange != "" {
		whereClause = append(whereClause, "exchange = ?")
		args = append(args, exchange)
	}

	// Add option type filter (CE/PE)
	if optionType != "" {
		whereClause = append(whereClause, "option_type = ?")
		args = append(args, optionType)
	}

	// Add strike price filter
	if strikePrice > 0 {
		whereClause = append(whereClause, "strike_price = ?")
		args = append(args, strikePrice)
	}

	// Add underlying symbol filter
	if symbol != "" {
		whereClause = append(whereClause, "underlying_symbol = ?")
		args = append(args, symbol)
	}

	// Add underlying security ID filter
	if underlyingSecurityID > 0 {
		whereClause = append(whereClause, "underlying_security_id = ?")
		args = append(args, underlyingSecurityID)
	}

	// Add expiry date filter (exact match)
	if !expiryDate.IsZero() {
		whereClause = append(whereClause, "DATE(expiry) = DATE(?)")
		args = append(args, expiryDate)
	}

	// Build WHERE clause
	whereStr := strings.Join(whereClause, " AND ")

	// Execute query
	query := fmt.Sprintf(`
		SELECT 
			exchange_security_id,
			exchange,
			segment,
			symbol,
			instrument,
			instrument_name,
			isin,
			underlying_security_id,
			underlying_symbol,
			expiry,
			expiry_flag,
			lot_units,
			tick_size,
			strike_price,
			option_type,
			instrument_type,
			bracket_flag,
			cover_flag,
			asm_gsm_flag,
			created_at,
			updated_at
		FROM sutram.instrument_master
		WHERE %s
		ORDER BY strike_price, expiry
	`, whereStr)

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying option contracts: %w", err)
	}

	// Check if exactly one record is found
	if len(instrumentMasters) == 0 {
		return nil, fmt.Errorf("no option contract found matching the specified criteria")
	}

	if len(instrumentMasters) > 1 {
		return nil, fmt.Errorf("multiple option contracts found (%d records), expected exactly one", len(instrumentMasters))
	}

	return instrumentMasters, nil
}

// FindFuturesContracts finds futures contracts based on filters
func (r *InstrumentMasterRepository) FindFuturesContracts(
	ctx context.Context,
	exchange string,
	expiryDate time.Time,
	underlyingSecurityID int32,
) ([]*model.InstrumentMaster, error) {
	// Build the base query with filters
	whereClause := []string{
		"(segment = 'DERIVATIVE' OR segment = 'COMMODITY')",
		"instrument IN ('STOCK_FUTURES', 'INDEX_FUTURES', 'COMMODITY_FUTURES')",
	}
	args := []interface{}{}

	// Add exchange filter
	if exchange != "" {
		whereClause = append(whereClause, "exchange = ?")
		args = append(args, exchange)
	}

	// Add underlying security ID filter
	if underlyingSecurityID > 0 {
		whereClause = append(whereClause, "underlying_security_id = ?")
		args = append(args, underlyingSecurityID)
	}

	// Add expiry date filter (exact match)
	if !expiryDate.IsZero() {
		whereClause = append(whereClause, "DATE(expiry) = DATE(?)")
		args = append(args, expiryDate)
	}

	// Build WHERE clause
	whereStr := strings.Join(whereClause, " AND ")

	// Execute query
	query := fmt.Sprintf(`
		SELECT 
			exchange_security_id,
			exchange,
			segment,
			symbol,
			instrument,
			instrument_name,
			isin,
			underlying_security_id,
			underlying_symbol,
			expiry,
			expiry_flag,
			lot_units,
			tick_size,
			strike_price,
			option_type,
			instrument_type,
			bracket_flag,
			cover_flag,
			asm_gsm_flag,
			created_at,
			updated_at
		FROM sutram.instrument_master
		WHERE %s
		ORDER BY expiry, symbol
	`, whereStr)

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying futures contracts: %w", err)
	}

	// Check if exactly one record is found
	if len(instrumentMasters) == 0 {
		return nil, fmt.Errorf("no futures contract found matching the specified criteria")
	}

	if len(instrumentMasters) > 1 {
		return nil, fmt.Errorf("multiple futures contracts found (%d records), expected exactly one", len(instrumentMasters))
	}

	return instrumentMasters, nil
}
