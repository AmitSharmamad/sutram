package clickhouse_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"algo/internal/modules/instrument/domain/model"
	repository "algo/internal/modules/instrument/infrastructure/persistence/clickhouse"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindByID(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	id := uint64(1)
	now := time.Now()
	
	rows := sqlmock.NewRows([]string{
		"id", "exchange_security_id", "exchange", "segment", "symbol", "instrument", 
		"start_date", "expiry", "exm_exch_id", "smst_security_id", "instrument_name", 
		"trading_symbol", "lot_units", "strike_price", "option_type", "tick_size", 
		"expiry_flag", "exch_instrument_type", "symbol_name",
	}).AddRow(
		id, 12345, "NSE", "EQUITY", "SBIN", "State Bank of India", 
		now, now.AddDate(0, 1, 0), 100, 200, "SBI Bank Limited", 
		"SBIN", 1, 0.0, "", 0.05, "N", "EQ", "SBIN",
	)
	
	mock.ExpectQuery("SELECT (.+) FROM sutram.instrument_master").
		WithArgs(id).
		WillReturnRows(rows)
	
	// Act
	result, err := repo.FindByID(ctx, id)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, int32(12345), result.ExchangeSecurityID)
	assert.Equal(t, "NSE", result.Exchange)
	assert.Equal(t, "EQUITY", result.Segment)
	assert.Equal(t, "SBIN", result.Symbol)
	assert.Equal(t, "State Bank of India", result.Name)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), result.ExmExchID)
	assert.Equal(t, int32(200), result.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", result.InstrumentName)
	assert.Equal(t, "SBIN", result.TradingSymbol)
	assert.Equal(t, int32(1), result.LotUnits)
	assert.Equal(t, 0.0, result.StrikePrice)
	assert.Equal(t, "", result.OptionType)
	assert.Equal(t, 0.05, result.TickSize)
	assert.Equal(t, "N", result.ExpiryFlag)
	assert.Equal(t, "EQ", result.ExchInstrumentType)
	assert.Equal(t, "SBIN", result.SymbolName)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindBySymbol(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	symbol := "SBIN"
	now := time.Now()
	
	rows := sqlmock.NewRows([]string{
		"id", "exchange_security_id", "exchange", "segment", "symbol", "instrument", 
		"start_date", "expiry", "exm_exch_id", "smst_security_id", "instrument_name", 
		"trading_symbol", "lot_units", "strike_price", "option_type", "tick_size", 
		"expiry_flag", "exch_instrument_type", "symbol_name",
	}).AddRow(
		1, 12345, "NSE", "EQUITY", symbol, "State Bank of India", 
		now, now.AddDate(0, 1, 0), 100, 200, "SBI Bank Limited", 
		symbol, 1, 0.0, "", 0.05, "N", "EQ", symbol,
	)
	
	mock.ExpectQuery("SELECT (.+) FROM sutram.instrument_master").
		WithArgs(symbol).
		WillReturnRows(rows)
	
	// Act
	result, err := repo.FindBySymbol(ctx, symbol)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), result.ExmExchID)
	assert.Equal(t, int32(200), result.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", result.InstrumentName)
	assert.Equal(t, symbol, result.TradingSymbol)
	assert.Equal(t, int32(1), result.LotUnits)
	assert.Equal(t, 0.0, result.StrikePrice)
	assert.Equal(t, "", result.OptionType)
	assert.Equal(t, 0.05, result.TickSize)
	assert.Equal(t, "N", result.ExpiryFlag)
	assert.Equal(t, "EQ", result.ExchInstrumentType)
	assert.Equal(t, symbol, result.SymbolName)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindByExchangeSecurityID(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	exchangeSecurityID := int32(12345)
	now := time.Now()
	
	rows := sqlmock.NewRows([]string{
		"id", "exchange_security_id", "exchange", "segment", "symbol", "instrument", 
		"start_date", "expiry", "exm_exch_id", "smst_security_id", "instrument_name", 
		"trading_symbol", "lot_units", "strike_price", "option_type", "tick_size", 
		"expiry_flag", "exch_instrument_type", "symbol_name",
	}).AddRow(
		1, exchangeSecurityID, "NSE", "EQUITY", "SBIN", "State Bank of India", 
		now, now.AddDate(0, 1, 0), 100, 200, "SBI Bank Limited", 
		"SBIN", 1, 0.0, "", 0.05, "N", "EQ", "SBIN",
	)
	
	mock.ExpectQuery("SELECT (.+) FROM sutram.instrument_master").
		WithArgs(exchangeSecurityID).
		WillReturnRows(rows)
	
	// Act
	result, err := repo.FindByExchangeSecurityID(ctx, exchangeSecurityID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exchangeSecurityID, result.ExchangeSecurityID)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), result.ExmExchID)
	assert.Equal(t, int32(200), result.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", result.InstrumentName)
	assert.Equal(t, "SBIN", result.TradingSymbol)
	assert.Equal(t, int32(1), result.LotUnits)
	assert.Equal(t, 0.0, result.StrikePrice)
	assert.Equal(t, "", result.OptionType)
	assert.Equal(t, 0.05, result.TickSize)
	assert.Equal(t, "N", result.ExpiryFlag)
	assert.Equal(t, "EQ", result.ExchInstrumentType)
	assert.Equal(t, "SBIN", result.SymbolName)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	now := time.Now()
	
	instrumentMaster := &model.InstrumentMaster{
		ID:                 1,
		ExchangeSecurityID: 12345,
		Exchange:           "NSE",
		Segment:            "EQUITY",
		Symbol:             "SBIN",
		Name:               "State Bank of India",
		StartDate:          now,
		Expiry:             now.AddDate(0, 1, 0),
		ExmExchID:          100,
		SmstSecurityID:     200,
		InstrumentName:     "SBI Bank Limited",
		TradingSymbol:      "SBIN",
		LotUnits:           1,
		StrikePrice:        0.0,
		OptionType:         "",
		TickSize:           0.05,
		ExpiryFlag:         "N",
		ExchInstrumentType: "EQ",
		SymbolName:         "SBIN",
	}
	
	mock.ExpectExec("INSERT INTO instrument_master").
		WithArgs(
			instrumentMaster.ID,
			instrumentMaster.ExchangeSecurityID,
			instrumentMaster.Exchange,
			instrumentMaster.Segment,
			instrumentMaster.Symbol,
			instrumentMaster.Name,
			instrumentMaster.StartDate,
			instrumentMaster.Expiry,
			instrumentMaster.ExmExchID,
			instrumentMaster.SmstSecurityID,
			instrumentMaster.InstrumentName,
			instrumentMaster.TradingSymbol,
			instrumentMaster.LotUnits,
			instrumentMaster.StrikePrice,
			instrumentMaster.OptionType,
			instrumentMaster.TickSize,
			instrumentMaster.ExpiryFlag,
			instrumentMaster.ExchInstrumentType,
			instrumentMaster.SymbolName,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Act
	err = repo.Create(ctx, instrumentMaster)
	
	// Assert
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	now := time.Now()
	
	instrumentMaster := &model.InstrumentMaster{
		ID:                 1,
		ExchangeSecurityID: 12345,
		Exchange:           "NSE",
		Segment:            "EQUITY",
		Symbol:             "SBIN",
		Name:               "State Bank of India",
		StartDate:          now,
		Expiry:             now.AddDate(0, 1, 0),
		ExmExchID:          100,
		SmstSecurityID:     200,
		InstrumentName:     "SBI Bank Limited",
		TradingSymbol:      "SBIN",
		LotUnits:           1,
		StrikePrice:        0.0,
		OptionType:         "",
		TickSize:           0.05,
		ExpiryFlag:         "N",
		ExchInstrumentType: "EQ",
		SymbolName:         "SBIN",
	}
	
	mock.ExpectExec("ALTER TABLE instrument_master UPDATE").
		WithArgs(
			instrumentMaster.ExchangeSecurityID,
			instrumentMaster.Exchange,
			instrumentMaster.Segment,
			instrumentMaster.Symbol,
			instrumentMaster.Name,
			instrumentMaster.StartDate,
			instrumentMaster.Expiry,
			instrumentMaster.ExmExchID,
			instrumentMaster.SmstSecurityID,
			instrumentMaster.InstrumentName,
			instrumentMaster.TradingSymbol,
			instrumentMaster.LotUnits,
			instrumentMaster.StrikePrice,
			instrumentMaster.OptionType,
			instrumentMaster.TickSize,
			instrumentMaster.ExpiryFlag,
			instrumentMaster.ExchInstrumentType,
			instrumentMaster.SymbolName,
			instrumentMaster.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Act
	err = repo.Update(ctx, instrumentMaster)
	
	// Assert
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindWithFilters(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	exchange := "NSE"
	segment := "EQUITY"
	offset := 0
	limit := 10
	now := time.Now()
	
	rows := sqlmock.NewRows([]string{
		"id", "exchange_security_id", "exchange", "segment", "symbol", "instrument", 
		"start_date", "expiry", "exm_exch_id", "smst_security_id", "instrument_name", 
		"trading_symbol", "lot_units", "strike_price", "option_type", "tick_size", 
		"expiry_flag", "exch_instrument_type", "symbol_name",
	}).
		AddRow(1, 12345, exchange, segment, "SBIN", "State Bank of India", 
			now, now.AddDate(0, 1, 0), 100, 200, "SBI Bank Limited", 
			"SBIN", 1, 0.0, "", 0.05, "N", "EQ", "SBIN").
		AddRow(2, 12346, exchange, segment, "HDFC", "HDFC Bank", 
			now, now.AddDate(0, 1, 0), 101, 201, "HDFC Bank Limited", 
			"HDFC", 1, 0.0, "", 0.05, "N", "EQ", "HDFC")
	
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	
	mock.ExpectQuery("SELECT (.+) FROM sutram.instrument_master").
		WithArgs(exchange, segment, limit, offset).
		WillReturnRows(rows)
		
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(exchange, segment).
		WillReturnRows(countRows)
	
	// Act
	results, count, err := repo.FindWithFilters(ctx, exchange, segment, offset, limit)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)
	
	// Check first result
	assert.Equal(t, uint64(1), results[0].ID)
	assert.Equal(t, "SBIN", results[0].Symbol)
	assert.Equal(t, int32(100), results[0].ExmExchID)
	assert.Equal(t, int32(200), results[0].SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", results[0].InstrumentName)
	assert.Equal(t, "SBIN", results[0].TradingSymbol)
	
	// Check second result
	assert.Equal(t, uint64(2), results[1].ID)
	assert.Equal(t, "HDFC", results[1].Symbol)
	assert.Equal(t, int32(101), results[1].ExmExchID)
	assert.Equal(t, int32(201), results[1].SmstSecurityID)
	assert.Equal(t, "HDFC Bank Limited", results[1].InstrumentName)
	assert.Equal(t, "HDFC", results[1].TradingSymbol)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindByIDs(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	ids := []uint64{1, 2}
	now := time.Now()
	
	rows := sqlmock.NewRows([]string{
		"id", "exchange_security_id", "exchange", "segment", "symbol", "instrument", 
		"start_date", "expiry", "exm_exch_id", "smst_security_id", "instrument_name", 
		"trading_symbol", "lot_units", "strike_price", "option_type", "tick_size", 
		"expiry_flag", "exch_instrument_type", "symbol_name",
	}).
		AddRow(1, 12345, "NSE", "EQUITY", "SBIN", "State Bank of India", 
			now, now.AddDate(0, 1, 0), 100, 200, "SBI Bank Limited", 
			"SBIN", 1, 0.0, "", 0.05, "N", "EQ", "SBIN").
		AddRow(2, 12346, "NSE", "EQUITY", "HDFC", "HDFC Bank", 
			now, now.AddDate(0, 1, 0), 101, 201, "HDFC Bank Limited", 
			"HDFC", 1, 0.0, "", 0.05, "N", "EQ", "HDFC")
	
	mock.ExpectQuery("SELECT (.+) FROM sutram.instrument_master").
		WillReturnRows(rows)
	
	// Act
	results, err := repo.FindByIDs(ctx, ids)
	
	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Check first result
	assert.Equal(t, uint64(1), results[0].ID)
	assert.Equal(t, int32(100), results[0].ExmExchID)
	assert.Equal(t, int32(200), results[0].SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", results[0].InstrumentName)
	assert.Equal(t, "SBIN", results[0].TradingSymbol)
	
	// Check second result
	assert.Equal(t, uint64(2), results[1].ID)
	assert.Equal(t, int32(101), results[1].ExmExchID)
	assert.Equal(t, int32(201), results[1].SmstSecurityID)
	assert.Equal(t, "HDFC Bank Limited", results[1].InstrumentName)
	assert.Equal(t, "HDFC", results[1].TradingSymbol)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestErrorHandling(t *testing.T) {
	// Arrange
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	
	sqlxDB := sqlx.NewDb(mockDB, "clickhouse")
	repo := repository.NewInstrumentMasterRepository(sqlxDB)
	
	ctx := context.Background()
	id := uint64(1)
	
	mock.ExpectQuery("SELECT (.+) FROM sutram.instrument_master").
		WithArgs(id).
		WillReturnError(errors.New("database error"))
	
	// Act
	result, err := repo.FindByID(ctx, id)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error querying instrument master")
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}