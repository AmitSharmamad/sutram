package clickhouse

import (
	"context"
	"fmt"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"

	"github.com/jmoiron/sqlx"
)

// InstrumentMasterRepository is a ClickHouse implementation of the InstrumentMasterRepository interface
type InstrumentMasterRepository struct {
	db *sqlx.DB
}

// NewInstrumentMasterRepository creates a new ClickHouse instrument master repository
func NewInstrumentMasterRepository(db *sqlx.DB) repository.InstrumentMasterRepository {
	return &InstrumentMasterRepository{
		db: db,
	}
}

// FindByID finds an instrument master by ID
func (r *InstrumentMasterRepository) FindByID(ctx context.Context, id uint64) (*model.InstrumentMaster, error) {
	query := `
		SELECT 
			id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry
		FROM instrument_master
		WHERE id = ?
		LIMIT 1
	`

	instrumentMaster := &model.InstrumentMaster{}
	err := r.db.GetContext(ctx, instrumentMaster, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument master by ID: %w", err)
	}

	return instrumentMaster, nil
}

// FindBySymbol finds an instrument master by symbol
func (r *InstrumentMasterRepository) FindBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error) {
	query := `
		SELECT 
			id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry
		FROM instrument_master
		WHERE symbol = ?
		LIMIT 1
	`

	instrumentMaster := &model.InstrumentMaster{}
	err := r.db.GetContext(ctx, instrumentMaster, query, symbol)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument master by symbol: %w", err)
	}

	return instrumentMaster, nil
}

// FindByExchange finds instrument masters by exchange
func (r *InstrumentMasterRepository) FindByExchange(ctx context.Context, exchange string) ([]*model.InstrumentMaster, error) {
	query := `
		SELECT 
			id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry
		FROM instrument_master
		WHERE exchange = ?
	`

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query, exchange)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument masters by exchange: %w", err)
	}

	return instrumentMasters, nil
}

// FindBySegment finds instrument masters by segment
func (r *InstrumentMasterRepository) FindBySegment(ctx context.Context, segment string) ([]*model.InstrumentMaster, error) {
	query := `
		SELECT 
			id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry
		FROM instrument_master
		WHERE segment = ?
	`

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query, segment)
	if err != nil {
		return nil, fmt.Errorf("error querying instrument masters by segment: %w", err)
	}

	return instrumentMasters, nil
}

// FindAll returns all instrument masters
func (r *InstrumentMasterRepository) FindAll(ctx context.Context) ([]*model.InstrumentMaster, error) {
	query := `
		SELECT 
			id, 
			exchange, 
			segment, 
			symbol, 
			instrument, 
			start_date, 
			expiry
		FROM instrument_master
	`

	var instrumentMasters []*model.InstrumentMaster
	err := r.db.SelectContext(ctx, &instrumentMasters, query)
	if err != nil {
		return nil, fmt.Errorf("error querying all instrument masters: %w", err)
	}

	return instrumentMasters, nil
}

// Create creates a new instrument master
func (r *InstrumentMasterRepository) Create(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	query := `
		INSERT INTO instrument_master (
			id,
			exchange,
			segment,
			symbol,
			instrument,
			start_date,
			expiry
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrumentMaster.ID,
		instrumentMaster.Exchange,
		instrumentMaster.Segment,
		instrumentMaster.Symbol,
		instrumentMaster.Name,
		instrumentMaster.StartDate,
		instrumentMaster.Expiry,
	)

	if err != nil {
		return fmt.Errorf("error creating instrument master: %w", err)
	}

	return nil
}

// Update updates an instrument master
func (r *InstrumentMasterRepository) Update(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
	// Note: ClickHouse typically doesn't support direct updates as in traditional RDBMS
	// This is a simplified version assuming you're using ALTER TABLE UPDATE syntax
	query := `
		ALTER TABLE instrument_master UPDATE
			exchange = ?,
			segment = ?,
			symbol = ?,
			instrument = ?,
			start_date = ?,
			expiry = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		instrumentMaster.Exchange,
		instrumentMaster.Segment,
		instrumentMaster.Symbol,
		instrumentMaster.Name,
		instrumentMaster.StartDate,
		instrumentMaster.Expiry,
		instrumentMaster.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating instrument master: %w", err)
	}

	return nil
}

// Delete deletes an instrument master
func (r *InstrumentMasterRepository) Delete(ctx context.Context, id uint64) error {
	query := `
		ALTER TABLE instrument_master DELETE WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting instrument master: %w", err)
	}

	return nil
}
