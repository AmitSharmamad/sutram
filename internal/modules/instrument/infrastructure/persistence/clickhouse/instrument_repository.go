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
		FROM instrument
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
		FROM instrument
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

// FindBySymbolAndTimeRange finds instruments by symbol within a time range
func (r *InstrumentRepository) FindBySymbolAndTimeRange(ctx context.Context, symbol string, start, end time.Time) ([]*model.Instrument, error) {
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
		FROM instrument
		WHERE 
			instrument = ? AND
			timestamp >= ? AND
			timestamp <= ?
		ORDER BY timestamp ASC
	`

	var instruments []*model.Instrument
	err := r.db.SelectContext(ctx, &instruments, query, symbol, start, end)
	if err != nil {
		return nil, fmt.Errorf("error querying instruments by time range: %w", err)
	}

	return instruments, nil
}

// Create creates a new instrument
func (r *InstrumentRepository) Create(ctx context.Context, instrument *model.Instrument) error {
	query := `
		INSERT INTO instrument (
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
