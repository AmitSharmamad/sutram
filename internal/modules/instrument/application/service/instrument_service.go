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

// GetInstrumentByID gets an instrument by ID
func (s *InstrumentService) GetInstrumentByID(ctx context.Context, id uint64) (*model.Instrument, error) {
	instrument, err := s.instrumentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument: %w", err)
	}
	if instrument == nil {
		return nil, ErrInstrumentNotFound
	}
	return instrument, nil
}

// GetInstrumentMasterByID gets an instrument master by ID
func (s *InstrumentService) GetInstrumentMasterByID(ctx context.Context, id uint64) (*model.InstrumentMaster, error) {
	instrumentMaster, err := s.instrumentMasterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding instrument master: %w", err)
	}
	if instrumentMaster == nil {
		return nil, ErrInstrumentNotFound
	}
	return instrumentMaster, nil
}

// GetInstrumentsBySymbol gets instruments by symbol
func (s *InstrumentService) GetInstrumentsBySymbol(ctx context.Context, symbol string) ([]*model.Instrument, error) {
	instruments, err := s.instrumentRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("error finding instruments by symbol: %w", err)
	}
	return instruments, nil
}

// GetInstrumentsByTimeRange gets instruments by ID within a time range
func (s *InstrumentService) GetInstrumentsByTimeRange(ctx context.Context, id uint64, start, end time.Time, timeFrame string) ([]*model.Instrument, error) {
	// Validate time range
	if start.After(end) {
		return nil, ErrInvalidTimeRange
	}

	// Default to '1m' if timeFrame is empty
	if timeFrame == "" {
		timeFrame = "1m"
	}

	instruments, err := s.instrumentRepo.FindByIDAndTimeRange(ctx, id, start, end, timeFrame)
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
func (s *InstrumentService) GetInstrumentMastersWithFilters(ctx context.Context, exchange, segment string, page, pageSize int) ([]*model.InstrumentMaster, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	instrumentMasters, totalCount, err := s.instrumentMasterRepo.FindWithFilters(ctx, exchange, segment, offset, limit)
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
	existing, err := s.instrumentRepo.FindByID(ctx, instrument.ID)
	if err != nil {
		return fmt.Errorf("error finding instrument: %w", err)
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
	if err := s.instrumentMasterRepo.Create(ctx, master); err != nil {
		return fmt.Errorf("error creating instrument master: %w", err)
	}
	return nil
}

// UpdateInstrumentMaster updates an instrument master
func (s *InstrumentService) UpdateInstrumentMaster(ctx context.Context, master *model.InstrumentMaster) error {
	// Check if instrument master exists
	existing, err := s.instrumentMasterRepo.FindByID(ctx, master.ID)
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

// DeleteInstrument deletes an instrument
func (s *InstrumentService) DeleteInstrument(ctx context.Context, id uint64) error {
	// Check if instrument exists
	existing, err := s.instrumentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding instrument: %w", err)
	}
	if existing == nil {
		return ErrInstrumentNotFound
	}

	if err := s.instrumentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting instrument: %w", err)
	}
	return nil
}

// DeleteInstrumentMaster deletes an instrument master
func (s *InstrumentService) DeleteInstrumentMaster(ctx context.Context, id uint64) error {
	// Check if instrument master exists
	existing, err := s.instrumentMasterRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding instrument master: %w", err)
	}
	if existing == nil {
		return ErrInstrumentNotFound
	}

	if err := s.instrumentMasterRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting instrument master: %w", err)
	}
	return nil
}

// GetInstrumentDataWithLookback gets instrument data for a given time range plus historical lookback data
// Combines both the requested date range data and lookback data into a single result
func (s *InstrumentService) GetInstrumentDataWithLookback(
	ctx context.Context,
	id uint64,
	start, end time.Time,
	timeFrame string,
	lookBack int,
) (*model.StockData, error) {
	// Validate time range
	if start.After(end) {
		return nil, ErrInvalidTimeRange
	}

	// Get data for the requested date range
	rangeInstruments, err := s.instrumentRepo.FindByIDAndTimeRange(ctx, id, start, end, timeFrame)
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
	lookBackInstruments, err := s.instrumentRepo.FindByIDAndTimeRange(ctx, id, lookBackStart, start, timeFrame)
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
