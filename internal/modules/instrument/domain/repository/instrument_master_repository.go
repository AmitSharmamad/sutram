package repository

import (
	"context"

	"algo/internal/modules/instrument/domain/model"
)

// InstrumentMasterRepository provides access to instrument master data
type InstrumentMasterRepository interface {
	// FindByID finds an instrument master by ID
	FindByID(ctx context.Context, id uint64) (*model.InstrumentMaster, error)

	// FindBySymbol finds an instrument master by symbol
	FindBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error)

	// FindByExchange finds instrument masters by exchange
	FindByExchange(ctx context.Context, exchange string) ([]*model.InstrumentMaster, error)

	// FindBySegment finds instrument masters by segment
	FindBySegment(ctx context.Context, segment string) ([]*model.InstrumentMaster, error)

	// FindAll returns all instrument masters
	FindAll(ctx context.Context) ([]*model.InstrumentMaster, error)

	// Create creates a new instrument master
	Create(ctx context.Context, instrumentMaster *model.InstrumentMaster) error

	// Update updates an instrument master
	Update(ctx context.Context, instrumentMaster *model.InstrumentMaster) error

	// Delete deletes an instrument master
	Delete(ctx context.Context, id uint64) error
}
