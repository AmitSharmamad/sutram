package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"algo/internal/modules/instrument/application/service"
	"algo/internal/modules/instrument/domain/model"
	"algo/tests/modules/instrument/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetInstrumentMasterByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	id := uint64(1)
	now := time.Now()

	expectedInstrumentMaster := &model.InstrumentMaster{
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

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindByIDFunc: func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
			assert.Equal(t, id, idParam)
			return expectedInstrumentMaster, nil
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	result, err := instrumentService.GetInstrumentMasterByID(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, expectedInstrumentMaster.Exchange, result.Exchange)
	assert.Equal(t, expectedInstrumentMaster.Symbol, result.Symbol)
	
	// Assert new fields
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
}

func TestGetInstrumentMasterByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	id := uint64(99)

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindByIDFunc: func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
			assert.Equal(t, id, idParam)
			return nil, errors.New("instrument master not found")
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	result, err := instrumentService.GetInstrumentMasterByID(ctx, id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetInstrumentMasterBySymbol_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	symbol := "SBIN"
	now := time.Now()

	expectedInstrumentMaster := &model.InstrumentMaster{
		ID:                 1,
		ExchangeSecurityID: 12345,
		Exchange:           "NSE",
		Segment:            "EQUITY",
		Symbol:             symbol,
		Name:               "State Bank of India",
		StartDate:          now,
		Expiry:             now.AddDate(0, 1, 0),
		ExmExchID:          100,
		SmstSecurityID:     200,
		InstrumentName:     "SBI Bank Limited",
		TradingSymbol:      symbol,
		LotUnits:           1,
		StrikePrice:        0.0,
		OptionType:         "",
		TickSize:           0.05,
		ExpiryFlag:         "N",
		ExchInstrumentType: "EQ",
		SymbolName:         symbol,
	}

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindBySymbolFunc: func(ctx context.Context, symbolParam string) (*model.InstrumentMaster, error) {
			assert.Equal(t, symbol, symbolParam)
			return expectedInstrumentMaster, nil
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	result, err := instrumentService.GetInstrumentMasterBySymbol(ctx, symbol)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	
	// Assert new fields
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
}

func TestCreateInstrumentMaster_Success(t *testing.T) {
	// Arrange
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

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		CreateFunc: func(ctx context.Context, master *model.InstrumentMaster) error {
			assert.Equal(t, instrumentMaster, master)
			return nil
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	err := instrumentService.CreateInstrumentMaster(ctx, instrumentMaster)

	// Assert
	assert.NoError(t, err)
}

func TestCreateInstrumentMaster_Error(t *testing.T) {
	// Arrange
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

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		CreateFunc: func(ctx context.Context, master *model.InstrumentMaster) error {
			return errors.New("database error")
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	err := instrumentService.CreateInstrumentMaster(ctx, instrumentMaster)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating instrument master")
}

func TestUpdateInstrumentMaster_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
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
	
	existingInstrumentMaster := &model.InstrumentMaster{
		ID:                 id,
		ExchangeSecurityID: 12345,
		Exchange:           "NSE",
		Segment:            "EQUITY",
		Symbol:             "SBIN",
		Name:               "State Bank of India",
		StartDate:          now,
		Expiry:             now.AddDate(0, 1, 0),
		// Old field values before update
		ExmExchID:          50,
		SmstSecurityID:     150,
		InstrumentName:     "Old Bank Limited",
		TradingSymbol:      "OLD",
		LotUnits:           5,
		StrikePrice:        10.0,
		OptionType:         "CE",
		TickSize:           0.1,
		ExpiryFlag:         "Y",
		ExchInstrumentType: "XX",
		SymbolName:         "OLDNAME",
	}

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindByIDFunc: func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
			assert.Equal(t, id, idParam)
			return existingInstrumentMaster, nil
		},
		UpdateFunc: func(ctx context.Context, master *model.InstrumentMaster) error {
			assert.Equal(t, instrumentMaster, master)
			return nil
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	err := instrumentService.UpdateInstrumentMaster(ctx, instrumentMaster)

	// Assert
	assert.NoError(t, err)
}

func TestUpdateInstrumentMaster_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
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

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindByIDFunc: func(ctx context.Context, idParam uint64) (*model.InstrumentMaster, error) {
			assert.Equal(t, id, idParam)
			return nil, errors.New("instrument master not found")
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	err := instrumentService.UpdateInstrumentMaster(ctx, instrumentMaster)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error finding instrument master")
}

func TestGetInstrumentMastersWithFilters_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	exchange := "NSE"
	segment := "EQUITY"
	page := 1
	pageSize := 10
	offset := 0
	limit := 10
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

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindWithFiltersFunc: func(ctx context.Context, exchangeParam, segmentParam string, offsetParam, limitParam int) ([]*model.InstrumentMaster, int, error) {
			assert.Equal(t, exchange, exchangeParam)
			assert.Equal(t, segment, segmentParam)
			assert.Equal(t, offset, offsetParam)
			assert.Equal(t, limit, limitParam)
			return instrumentMasters, totalCount, nil
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	results, count, err := instrumentService.GetInstrumentMastersWithFilters(ctx, exchange, segment, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, totalCount, count)
	assert.Len(t, results, 2)
	
	// Check first result
	assert.Equal(t, uint64(1), results[0].ID)
	assert.Equal(t, "SBIN", results[0].Symbol)
	assert.Equal(t, int32(100), results[0].ExmExchID)
	assert.Equal(t, int32(200), results[0].SmstSecurityID)
	assert.Equal(t, "SBI Bank Limited", results[0].InstrumentName)
	assert.Equal(t, "SBIN", results[0].TradingSymbol)
	assert.Equal(t, int32(1), results[0].LotUnits)
	assert.Equal(t, 0.0, results[0].StrikePrice)
	assert.Equal(t, "", results[0].OptionType)
	assert.Equal(t, 0.05, results[0].TickSize)
	assert.Equal(t, "N", results[0].ExpiryFlag)
	assert.Equal(t, "EQ", results[0].ExchInstrumentType)
	assert.Equal(t, "SBIN", results[0].SymbolName)
	
	// Check second result
	assert.Equal(t, uint64(2), results[1].ID)
	assert.Equal(t, "HDFC", results[1].Symbol)
	assert.Equal(t, int32(101), results[1].ExmExchID)
	assert.Equal(t, int32(201), results[1].SmstSecurityID)
	assert.Equal(t, "HDFC Bank Limited", results[1].InstrumentName)
	assert.Equal(t, "HDFC", results[1].TradingSymbol)
}

func TestGetInstrumentMastersByIDs_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ids := []uint64{1, 2}
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

	mockInstrumentRepo := &mocks.MockInstrumentRepository{}
	mockInstrumentMasterRepo := &mocks.MockInstrumentMasterRepository{
		FindByIDsFunc: func(ctx context.Context, idsParam []uint64) ([]*model.InstrumentMaster, error) {
			assert.Equal(t, ids, idsParam)
			return instrumentMasters, nil
		},
	}

	instrumentService := service.NewInstrumentService(mockInstrumentRepo, mockInstrumentMasterRepo)

	// Act
	results, err := instrumentService.GetInstrumentMastersByIDs(ctx, ids)

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
}