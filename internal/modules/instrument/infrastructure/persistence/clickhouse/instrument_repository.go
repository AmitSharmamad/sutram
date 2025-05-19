package clickhouse

import (
	"context"
	"fmt"
	"time"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"

	"github.com/jmoiron/sqlx"
)

// InstrumentRepository is a ClickHouse implementation of the InstrumentRepository interface
type InstrumentRepository struct {
	db *sqlx.DB
}

// NewInstrumentRepository creates a new ClickHouse instrument repository
func NewInstrumentRepository(db *sqlx.DB) repository.InstrumentRepository {
	return &InstrumentRepository{
		db: db,
	}
}

// FindByID finds an instrument by ID
func (r *InstrumentRepository) FindByID(ctx context.Context, id uint64) (*model.Instrument, error) {
	query := `
		SELECT 
			id, 
			instrument, 
			open, 
			high, 
			low, 
			close, 
			volume, 
			timestamp, 
			created_at, 
			updated_at
		FROM sutram.instrument
		WHERE id = ?
		LIMIT 1
	`

	instrument := &model.Instrument{}
	err := r.db.GetContext(ctx, instrument, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument by ID: %w", err)
	}

	return instrument, nil
}

// FindBySymbol finds instruments by symbol
func (r *InstrumentRepository) FindBySymbol(ctx context.Context, symbol string) ([]*model.Instrument, error) {
	query := `
		SELECT 
			id, 
			instrument, 
			open, 
			high, 
			low, 
			close, 
			volume, 
			timestamp, 
			created_at, 
			updated_at
		FROM sutram.instrument
		WHERE instrument = ?
		ORDER BY timestamp DESC
		LIMIT 1000
	`

	var instruments []*model.Instrument
	err := r.db.SelectContext(ctx, &instruments, query, symbol)
	if err != nil {
		return nil, fmt.Errorf("error querying instruments by symbol: %w", err)
	}

	return instruments, nil
}

// FindByIDAndTimeRange finds instruments by ID within a time range
// timeFrame parameter allows data aggregation by different time intervals ('1s', '1m', '3m', '5m', '15m', '30m', '1h', '3h', '1d', '1M')
func (r *InstrumentRepository) FindByIDAndTimeRange(ctx context.Context, id uint64, start, end time.Time, timeFrame string) ([]*model.Instrument, error) {
	query := `
		SELECT
			id,
			instrument,
			CASE
				WHEN ? = '1s' THEN toTimeZone(timestamp, 'Asia/Kolkata')
				WHEN ? = '1m' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 1 MINUTE)
				WHEN ? = '3m' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 3 MINUTE)
				WHEN ? = '5m' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 5 MINUTE)
				WHEN ? = '15m' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 15 MINUTE)
				WHEN ? = '30m' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 30 MINUTE)
				WHEN ? = '1h' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 1 HOUR)
				WHEN ? = '3h' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 3 HOUR)
				WHEN ? = '1d' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 1 DAY)
				WHEN ? = '1M' THEN toStartOfInterval(toTimeZone(timestamp, 'Asia/Kolkata'), INTERVAL 1 MONTH)
			END AS timestamp,
			MIN(low) AS low,
			MAX(high) AS high,
			first_value(open) AS open,
			last_value(close) AS close,
			SUM(volume) AS volume,
			MIN(created_at) AS created_at,
			MAX(updated_at) AS updated_at
		FROM sutram.instrument
		WHERE 
			id = ?
		GROUP BY
			id,
			instrument,
			timestamp
		ORDER BY
			timestamp ASC
	`

	var instruments []*model.Instrument
	err := r.db.SelectContext(ctx, &instruments, query,
		timeFrame, timeFrame, timeFrame, timeFrame, timeFrame,
		timeFrame, timeFrame, timeFrame, timeFrame, timeFrame,
		id)

	if err != nil {
		return nil, fmt.Errorf("error querying instruments by time range: %w", err)
	}

	return instruments, nil
}

// Create creates a new instrument
func (r *InstrumentRepository) Create(ctx context.Context, instrument *model.Instrument) error {
	query := `
		INSERT INTO sutram.instrument (
			id,
			instrument,
			open,
			high,
			low,
			close,
			volume,
			timestamp,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrument.ID,
		instrument.Symbol,
		instrument.Open,
		instrument.High,
		instrument.Low,
		instrument.Close,
		instrument.Volume,
		instrument.Timestamp,
		instrument.CreatedAt,
		instrument.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating instrument: %w", err)
	}

	return nil
}

// Update updates an instrument
func (r *InstrumentRepository) Update(ctx context.Context, instrument *model.Instrument) error {
	// Note: In ClickHouse, we don't really update rows, we insert new rows with updated values
	// This is a simplified version assuming you have a way to handle this in your application
	query := `
		ALTER TABLE instrument UPDATE
			open = ?,
			high = ?,
			low = ?,
			close = ?,
			volume = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrument.Open,
		instrument.High,
		instrument.Low,
		instrument.Close,
		instrument.Volume,
		instrument.UpdatedAt,
		instrument.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating instrument: %w", err)
	}

	return nil
}

// Delete deletes an instrument
func (r *InstrumentRepository) Delete(ctx context.Context, id uint64) error {
	query := `
		ALTER TABLE instrument DELETE WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting instrument: %w", err)
	}

	return nil
}
