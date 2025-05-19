package repository

import (
	"context"
	"time"

	"algo/internal/modules/instrument/domain/model"
)

// InstrumentRepository provides access to instrument data
type InstrumentRepository interface {
	// FindByExchangeSecurityID finds instruments by exchange security ID
	FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) ([]*model.Instrument, error)

	// FindByExchangeSecurityIDAndTimeRange finds instruments by exchange security ID within a time range and aggregated by timeFrame
	// timeFrame parameter allows data aggregation by different time intervals ('1s', '1m', '3m', '5m', '15m', '30m', '1h', '3h', '1d', '1M')
	FindByExchangeSecurityIDAndTimeRange(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string) ([]*model.Instrument, error)

	// Create creates a new instrument
	Create(ctx context.Context, instrument *model.Instrument) error

	// Update updates an instrument
	Update(ctx context.Context, instrument *model.Instrument) error

	// Delete deletes instruments by exchange security ID
	Delete(ctx context.Context, exchangeSecurityID int32) error
}
