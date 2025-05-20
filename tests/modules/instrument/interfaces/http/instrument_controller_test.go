package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"algo/internal/modules/instrument/application/dto"
	"algo/internal/modules/instrument/application/service"
	"algo/internal/modules/instrument/domain/model"
	http_controller "algo/internal/modules/instrument/interfaces/http"
	"algo/tests/modules/instrument/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest() (*gin.Engine, *mocks.MockInstrumentRepository, *mocks.MockInstrumentMasterRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{}
	
	return router, mockInstrumentRepo, mockInstrumentMasterRepo
}

func TestGetInstrumentMasterByID_Success(t *testing.T) {
	// Arrange
	router, mockInstrumentRepo, mockInstrumentMasterRepo := setupTest()
	
	id := uint64(1)
	now := time.Now()
	instrumentMaster := &model.InstrumentMaster{
		ID:                 id,
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
	
	mockInstrumentMasterRepo.FindByIDFunc = func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
		assert.Equal(t, id, idParam)
		return instrumentMaster, nil
	}
	
	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)
	controller := http_controller.NewInstrumentController(instrumentService)
	
	router.GET("/api/instrument-masters/:id", controller.GetInstrumentMasterByID)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/instrument-masters/%d", id), nil)
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.True(t, response.Success)
	
	// Extract data
	dataMap, ok := response.Data.(map[string]interface{})
	require.True(t, ok)
	
	assert.Equal(t, float64(id), dataMap["id"])
	assert.Equal(t, float64(12345), dataMap["exchange_security_id"])
	assert.Equal(t, "NSE", dataMap["exchange"])
	assert.Equal(t, "EQUITY", dataMap["segment"])
	assert.Equal(t, "SBIN", dataMap["symbol"])
	
	// Assert new fields
	assert.Equal(t, float64(100), dataMap["exm_exch_id"])
	assert.Equal(t, float64(200), dataMap["smst_security_id"])
	assert.Equal(t, "SBI Bank Limited", dataMap["instrument_name"])
	assert.Equal(t, "SBIN", dataMap["trading_symbol"])
	assert.Equal(t, float64(1), dataMap["lot_units"])
	assert.Equal(t, float64(0), dataMap["strike_price"])
	assert.Equal(t, "", dataMap["option_type"])
	assert.Equal(t, float64(0.05), dataMap["tick_size"])
	assert.Equal(t, "N", dataMap["expiry_flag"])
	assert.Equal(t, "EQ", dataMap["exch_instrument_type"])
	assert.Equal(t, "SBIN", dataMap["symbol_name"])
}

func TestGetInstrumentMasterByID_NotFound(t *testing.T) {
	// Arrange
	router, mockInstrumentRepo, mockInstrumentMasterRepo := setupTest()
	
	id := uint64(99)
	
	mockInstrumentMasterRepo.FindByIDFunc = func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
		assert.Equal(t, id, idParam)
		return nil, errors.New("instrument master not found")
	}
	
	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)
	controller := http_controller.NewInstrumentController(instrumentService)
	
	router.GET("/api/instrument-masters/:id", controller.GetInstrumentMasterByID)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/instrument-masters/%d", id), nil)
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response dto.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "instrument master not found")
}

func TestCreateInstrumentMaster_Success(t *testing.T) {
	// Arrange
	router, mockInstrumentRepo, mockInstrumentMasterRepo := setupTest()
	
	now := time.Now()
	createRequest := dto.CreateInstrumentMasterRequest{
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
	
	// Mock the repository create function to simulate successful creation
	mockInstrumentMasterRepo.CreateFunc = func(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
		assert.Equal(t, createRequest.Exchange, instrumentMaster.Exchange)
		assert.Equal(t, createRequest.Segment, instrumentMaster.Segment)
		assert.Equal(t, createRequest.Symbol, instrumentMaster.Symbol)
		assert.Equal(t, createRequest.Name, instrumentMaster.Name)
		assert.Equal(t, createRequest.StartDate, instrumentMaster.StartDate)
		assert.Equal(t, createRequest.Expiry, instrumentMaster.Expiry)
		
		// Verify new fields
		assert.Equal(t, createRequest.ExmExchID, instrumentMaster.ExmExchID)
		assert.Equal(t, createRequest.SmstSecurityID, instrumentMaster.SmstSecurityID)
		assert.Equal(t, createRequest.InstrumentName, instrumentMaster.InstrumentName)
		assert.Equal(t, createRequest.TradingSymbol, instrumentMaster.TradingSymbol)
		assert.Equal(t, createRequest.LotUnits, instrumentMaster.LotUnits)
		assert.Equal(t, createRequest.StrikePrice, instrumentMaster.StrikePrice)
		assert.Equal(t, createRequest.OptionType, instrumentMaster.OptionType)
		assert.Equal(t, createRequest.TickSize, instrumentMaster.TickSize)
		assert.Equal(t, createRequest.ExpiryFlag, instrumentMaster.ExpiryFlag)
		assert.Equal(t, createRequest.ExchInstrumentType, instrumentMaster.ExchInstrumentType)
		assert.Equal(t, createRequest.SymbolName, instrumentMaster.SymbolName)
		
		// Set the ID as if it was auto-generated
		instrumentMaster.ID = 1
		return nil
	}
	
	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)
	controller := http_controller.NewInstrumentController(instrumentService)
	
	router.POST("/api/instrument-masters", controller.CreateInstrumentMaster)
	
	// Serialize request to JSON
	requestBody, err := json.Marshal(createRequest)
	require.NoError(t, err)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/instrument-masters", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response dto.ApiResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.True(t, response.Success)
	assert.Equal(t, "Instrument master created successfully", response.Message)
	
	// Extract data
	dataMap, ok := response.Data.(map[string]interface{})
	require.True(t, ok)
	
	assert.Equal(t, float64(1), dataMap["id"]) // The auto-generated ID
	assert.Equal(t, "NSE", dataMap["exchange"])
	assert.Equal(t, "EQUITY", dataMap["segment"])
	assert.Equal(t, "SBIN", dataMap["symbol"])
	
	// Assert new fields in response
	assert.Equal(t, float64(100), dataMap["exm_exch_id"])
	assert.Equal(t, float64(200), dataMap["smst_security_id"])
	assert.Equal(t, "SBI Bank Limited", dataMap["instrument_name"])
	assert.Equal(t, "SBIN", dataMap["trading_symbol"])
	assert.Equal(t, float64(1), dataMap["lot_units"])
	assert.Equal(t, float64(0), dataMap["strike_price"])
	assert.Equal(t, "", dataMap["option_type"])
	assert.Equal(t, float64(0.05), dataMap["tick_size"])
	assert.Equal(t, "N", dataMap["expiry_flag"])
	assert.Equal(t, "EQ", dataMap["exch_instrument_type"])
	assert.Equal(t, "SBIN", dataMap["symbol_name"])
}

func TestUpdateInstrumentMaster_Success(t *testing.T) {
	// Arrange
	router, mockInstrumentRepo, mockInstrumentMasterRepo := setupTest()
	
	id := uint64(1)
	now := time.Now()
	updateRequest := dto.UpdateInstrumentMasterRequest{
		Exchange:           "NSE",
		Segment:            "EQUITY",
		Symbol:             "SBIN",
		Name:               "State Bank of India",
		StartDate:          now,
		ExmExchID:          200, // Updated value
		SmstSecurityID:     300, // Updated value
		InstrumentName:     "Updated Bank Limited", // Updated value
		TradingSymbol:      "UPD", // Updated value
		LotUnits:           2, // Updated value
		StrikePrice:        10.5, // Updated value
		OptionType:         "CE", // Updated value
		TickSize:           0.1, // Updated value
		ExpiryFlag:         "Y", // Updated value
		ExchInstrumentType: "OPT", // Updated value
		SymbolName:         "UPDATED", // Updated value
	}
	
	// The existing instrument master
	existingInstrumentMaster := &model.InstrumentMaster{
		ID:                 id,
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
	
	// Mock the repository functions
	mockInstrumentMasterRepo.FindByIDFunc = func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
		assert.Equal(t, id, idParam)
		return existingInstrumentMaster, nil
	}
	
	mockInstrumentMasterRepo.UpdateFunc = func(ctx context.Context, instrumentMaster *model.InstrumentMaster) error {
		assert.Equal(t, id, instrumentMaster.ID)
		
		// Verify fields were updated
		assert.Equal(t, updateRequest.Exchange, instrumentMaster.Exchange)
		assert.Equal(t, updateRequest.Symbol, instrumentMaster.Symbol)
		
		// Verify new fields were updated
		assert.Equal(t, updateRequest.ExmExchID, instrumentMaster.ExmExchID)
		assert.Equal(t, updateRequest.SmstSecurityID, instrumentMaster.SmstSecurityID)
		assert.Equal(t, updateRequest.InstrumentName, instrumentMaster.InstrumentName)
		assert.Equal(t, updateRequest.TradingSymbol, instrumentMaster.TradingSymbol)
		assert.Equal(t, updateRequest.LotUnits, instrumentMaster.LotUnits)
		assert.Equal(t, updateRequest.StrikePrice, instrumentMaster.StrikePrice)
		assert.Equal(t, updateRequest.OptionType, instrumentMaster.OptionType)
		assert.Equal(t, updateRequest.TickSize, instrumentMaster.TickSize)
		assert.Equal(t, updateRequest.ExpiryFlag, instrumentMaster.ExpiryFlag)
		assert.Equal(t, updateRequest.ExchInstrumentType, instrumentMaster.ExchInstrumentType)
		assert.Equal(t, updateRequest.SymbolName, instrumentMaster.SymbolName)
		
		return nil
	}
	
	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)
	controller := http_controller.NewInstrumentController(instrumentService)
	
	router.PUT("/api/instrument-masters/:id", controller.UpdateInstrumentMaster)
	
	// Serialize request to JSON
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/instrument-masters/%d", id), bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.ApiResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.True(t, response.Success)
	assert.Equal(t, "Instrument master updated successfully", response.Message)
	
	// Extract data
	dataMap, ok := response.Data.(map[string]interface{})
	require.True(t, ok)
	
	assert.Equal(t, float64(id), dataMap["id"])
	
	// Assert updated fields in response
	assert.Equal(t, float64(200), dataMap["exm_exch_id"])
	assert.Equal(t, float64(300), dataMap["smst_security_id"])
	assert.Equal(t, "Updated Bank Limited", dataMap["instrument_name"])
	assert.Equal(t, "UPD", dataMap["trading_symbol"])
	assert.Equal(t, float64(2), dataMap["lot_units"])
	assert.Equal(t, float64(10.5), dataMap["strike_price"])
	assert.Equal(t, "CE", dataMap["option_type"])
	assert.Equal(t, float64(0.1), dataMap["tick_size"])
	assert.Equal(t, "Y", dataMap["expiry_flag"])
	assert.Equal(t, "OPT", dataMap["exch_instrument_type"])
	assert.Equal(t, "UPDATED", dataMap["symbol_name"])
}

func TestGetAllInstrumentMasters_Success(t *testing.T) {
	// Arrange
	router, mockInstrumentRepo, mockInstrumentMasterRepo := setupTest()
	
	now := time.Now()
	instrumentMasters := []*model.InstrumentMaster{
		{
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
			StrikePrice:        0.0,
			OptionType:         "",
			TickSize:           0.05,
			ExpiryFlag:         "N",
			ExchInstrumentType: "EQ",
			SymbolName:         "SBIN",
		},
		{
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
			StrikePrice:        0.0,
			OptionType:         "",
			TickSize:           0.05,
			ExpiryFlag:         "N",
			ExchInstrumentType: "EQ",
			SymbolName:         "HDFC",
		},
	}
	totalCount := 2
	
	// Mock the repository function
	mockInstrumentMasterRepo.FindWithFiltersFunc = func(ctx context.Context, exchange, segment string, offset, limit int) ([]*model.InstrumentMaster, int, error) {
		return instrumentMasters, totalCount, nil
	}
	
	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)
	controller := http_controller.NewInstrumentController(instrumentService)
	
	router.GET("/api/instrument-masters", controller.GetAllInstrumentMasters)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/instrument-masters", nil)
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.True(t, response.Success)
	
	// Extract data
	dataMap, ok := response.Data.(map[string]interface{})
	require.True(t, ok)
	
	// Check pagination info
	assert.Equal(t, float64(1), dataMap["page"])
	assert.Equal(t, float64(10), dataMap["pageSize"]) // Default page size
	assert.Equal(t, float64(2), dataMap["totalCount"])
	assert.Equal(t, float64(1), dataMap["totalPages"])
	
	// Check data array
	data, ok := dataMap["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, data, 2)
	
	// Check first item
	firstItem := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstItem["id"])
	assert.Equal(t, "SBIN", firstItem["symbol"])
	
	// Check new fields
	assert.Equal(t, float64(100), firstItem["exm_exch_id"])
	assert.Equal(t, float64(200), firstItem["smst_security_id"])
	assert.Equal(t, "SBI Bank Limited", firstItem["instrument_name"])
	assert.Equal(t, "SBIN", firstItem["trading_symbol"])
	assert.Equal(t, float64(1), firstItem["lot_units"])
	assert.Equal(t, float64(0), firstItem["strike_price"])
	assert.Equal(t, "", firstItem["option_type"])
	assert.Equal(t, float64(0.05), firstItem["tick_size"])
	assert.Equal(t, "N", firstItem["expiry_flag"])
	assert.Equal(t, "EQ", firstItem["exch_instrument_type"])
	assert.Equal(t, "SBIN", firstItem["symbol_name"])
	
	// Check second item
	secondItem := data[1].(map[string]interface{})
	assert.Equal(t, float64(2), secondItem["id"])
	assert.Equal(t, "HDFC", secondItem["symbol"])
	
	// Check new fields
	assert.Equal(t, float64(101), secondItem["exm_exch_id"])
	assert.Equal(t, float64(201), secondItem["smst_security_id"])
	assert.Equal(t, "HDFC Bank Limited", secondItem["instrument_name"])
	assert.Equal(t, "HDFC", secondItem["trading_symbol"])
}