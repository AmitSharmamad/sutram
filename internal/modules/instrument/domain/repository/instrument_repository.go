package repository

import (
	"context"
	"time"

	"algo/internal/modules/instrument/domain/model"
)

// InstrumentRepository provides access to instrument data
type InstrumentRepository interface {
	// FindByID finds an instrument by ID
	FindByID(ctx context.Context, id uint64) (*model.Instrument, error)

	// FindBySymbol finds instruments by symbol
	FindBySymbol(ctx context.Context, symbol string) ([]*model.Instrument, error)

	// FindBySymbolAndTimeRange finds instruments by symbol within a time range
	FindBySymbolAndTimeRange(ctx context.Context, symbol string, start, end time.Time) ([]*model.Instrument, error)

	// Create creates a new instrument
	Create(ctx context.Context, instrument *model.Instrument) error

	// Update updates an instrument
	Update(ctx context.Context, instrument *model.Instrument) error

	// Delete deletes an instrument
	Delete(ctx context.Context, id uint64) error
}
