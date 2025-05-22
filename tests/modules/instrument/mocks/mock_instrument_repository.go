package mocks

import (
	"context"
	"time"

	"algo/internal/modules/instrument/domain/model"
)

// MockInstrumentRepository is a mock implementation of the InstrumentRepository interface
type MockInstrumentRepository struct {
	FindByExchangeSecurityIDFunc                func(ctx context.Context, exchangeSecurityID int32) ([]*model.Instrument, error)
	FindByExchangeSecurityIDAndTimeRangeFunc    func(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string) ([]*model.Instrument, error)
	CreateFunc                                 func(ctx context.Context, instrument *model.Instrument) error
	UpdateFunc                                 func(ctx context.Context, instrument *model.Instrument) error
	DeleteFunc                                 func(ctx context.Context, exchangeSecurityID int32) error
}

// FindByExchangeSecurityID implements the repository interface
func (m *MockInstrumentRepository) FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) ([]*model.Instrument, error) {
	if m.FindByExchangeSecurityIDFunc != nil {
		return m.FindByExchangeSecurityIDFunc(ctx, exchangeSecurityID)
	}
	return nil, nil
}

// FindByExchangeSecurityIDAndTimeRange implements the repository interface
func (m *MockInstrumentRepository) FindByExchangeSecurityIDAndTimeRange(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string) ([]*model.Instrument, error) {
	if m.FindByExchangeSecurityIDAndTimeRangeFunc != nil {
		return m.FindByExchangeSecurityIDAndTimeRangeFunc(ctx, exchangeSecurityID, start, end, timeFrame)
	}
	return nil, nil
}

// Create implements the repository interface
func (m *MockInstrumentRepository) Create(ctx context.Context, instrument *model.Instrument) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, instrument)
	}
	return nil
}

// Update implements the repository interface
func (m *MockInstrumentRepository) Update(ctx context.Context, instrument *model.Instrument) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, instrument)
	}
	return nil
}

// Delete implements the repository interface
func (m *MockInstrumentRepository) Delete(ctx context.Context, exchangeSecurityID int32) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, exchangeSecurityID)
	}
	return nil
}