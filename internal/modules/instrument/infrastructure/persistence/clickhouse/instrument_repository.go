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

// FindByExchangeSecurityID finds instruments by exchange security ID
func (r *InstrumentRepository) FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) ([]*model.Instrument, error) {
	query := `
		SELECT 
			exchange_security_id, 
			open, 
			high, 
			low, 
			close, 
			volume, 
			timestamp, 
			created_at, 
			updated_at
		FROM sutram.instrument
		WHERE exchange_security_id = ?
		ORDER BY timestamp DESC
	`

	var instruments []*model.Instrument
	err := r.db.SelectContext(ctx, &instruments, query, exchangeSecurityID)
	if err != nil {
		return nil, fmt.Errorf("error querying instruments by exchange security ID: %w", err)
	}

	return instruments, nil
}

// FindByExchangeSecurityIDAndTimeRange finds instruments by exchange security ID within a time range
// timeFrame parameter allows data aggregation by different time intervals ('1s', '1m', '3m', '5m', '15m', '30m', '1h', '3h', '1d', '1M')
func (r *InstrumentRepository) FindByExchangeSecurityIDAndTimeRange(ctx context.Context, exchangeSecurityID int32, start, end time.Time, timeFrame string) ([]*model.Instrument, error) {
	query := `
		SELECT
			exchange_security_id,
			CASE
				WHEN ? = '1s' THEN toTimeZone(timestamp, 'Asia/Kolkata')
				WHEN ? = '1m' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 1 MINUTE
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
				WHEN ? = '3m' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 3 MINUTE
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
				WHEN ? = '5m' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 5 MINUTE
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
				WHEN ? = '15m' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 15 MINUTE
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
				WHEN ? = '30m' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 30 MINUTE
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
				WHEN ? = '1h' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 1 HOUR
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
				WHEN ? = '3h' THEN toStartOfInterval(
					toTimeZone(timestamp, 'Asia/Kolkata') - INTERVAL 9 HOUR - INTERVAL 15 MINUTE,
					INTERVAL 3 HOUR
				) + INTERVAL 9 HOUR + INTERVAL 15 MINUTE
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
			exchange_security_id = ? AND
			timestamp BETWEEN ? AND ?
		GROUP BY
			exchange_security_id,
			timestamp
		ORDER BY
			timestamp ASC
	`

	var instruments []*model.Instrument
	err := r.db.SelectContext(ctx, &instruments, query,
		timeFrame, timeFrame, timeFrame, timeFrame, timeFrame,
		timeFrame, timeFrame, timeFrame, timeFrame, timeFrame,
		exchangeSecurityID, start, end)

	if err != nil {
		return nil, fmt.Errorf("error querying instruments by exchange security ID and time range: %w", err)
	}

	return instruments, nil
}

// Create creates a new instrument
func (r *InstrumentRepository) Create(ctx context.Context, instrument *model.Instrument) error {
	query := `
		INSERT INTO sutram.instrument (
			exchange_security_id,
			open,
			high,
			low,
			close,
			volume,
			timestamp,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrument.ExchangeSecurityID,
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
		WHERE exchange_security_id = ? AND timestamp = ?
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
		instrument.ExchangeSecurityID,
		instrument.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("error updating instrument: %w", err)
	}

	return nil
}

// Delete deletes instruments by exchange security ID
func (r *InstrumentRepository) Delete(ctx context.Context, exchangeSecurityID int32) error {
	query := `
		ALTER TABLE instrument DELETE WHERE exchange_security_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, exchangeSecurityID)
	if err != nil {
		return fmt.Errorf("error deleting instruments: %w", err)
	}

	return nil
}

// FindByExchangeSecurityIDPaginated finds instruments by exchange security ID with pagination
func (r *InstrumentRepository) FindByExchangeSecurityIDPaginated(ctx context.Context, exchangeSecurityID int32, page, pageSize int) ([]*model.Instrument, int64, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Ensure pageSize doesn't exceed maximum
	if pageSize > 1000 {
		pageSize = 1000
	}

	// First, get the total count
	countQuery := `
		SELECT COUNT(*) 
		FROM sutram.instrument
		WHERE exchange_security_id = ?
	`
	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, exchangeSecurityID)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting instruments: %w", err)
	}

	// If no records, return early
	if totalCount == 0 {
		return []*model.Instrument{}, 0, nil
	}

	// Get paginated data
	query := `
		SELECT 
			exchange_security_id, 
			open, 
			high, 
			low, 
			close, 
			volume, 
			timestamp, 
			created_at, 
			updated_at
		FROM sutram.instrument
		WHERE exchange_security_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	var instruments []*model.Instrument
	err = r.db.SelectContext(ctx, &instruments, query, exchangeSecurityID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying paginated instruments: %w", err)
	}

	return instruments, totalCount, nil
}
