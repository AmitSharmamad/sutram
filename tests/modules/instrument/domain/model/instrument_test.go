package model_test

import (
	"testing"
	"time"

	"algo/internal/modules/instrument/domain/model"
	
	"github.com/stretchr/testify/assert"
)

func TestInstrumentMasterCreation(t *testing.T) {
	// Arrange
	now := time.Now()
	
	// Act
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

	// Assert
	assert.Equal(t, uint64(1), instrumentMaster.ID)
	assert.Equal(t, int32(12345), instrumentMaster.ExchangeSecurityID)
	assert.Equal(t, "NSE", instrumentMaster.Exchange)
	assert.Equal(t, "EQUITY", instrumentMaster.Segment)
	assert.Equal(t, "SBIN", instrumentMaster.Symbol)
	assert.Equal(t, "State Bank of India", instrumentMaster.Name)
	assert.Equal(t, now, instrumentMaster.StartDate)
	assert.Equal(t, now.AddDate(0, 1, 0), instrumentMaster.Expiry)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), instrumentMaster.ExmExchID)
	assert.Equal(t, int32(200), instrumentMaster.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", instrumentMaster.InstrumentName)
	assert.Equal(t, "SBIN", instrumentMaster.TradingSymbol)
	assert.Equal(t, int32(1), instrumentMaster.LotUnits)
	assert.Equal(t, 0.0, instrumentMaster.StrikePrice)
	assert.Equal(t, "", instrumentMaster.OptionType)
	assert.Equal(t, 0.05, instrumentMaster.TickSize)
	assert.Equal(t, "N", instrumentMaster.ExpiryFlag)
	assert.Equal(t, "EQ", instrumentMaster.ExchInstrumentType)
	assert.Equal(t, "SBIN", instrumentMaster.SymbolName)
}

func TestInstrumentWithMetadataCreation(t *testing.T) {
	// Arrange
	now := time.Now()
	
	instrument := &model.Instrument{
		ExchangeSecurityID: 12345,
		Open:               100.0,
		High:               105.0,
		Low:                95.0,
		Close:              102.0,
		Volume:             1000000.0,
		Timestamp:          now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
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
	
	// Act
	instrumentWithMetadata := &model.InstrumentWithMetadata{
		Instrument:       instrument,
		InstrumentMaster: instrumentMaster,
	}
	
	// Assert
	assert.NotNil(t, instrumentWithMetadata)
	assert.Equal(t, instrument, instrumentWithMetadata.Instrument)
	assert.Equal(t, instrumentMaster, instrumentWithMetadata.InstrumentMaster)
	
	// Assert fields in nested objects
	assert.Equal(t, int32(12345), instrumentWithMetadata.Instrument.ExchangeSecurityID)
	assert.Equal(t, int32(12345), instrumentWithMetadata.InstrumentMaster.ExchangeSecurityID)
	
	// Assert new fields in nested objects
	assert.Equal(t, int32(100), instrumentWithMetadata.InstrumentMaster.ExmExchID)
	assert.Equal(t, int32(200), instrumentWithMetadata.InstrumentMaster.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", instrumentWithMetadata.InstrumentMaster.InstrumentName)
	assert.Equal(t, "SBIN", instrumentWithMetadata.InstrumentMaster.TradingSymbol)
	assert.Equal(t, int32(1), instrumentWithMetadata.InstrumentMaster.LotUnits)
	assert.Equal(t, 0.0, instrumentWithMetadata.InstrumentMaster.StrikePrice)
	assert.Equal(t, "", instrumentWithMetadata.InstrumentMaster.OptionType)
	assert.Equal(t, 0.05, instrumentWithMetadata.InstrumentMaster.TickSize)
	assert.Equal(t, "N", instrumentWithMetadata.InstrumentMaster.ExpiryFlag)
	assert.Equal(t, "EQ", instrumentWithMetadata.InstrumentMaster.ExchInstrumentType)
	assert.Equal(t, "SBIN", instrumentWithMetadata.InstrumentMaster.SymbolName)
}

func TestStockDataCreation(t *testing.T) {
	// Arrange
	now := time.Now()
	
	inst1 := model.Instrument{
		ExchangeSecurityID: 12345,
		Open:               100.0,
		High:               105.0,
		Low:                95.0,
		Close:              102.0,
		Volume:             1000000.0,
		Timestamp:          now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	inst2 := model.Instrument{
		ExchangeSecurityID: 12345,
		Open:               102.0,
		High:               107.0,
		Low:                101.0,
		Close:              106.0,
		Volume:             1200000.0,
		Timestamp:          now.Add(time.Hour),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	// Act
	stockData := model.StockData{inst1, inst2}
	
	// Assert
	assert.Len(t, stockData, 2)
	assert.Equal(t, int32(12345), stockData[0].ExchangeSecurityID)
	assert.Equal(t, 100.0, stockData[0].Open)
	assert.Equal(t, 102.0, stockData[1].Open)
}