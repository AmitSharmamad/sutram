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

	// FindByIDAndTimeRange finds instruments by ID within a time range and aggregated by timeFrame
	// timeFrame parameter allows data aggregation by different time intervals ('1s', '1m', '3m', '5m', '15m', '30m', '1h', '3h', '1d', '1M')
	FindByIDAndTimeRange(ctx context.Context, id uint64, start, end time.Time, timeFrame string) ([]*model.Instrument, error)

	// Create creates a new instrument
	Create(ctx context.Context, instrument *model.Instrument) error

	// Update updates an instrument
	Update(ctx context.Context, instrument *model.Instrument) error

	// Delete deletes an instrument
	Delete(ctx context.Context, id uint64) error
}
