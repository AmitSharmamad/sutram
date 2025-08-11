package repository

import (
	"context"
	"time"

	"algo/internal/modules/instrument/domain/model"
)

// InstrumentMasterRepository provides access to instrument master data
type InstrumentMasterRepository interface {
	// FindByExchangeSecurityID finds an instrument master by exchange security ID (primary key)
	FindByExchangeSecurityID(ctx context.Context, exchangeSecurityID int32) (*model.InstrumentMaster, error)

	// FindBySymbol finds an instrument master by symbol
	FindBySymbol(ctx context.Context, symbol string) (*model.InstrumentMaster, error)

	// FindAll returns all instrument masters
	FindAll(ctx context.Context) ([]*model.InstrumentMaster, error)

	// FindWithFilters returns instrument masters with filtering and pagination
	FindWithFilters(ctx context.Context, exchange, segment, query string, offset, limit int) ([]*model.InstrumentMaster, int, error)

	// Create creates a new instrument master
	Create(ctx context.Context, instrumentMaster *model.InstrumentMaster) error

	// Update updates an instrument master by exchange security ID
	Update(ctx context.Context, instrumentMaster *model.InstrumentMaster) error

	// Delete deletes an instrument master by exchange security ID
	Delete(ctx context.Context, exchangeSecurityID int32) error

	// FindByExchangeSecurityIDs finds multiple instrument masters by their exchange security IDs
	FindByExchangeSecurityIDs(ctx context.Context, exchangeSecurityIDs []int32) ([]*model.InstrumentMaster, error)

	// FindOptionsContracts finds options contracts based on filters
	FindOptionsContracts(
		ctx context.Context,
		strikePrice float64,
		optionType string,
		exchange string,
		symbol string,
		expiryDate time.Time,
		underlyingSecurityID int32,
	) ([]*model.InstrumentMaster, error)

	// FindFuturesContracts finds futures contracts based on filters
	FindFuturesContracts(
		ctx context.Context,
		exchange string,
		expiryDate time.Time,
		underlyingSecurityID int32,
	) ([]*model.InstrumentMaster, error)
}
