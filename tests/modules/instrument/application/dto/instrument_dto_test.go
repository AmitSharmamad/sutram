package dto_test

import (
	"testing"
	"time"

	"algo/internal/modules/instrument/application/dto"
	
	"github.com/stretchr/testify/assert"
)

func TestInstrumentMasterDTOCreation(t *testing.T) {
	// Arrange
	now := time.Now()
	
	// Act
	instrumentMasterDTO := &dto.InstrumentMasterDTO{
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
	assert.Equal(t, uint64(1), instrumentMasterDTO.ID)
	assert.Equal(t, int32(12345), instrumentMasterDTO.ExchangeSecurityID)
	assert.Equal(t, "NSE", instrumentMasterDTO.Exchange)
	assert.Equal(t, "EQUITY", instrumentMasterDTO.Segment)
	assert.Equal(t, "SBIN", instrumentMasterDTO.Symbol)
	assert.Equal(t, "State Bank of India", instrumentMasterDTO.Name)
	assert.Equal(t, now, instrumentMasterDTO.StartDate)
	assert.Equal(t, now.AddDate(0, 1, 0), instrumentMasterDTO.Expiry)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), instrumentMasterDTO.ExmExchID)
	assert.Equal(t, int32(200), instrumentMasterDTO.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", instrumentMasterDTO.InstrumentName)
	assert.Equal(t, "SBIN", instrumentMasterDTO.TradingSymbol)
	assert.Equal(t, int32(1), instrumentMasterDTO.LotUnits)
	assert.Equal(t, 0.0, instrumentMasterDTO.StrikePrice)
	assert.Equal(t, "", instrumentMasterDTO.OptionType)
	assert.Equal(t, 0.05, instrumentMasterDTO.TickSize)
	assert.Equal(t, "N", instrumentMasterDTO.ExpiryFlag)
	assert.Equal(t, "EQ", instrumentMasterDTO.ExchInstrumentType)
	assert.Equal(t, "SBIN", instrumentMasterDTO.SymbolName)
}

func TestCreateInstrumentMasterRequestCreation(t *testing.T) {
	// Arrange
	now := time.Now()
	
	// Act
	createRequest := &dto.CreateInstrumentMasterRequest{
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
	assert.Equal(t, "NSE", createRequest.Exchange)
	assert.Equal(t, "EQUITY", createRequest.Segment)
	assert.Equal(t, "SBIN", createRequest.Symbol)
	assert.Equal(t, "State Bank of India", createRequest.Name)
	assert.Equal(t, now, createRequest.StartDate)
	assert.Equal(t, now.AddDate(0, 1, 0), createRequest.Expiry)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), createRequest.ExmExchID)
	assert.Equal(t, int32(200), createRequest.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", createRequest.InstrumentName)
	assert.Equal(t, "SBIN", createRequest.TradingSymbol)
	assert.Equal(t, int32(1), createRequest.LotUnits)
	assert.Equal(t, 0.0, createRequest.StrikePrice)
	assert.Equal(t, "", createRequest.OptionType)
	assert.Equal(t, 0.05, createRequest.TickSize)
	assert.Equal(t, "N", createRequest.ExpiryFlag)
	assert.Equal(t, "EQ", createRequest.ExchInstrumentType)
	assert.Equal(t, "SBIN", createRequest.SymbolName)
}

func TestUpdateInstrumentMasterRequestCreation(t *testing.T) {
	// Arrange
	now := time.Now()
	
	// Act
	updateRequest := &dto.UpdateInstrumentMasterRequest{
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
	assert.Equal(t, "NSE", updateRequest.Exchange)
	assert.Equal(t, "EQUITY", updateRequest.Segment)
	assert.Equal(t, "SBIN", updateRequest.Symbol)
	assert.Equal(t, "State Bank of India", updateRequest.Name)
	assert.Equal(t, now, updateRequest.StartDate)
	assert.Equal(t, now.AddDate(0, 1, 0), updateRequest.Expiry)
	
	// Assert new fields added in PR
	assert.Equal(t, int32(100), updateRequest.ExmExchID)
	assert.Equal(t, int32(200), updateRequest.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", updateRequest.InstrumentName)
	assert.Equal(t, "SBIN", updateRequest.TradingSymbol)
	assert.Equal(t, int32(1), updateRequest.LotUnits)
	assert.Equal(t, 0.0, updateRequest.StrikePrice)
	assert.Equal(t, "", updateRequest.OptionType)
	assert.Equal(t, 0.05, updateRequest.TickSize)
	assert.Equal(t, "N", updateRequest.ExpiryFlag)
	assert.Equal(t, "EQ", updateRequest.ExchInstrumentType)
	assert.Equal(t, "SBIN", updateRequest.SymbolName)
}

func TestPaginatedResponseWithInstrumentMasters(t *testing.T) {
	// Arrange
	now := time.Now()
	
	instrumentMasters := []interface{}{
		&dto.InstrumentMasterDTO{
			ID:                 1,
			ExchangeSecurityID: 12345,
			Exchange:           "NSE",
			Segment:            "EQUITY",
			Symbol:             "SBIN",
			Name:               "State Bank of India",
			StartDate:          now,
			ExmExchID:          100,
			SmstSecurityID:     200,
			InstrumentName:     "SBI Bank Limited",
			TradingSymbol:      "SBIN",
			LotUnits:           1,
			TickSize:           0.05,
			ExpiryFlag:         "N",
			ExchInstrumentType: "EQ",
			SymbolName:         "SBIN",
		},
		&dto.InstrumentMasterDTO{
			ID:                 2,
			ExchangeSecurityID: 12346,
			Exchange:           "NSE",
			Segment:            "EQUITY",
			Symbol:             "HDFC",
			Name:               "HDFC Bank",
			StartDate:          now,
			ExmExchID:          101,
			SmstSecurityID:     201,
			InstrumentName:     "HDFC Bank Limited",
			TradingSymbol:      "HDFC",
			LotUnits:           1,
			TickSize:           0.05,
			ExpiryFlag:         "N",
			ExchInstrumentType: "EQ",
			SymbolName:         "HDFC",
		},
	}
	
	// Act
	paginatedResponse := &dto.PaginatedResponse{
		Data:       instrumentMasters,
		Page:       1,
		PageSize:   10,
		TotalCount: 2,
		TotalPages: 1,
	}
	
	// Assert
	assert.Equal(t, 1, paginatedResponse.Page)
	assert.Equal(t, 10, paginatedResponse.PageSize)
	assert.Equal(t, 2, paginatedResponse.TotalCount)
	assert.Equal(t, 1, paginatedResponse.TotalPages)
	
	// Assert data conversion
	data := paginatedResponse.Data.([]interface{})
	assert.Len(t, data, 2)
	
	master1 := data[0].(*dto.InstrumentMasterDTO)
	assert.Equal(t, uint64(1), master1.ID)
	assert.Equal(t, int32(100), master1.ExmExchID)
	assert.Equal(t, int32(200), master1.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", master1.InstrumentName)
	assert.Equal(t, "SBIN", master1.TradingSymbol)
	
	master2 := data[1].(*dto.InstrumentMasterDTO)
	assert.Equal(t, uint64(2), master2.ID)
	assert.Equal(t, int32(101), master2.ExmExchID)
	assert.Equal(t, int32(201), master2.SmstSecurityID)
	assert.Equal(t, "HDFC Bank Limited", master2.InstrumentName)
	assert.Equal(t, "HDFC", master2.TradingSymbol)
}

func TestApiResponseWithInstrumentMaster(t *testing.T) {
	// Arrange
	now := time.Now()
	
	instrumentMasterDTO := &dto.InstrumentMasterDTO{
		ID:                 1,
		ExchangeSecurityID: 12345,
		Exchange:           "NSE",
		Segment:            "EQUITY",
		Symbol:             "SBIN",
		Name:               "State Bank of India",
		StartDate:          now,
		ExmExchID:          100,
		SmstSecurityID:     200,
		InstrumentName:     "SBI Bank Limited",
		TradingSymbol:      "SBIN",
		LotUnits:           1,
		TickSize:           0.05,
		ExpiryFlag:         "N",
		ExchInstrumentType: "EQ",
		SymbolName:         "SBIN",
	}
	
	// Act
	apiResponse := &dto.ApiResponse{
		Success: true,
		Message: "Instrument master retrieved successfully",
		Data:    instrumentMasterDTO,
	}
	
	// Assert
	assert.True(t, apiResponse.Success)
	assert.Equal(t, "Instrument master retrieved successfully", apiResponse.Message)
	
	// Assert data conversion
	responseData := apiResponse.Data.(*dto.InstrumentMasterDTO)
	assert.Equal(t, uint64(1), responseData.ID)
	assert.Equal(t, int32(12345), responseData.ExchangeSecurityID)
	assert.Equal(t, "NSE", responseData.Exchange)
	
	// Assert new fields
	assert.Equal(t, int32(100), responseData.ExmExchID)
	assert.Equal(t, int32(200), responseData.SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", responseData.InstrumentName)
	assert.Equal(t, "SBIN", responseData.TradingSymbol)
	assert.Equal(t, int32(1), responseData.LotUnits)
	assert.Equal(t, 0.05, responseData.TickSize)
	assert.Equal(t, "N", responseData.ExpiryFlag)
	assert.Equal(t, "EQ", responseData.ExchInstrumentType)
	assert.Equal(t, "SBIN", responseData.SymbolName)
}