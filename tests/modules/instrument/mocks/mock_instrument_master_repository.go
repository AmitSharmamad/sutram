package mocks

import (
	"context"

	"algo/internal/modules/instrument/domain/model"
)

// MockInstrumentMasterRepository is a mock implementation of the InstrumentMasterRepository interface
type MockInstrumentMasterRepository struct {
	FindByIDFunc              func(ctx context.Context, id uint64) (*model.InstrumentMaster, error)
	FindBySymbolFunc          func(ctx context.Context, symbol string) (*model.InstrumentMaster, error)
	FindByExchangeFunc        func(ctx context.Context, exchange string) ([]*model.InstrumentMaster, error)
	FindBySegmentFunc         func(ctx context.Context, segment string) ([]*model.InstrumentMaster, error)
	FindAllFunc               func(ctx context.Context) ([]*model.InstrumentMaster, error)
	FindWithFiltersFunc       func(ctx context.Context, exchange, segment string, offset, limit int) ([]*model.InstrumentMaster, int, error)
	CreateFunc                func(ctx context.Context, instrumentMaster *model.InstrumentMaster) error
	UpdateFunc                func(ctx context.Context, instrumentMaster *model.InstrumentMaster) error
	DeleteFunc                func(ctx context.Context, id uint64) error
	FindByExchangeSecurityIDFunc func(ctx context.Context, exchangeSecurityID int32) (*model.InstrumentMaster, error)
	FindByIDsFunc             func(ctx context.Context, ids []uint64) ([]*model.InstrumentMaster, error)
}

// FindByID implements the repository interface
func (m *MockInstrumentMasterRepository) FindByID(ctx context.Context, id uint64) (*model.InstrumentMaster, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

// FindBySymbol implements the repository interface
func (m *MockInstrumentMasterRepository) FindBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error) {
	if m.FindBySymbolFunc != nil {
		return m.FindBySymbolFunc(ctx, symbol)
	}
	return nil, nil
}

// FindByExchange implements the repository interface
func (m *MockInstrumentMasterRepository) FindByExchange(ctx context.Context, exchange string) ([]*model.InstrumentMaster, error) {
	if m.FindByExchangeFunc != nil {
		return m.FindByExchangeFunc(ctx, exchange)
	}
	return nil, nil
}

// FindBySegment implements the repository interface
func (m *MockInstrumentMasterRepository) FindBySegment(ctx context.Context, segment string) ([]*model.InstrumentMaster, error) {
	if m.FindBySegmentFunc != nil {
		return m.FindBySegmentFunc(ctx, segment)
	}
	return nil, nil
}

// FindAll implements the repository interface
func (m *MockInstrumentMasterRepository) FindAll(ctx context.Context) ([]*model.InstrumentMaster, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

// FindWithFilters implements the repository interface
func (m *MockInstrumentMasterRepository) FindWithFilters(ctx context.Context, exchange, segment string, offset, limit int) ([]*model.InstrumentMaster, int, error) {
	if m.FindWithFiltersFunc != nil {
		return m.FindWithFiltersFunc(ctx, exchange, segment, offset, limit)
	}
	return nil, 0, nil
}

// Create implements the repository interface
func (m *MockInstrumentMasterRepository) Create(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, instrumentMaster)
	}
	return nil
}

// Update implements the repository interface
func (m *MockInstrumentMasterRepository) Update(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, instrumentMaster)
	}
	return nil
}

// Delete implements the repository interface
func (m *MockInstrumentMasterRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// FindByExchangeSecurityID implements the repository interface
func (m *MockInstrumentMasterRepository) FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) (*model.InstrumentMaster, error) {
	if m.FindByExchangeSecurityIDFunc != nil {
		return m.FindByExchangeSecurityIDFunc(ctx, exchangeSecurityID)
	}
	return nil, nil
}

// FindByIDs implements the repository interface
func (m *MockInstrumentMasterRepository) FindByIDs(ctx context.Context, ids []uint64) ([]*model.InstrumentMaster, error) {
	if m.FindByIDsFunc != nil {
		return m.FindByIDsFunc(ctx, ids)
	}
	return nil, nil
}