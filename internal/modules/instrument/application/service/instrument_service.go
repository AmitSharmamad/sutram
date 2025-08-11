package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"
)

// Common errors
var (
	ErrInstrumentNotFound = errors.New("instrument not found")
	ErrInvalidTimeRange   = errors.New("invalid time range")
)

// InstrumentService handles instrument-related business logic
type InstrumentService struct {
	instrumentRepo       repository.InstrumentRepository
	instrumentMasterRepo repository.InstrumentMasterRepository
}

// NewInstrumentService creates a new instrument service
func NewInstrumentService(
	instrumentRepo repository.InstrumentRepository,
	instrumentMasterRepo repository.InstrumentMasterRepository,
) *InstrumentService {
	return &InstrumentService{
		instrumentRepo:       instrumentRepo,
		instrumentMasterRepo: instrumentMasterRepo,
	}
}

// GetInstrumentsByExchangeSecurityID gets instruments by exchange security ID
func (s *InstrumentService) GetInstrumentsByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) ([]*model.Instrument, error) {
	instruments, err := s.instrumentRepo.FindByExchangeSecurityID(ctx, exchangeSecurityID)
	if err != nil {
		return nil, fmt.Errorf("error finding instruments: %w", err)
	}
	if len(instruments) == 0 {
		return nil, ErrInstrumentNotFound
	}
	return instruments, nil
}

// GetInstrumentsByExchangeSecurityIDPaginated gets paginated instruments by exchange security ID
func (s *InstrumentService) GetInstrumentsByExchangeSecurityIDPaginated(ctx context.Context, exchangeSecurityID int32, page, pageSize int) ([]*model.Instrument, int64, error) {
	// Validate and normalize pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	instruments, totalCount, err := s.instrumentRepo.FindByExchangeSecurityIDPaginated(ctx, exchangeSecurityID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding paginated instruments: %w", err)
	}

	if len(instruments) == 0 && page == 1 {
		return nil, 0, ErrInstrumentNotFound
	}

	return instruments, totalCount, nil
}

// GetInstrumentMasterByExchangeSecurityID gets an instrument master by exchange security ID
func (s *InstrumentService) GetInstrumentMasterByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) (*model.InstrumentMaster, error) {
	instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, exchangeSecurityID)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument master by exchange security ID: %w", err)
	}
	if instrumentMaster == nil {
		return nil, ErrInstrumentNotFound
	}
	return instrumentMaster, nil
}

// GetInstrumentsByTimeRange gets instruments by exchange security ID within a time range
func (s *InstrumentService) GetInstrumentsByTimeRange(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string) ([]*model.Instrument, error) {
	// Validate time range
	if start.After(end) {
		return nil, ErrInvalidTimeRange
	}

	// Default to '1m' if timeFrame is empty
	if timeFrame == "" {
		timeFrame = "1m"
	}

	instruments, err := s.instrumentRepo.FindByExchangeSecurityIDAndTimeRange(ctx, exchangeSecurityID, start, end, timeFrame)
	if err != nil {
		return nil, fmt.Errorf("error finding instruments by time range: %w", err)
	}
	return instruments, nil
}

// GetInstrumentMasterBySymbol gets an instrument master by symbol
func (s *InstrumentService) GetInstrumentMasterBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error) {
	instrumentMaster, err := s.instrumentMasterRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument master by symbol: %w", err)
	}
	if instrumentMaster == nil {
		return nil, ErrInstrumentNotFound
	}
	return instrumentMaster, nil
}

// GetAllInstrumentMasters gets all instrument masters
func (s *InstrumentService) GetAllInstrumentMasters(ctx context.Context) ([]*model.InstrumentMaster, error) {
	instrumentMasters, err := s.instrumentMasterRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding all instrument masters: %w", err)
	}
	return instrumentMasters, nil
}

// GetInstrumentMastersWithFilters gets instrument masters with filtering and pagination
func (s *InstrumentService) GetInstrumentMastersWithFilters(ctx context.Context, exchange string, segment string, query string, page, pageSize int) ([]*model.InstrumentMaster, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	instrumentMasters, totalCount, err := s.instrumentMasterRepo.FindWithFilters(ctx, exchange, segment, query, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding instrument masters with filters: %w", err)
	}

	return instrumentMasters, totalCount, nil
}

// CreateInstrument creates a new instrument
func (s *InstrumentService) CreateInstrument(ctx context.Context, instrument *model.Instrument) error {
	now := time.Now()
	instrument.CreatedAt = now
	instrument.UpdatedAt = now

	if err := s.instrumentRepo.Create(ctx, instrument); err != nil {
		return fmt.Errorf("error creating instrument: %w", err)
	}
	return nil
}

// UpdateInstrument updates an instrument
func (s *InstrumentService) UpdateInstrument(ctx context.Context, instrument *model.Instrument) error {
	// Check if instrument exists
	existingInstruments, err := s.instrumentRepo.FindByExchangeSecurityID(ctx, instrument.ExchangeSecurityID)
	if err != nil {
		return fmt.Errorf("error finding instrument: %w", err)
	}
	if len(existingInstruments) == 0 {
		return ErrInstrumentNotFound
	}

	// Find the specific instrument with matching timestamp
	var existing *model.Instrument
	for _, inst := range existingInstruments {
		if inst.Timestamp.Equal(instrument.Timestamp) {
			existing = inst
			break
		}
	}

	if existing == nil {
		return ErrInstrumentNotFound
	}

	// Update timestamp
	instrument.UpdatedAt = time.Now()
	instrument.CreatedAt = existing.CreatedAt

	if err := s.instrumentRepo.Update(ctx, instrument); err != nil {
		return fmt.Errorf("error updating instrument: %w", err)
	}
	return nil
}

// CreateInstrumentMaster creates a new instrument master
func (s *InstrumentService) CreateInstrumentMaster(ctx context.Context, master *model.InstrumentMaster) error {
	// Remove segment validation since it's now a string
	// Remove instrument type validation as well since the field structure changed

	if err := s.instrumentMasterRepo.Create(ctx, master); err != nil {
		return fmt.Errorf("error creating instrument master: %w", err)
	}
	return nil
}

// UpdateInstrumentMaster updates an instrument master
func (s *InstrumentService) UpdateInstrumentMaster(ctx context.Context, master *model.InstrumentMaster) error {
	// Validate enums before update
	// Remove segment validation since it's now a string
	// if !master.ExchInstrumentType.IsValid() {
	// 	return fmt.Errorf("invalid instrument type: %s", master.ExchInstrumentType)
	// }

	// Check if instrument master exists by exchange security ID
	existing, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, master.ExchangeSecurityID)
	if err != nil {
		return fmt.Errorf("error finding instrument master: %w", err)
	}
	if existing == nil {
		return ErrInstrumentNotFound
	}

	if err := s.instrumentMasterRepo.Update(ctx, master); err != nil {
		return fmt.Errorf("error updating instrument master: %w", err)
	}
	return nil
}

// DeleteInstrument deletes instruments by exchange security ID
func (s *InstrumentService) DeleteInstrument(ctx context.Context, exchangeSecurityID int32) error {
	// Check if instruments exist
	existing, err := s.instrumentRepo.FindByExchangeSecurityID(ctx, exchangeSecurityID)
	if err != nil {
		return fmt.Errorf("error finding instruments: %w", err)
	}
	if len(existing) == 0 {
		return ErrInstrumentNotFound
	}

	if err := s.instrumentRepo.Delete(ctx, exchangeSecurityID); err != nil {
		return fmt.Errorf("error deleting instruments: %w", err)
	}
	return nil
}

// DeleteInstrumentMaster deletes an instrument master by exchange security ID
func (s *InstrumentService) DeleteInstrumentMaster(ctx context.Context, exchangeSecurityID int32) error {
	// Check if instrument master exists
	existing, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, exchangeSecurityID)
	if err != nil {
		return fmt.Errorf("error finding instrument master: %w", err)
	}
	if existing == nil {
		return ErrInstrumentNotFound
	}

	if err := s.instrumentMasterRepo.Delete(ctx, exchangeSecurityID); err != nil {
		return fmt.Errorf("error deleting instrument master: %w", err)
	}
	return nil
}

// GetInstrumentDataWithLookback gets instrument data for a given time range plus historical lookback data
// Combines both the requested date range data and lookback data into a single result
func (s *InstrumentService) GetInstrumentDataWithLookback(
	ctx context.Context,
	exchangeSecurityID int32,
	start, end time.Time,
	timeFrame string,
	lookBack int,
) (*model.StockData, error) {
	// Validate time range
	if start.After(end) {
		return nil, ErrInvalidTimeRange
	}

	// Get data for the requested date range
	rangeInstruments, err := s.instrumentRepo.FindByExchangeSecurityIDAndTimeRange(ctx, exchangeSecurityID, start, end, timeFrame)
	if err != nil {
		return nil, fmt.Errorf("error finding instruments for date range: %w", err)
	}

	// Convert pointers to values for StockData
	dateRangeData := make(model.StockData, len(rangeInstruments))
	for i, instrument := range rangeInstruments {
		dateRangeData[i] = *instrument
	}

	// Calculate lookback start time (this would depend on timeFrame implementation)
	// For now, we'll assume lookBack is in days
	lookBackStart := start.AddDate(0, 0, -lookBack)

	// Get lookback data
	lookBackInstruments, err := s.instrumentRepo.FindByExchangeSecurityIDAndTimeRange(ctx, exchangeSecurityID, lookBackStart, start, timeFrame)
	if err != nil {
		// If there's an error with lookback, just return the date range data
		return &dateRangeData, nil
	}

	// Convert pointers to values for StockData
	lookBackData := make(model.StockData, len(lookBackInstruments))
	for i, instrument := range lookBackInstruments {
		lookBackData[i] = *instrument
	}

	// Reverse the lookback data (older to newer)
	for i, j := 0, len(lookBackData)-1; i < j; i, j = i+1, j-1 {
		lookBackData[i], lookBackData[j] = lookBackData[j], lookBackData[i]
	}

	// Combine lookback and date range data
	result := append(make(model.StockData, 0), lookBackData...)
	result = append(result, dateRangeData...)

	return &result, nil
}

// GetInstrumentsByExchangeSecurityIDWithLookback gets the exact number of historical data points before the given start time
func (s *InstrumentService) GetInstrumentsByExchangeSecurityIDWithLookback(
	ctx context.Context,
	exchangeSecurityID int32,
	startTime time.Time,
	timeFrame string,
	lookback int,
) ([]*model.Instrument, error) {
	// Default to '1m' if timeFrame is empty
	if timeFrame == "" {
		timeFrame = "1m"
	}

	instruments, err := s.instrumentRepo.FindByExchangeSecurityIDWithLookback(ctx, exchangeSecurityID, startTime, timeFrame, lookback)
	if err != nil {
		return nil, fmt.Errorf("error finding instruments with lookback: %w", err)
	}
	return instruments, nil
}

// GetInstrumentMastersByExchangeSecurityIDs gets multiple instrument masters by their exchange security IDs
func (s *InstrumentService) GetInstrumentMastersByExchangeSecurityIDs(ctx context.Context, exchangeSecurityIDs []int32) ([]*model.InstrumentMaster, error) {
	instrumentMasters, err := s.instrumentMasterRepo.FindByExchangeSecurityIDs(ctx, exchangeSecurityIDs)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument masters by exchange security IDs: %w", err)
	}
	return instrumentMasters, nil
}

// ResolveOptionsContractsFromMaster resolves options contracts from instrument_master using provided filters
func (s *InstrumentService) ResolveOptionsContractsFromMaster(
	ctx context.Context,
	strikePrice float64,
	optionType string,
	expiry time.Time,
	parentSecurityID *int32,
) (*model.InstrumentMaster, error) {

	// NOTE: below code is dummy, remove it once the development is done
	instrumentExchangeDummy, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, 71950)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument exchange dummy: %w", err)
	}

	if instrumentExchangeDummy != nil {
		return instrumentExchangeDummy, nil
	}

	// NOTE: Actual code starts from here
	parentSecurityIDValue := int32(*parentSecurityID)
	instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, parentSecurityIDValue)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument master by exchange security ID: %w", err)
	}
	instrumentMasterSymbol := instrumentMaster.Symbol
	instrumentExchange := instrumentMaster.Exchange

	instruments, err := s.instrumentMasterRepo.FindOptionsContracts(
		ctx, strikePrice, optionType, instrumentExchange, instrumentMasterSymbol, expiry, *parentSecurityID,
	)
	if err != nil {
		return nil, fmt.Errorf("error resolving options contracts from master: %w", err)
	}
	if len(instruments) == 1 {
		return instruments[0], nil
	} else {
		return nil, ErrInstrumentNotFound
	}
}

// ResolveFuturesContractsFromMaster resolves futures contracts from instrument_master using provided filters
func (s *InstrumentService) ResolveFuturesContractsFromMaster(
	ctx context.Context,
	exchange string,
	expiry time.Time,
	parentSecurityID *int32,
) (*model.InstrumentMaster, error) {

	// NOTE: Actual code starts from here
	instruments, err := s.instrumentMasterRepo.FindFuturesContracts(
		ctx, exchange, expiry, *parentSecurityID,
	)
	if err != nil {
		return nil, fmt.Errorf("error resolving futures contracts from master: %w", err)
	}
	if len(instruments) == 1 {
		return instruments[0], nil
	} else {
		return nil, ErrInstrumentNotFound
	}
}
