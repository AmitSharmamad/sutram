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

	// FindByExchangeSecurityIDPaginated finds instruments by exchange security ID with pagination
	FindByExchangeSecurityIDPaginated(ctx context.Context, exchangeSecurityID int32, page, pageSize int) ([]*model.Instrument, int64, error)

	// FindByExchangeSecurityIDAndTimeRange finds instruments by exchange security ID within a time range and aggregated by timeFrame
	// timeFrame parameter allows data aggregation by different time intervals ('1s', '1m', '3m', '5m', '15m', '30m', '1h', '3h', '1d', '1M')
	FindByExchangeSecurityIDAndTimeRange(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string) ([]*model.Instrument, error)

	// FindByExchangeSecurityIDAndTimePoints finds instruments by exchange security ID for specific time points within a date range
	// timePoints parameter accepts a list of specific times (e.g., ["09:30:00", "15:00:00", "12:00:00"])
	// timeFrame parameter allows data aggregation by different time intervals ('1s', '1m', '3m', '5m', '15m', '30m', '1h', '3h', '1d', '1M')
	FindByExchangeSecurityIDAndTimePoints(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string, timePoints []string) ([]*model.Instrument, error)

	// FindByExchangeSecurityIDWithLookback gets the exact number of historical data points before the given start time
	// Returns the data in descending order (newest first) limited by lookback count
	FindByExchangeSecurityIDWithLookback(ctx context.Context, exchangeSecurityID int32, startTime time.Time, timeFrame string, lookback int) ([]*model.Instrument, error)

	// Create creates a new instrument
	Create(ctx context.Context, instrument *model.Instrument) error

	// Update updates an instrument
	Update(ctx context.Context, instrument *model.Instrument) error

	// Delete deletes instruments by exchange security ID
	Delete(ctx context.Context, exchangeSecurityID int32) error
}
