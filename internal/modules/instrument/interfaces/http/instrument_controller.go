package http

import (
	"net/http"
	"strconv"
	"strings"

	// Import time for parsing timestamps in other functions

	"algo/internal/modules/instrument/application/dto"
	"algo/internal/modules/instrument/application/service"
	"algo/internal/modules/instrument/domain/model"

	"github.com/gin-gonic/gin"
)

// InstrumentController handles instrument-related HTTP requests
type InstrumentController struct {
	instrumentService *service.InstrumentService
}

// NewInstrumentController creates a new instrument controller
func NewInstrumentController(instrumentService *service.InstrumentService) *InstrumentController {
	return &InstrumentController{
		instrumentService: instrumentService,
	}
}

// GetInstrumentByExchangeSecurityID gets instruments by exchange security ID with pagination
func (c *InstrumentController) GetInstrumentByExchangeSecurityID(ctx *gin.Context) {
	// Get exchange security ID from URL parameter
	idParam := ctx.Param("id")
	exchangeSecurityID, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid exchange security ID",
		})
		return
	}

	// Get pagination parameters
	var pagination dto.GetInstrumentsByExchangeSecurityIDRequest
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid pagination parameters",
		})
		return
	}

	// Default values for pagination
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.PageSize <= 0 {
		pagination.PageSize = 50
	}
	if pagination.PageSize > 1000 {
		pagination.PageSize = 1000
	}

	// Get the paginated instruments
	instruments, totalCount, err := c.instrumentService.GetInstrumentsByExchangeSecurityIDPaginated(
		ctx.Request.Context(),
		int32(exchangeSecurityID),
		pagination.Page,
		pagination.PageSize,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert instrument models to DTOs
	instrumentDTOs := make([]*dto.InstrumentDTO, len(instruments))
	for i, instrument := range instruments {
		instrumentDTOs[i] = modelToDTO(instrument)
	}

	// Calculate total pages
	totalPages := int(totalCount) / pagination.PageSize
	if int(totalCount)%pagination.PageSize > 0 {
		totalPages++
	}

	// Create pagination response
	paginationInfo := dto.PaginationDTO{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}

	// Prepare data in the required format - separate arrays for each property
	openValues := make([]float64, len(instruments))
	highValues := make([]float64, len(instruments))
	lowValues := make([]float64, len(instruments))
	closeValues := make([]float64, len(instruments))
	volumeValues := make([]float64, len(instruments))
	// Using Unix timestamps (seconds since epoch) for better compatibility with frontend
	timestampValues := make([]int64, len(instruments))

	for i, instrument := range instruments {
		openValues[i] = instrument.Open
		highValues[i] = instrument.High
		lowValues[i] = instrument.Low
		closeValues[i] = instrument.Close
		volumeValues[i] = instrument.Volume
		timestampValues[i] = instrument.Timestamp.Unix()
	}

	// Create response using the new required format
	instrumentData := struct {
		Open      []float64 `json:"open"`
		High      []float64 `json:"high"`
		Low       []float64 `json:"low"`
		Close     []float64 `json:"close"`
		Volume    []float64 `json:"volume"`
		Timestamp []int64   `json:"timestamp"`
	}{
		Open:      openValues,
		High:      highValues,
		Low:       lowValues,
		Close:     closeValues,
		Volume:    volumeValues,
		Timestamp: timestampValues,
	}

	// Wrap the data in the outer structure
	response := struct {
		Data       interface{}       `json:"data"`
		Pagination dto.PaginationDTO `json:"pagination"`
	}{
		Data:       instrumentData,
		Pagination: paginationInfo,
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    response,
	})
}

// GetInstrumentMasterByExchangeSecurityID gets an instrument master by exchange security ID
func (c *InstrumentController) GetInstrumentMasterByExchangeSecurityID(ctx *gin.Context) {
	// Get exchange security ID from URL parameter
	idParam := ctx.Param("exchange_security_id")
	exchangeSecurityID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid exchange security ID",
		})
		return
	}

	// Get the instrument master
	instrumentMaster, err := c.instrumentService.GetInstrumentMasterByExchangeSecurityID(ctx.Request.Context(), exchangeSecurityID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert instrument master model to DTO
	instrumentMasterDTO := masterModelToDTO(instrumentMaster)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    instrumentMasterDTO,
	})
}

// GetInstrumentsByTimeRange gets instruments by exchange security ID within a time range
func (c *InstrumentController) GetInstrumentsByTimeRange(ctx *gin.Context) {
	var req dto.GetInstrumentsByTimeRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get instruments by time range
	timeFrame := req.TimeFrame
	if timeFrame == "" {
		timeFrame = "1m" // Default to 1-minute timeframe if not specified
	}

	instruments, err := c.instrumentService.GetInstrumentsByTimeRange(
		ctx.Request.Context(),
		req.ExchangeSecurityID,
		req.StartTime,
		req.EndTime,
		timeFrame,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get instruments: " + err.Error(),
		})
		return
	}

	// Convert instrument models to DTOs
	instrumentDTOs := make([]*dto.InstrumentDTO, len(instruments))
	for i, instrument := range instruments {
		instrumentDTOs[i] = modelToDTO(instrument)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    instrumentDTOs,
	})
}

// GetAllInstrumentMasters gets all instrument masters with pagination and filtering
func (c *InstrumentController) GetAllInstrumentMasters(ctx *gin.Context) {
	var filter dto.InstrumentMasterFilterDTO
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	// Get instrument masters with filtering and pagination
	instrumentMasters, totalCount, err := c.instrumentService.GetInstrumentMastersWithFilters(
		ctx.Request.Context(),
		filter.Exchange,
		filter.Segment,
		filter.Query,
		filter.Page,
		filter.PageSize,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get instrument masters: " + err.Error(),
		})
		return
	}

	// Convert instrument master models to DTOs
	instrumentMasterDTOs := make([]*dto.InstrumentMasterDTO, len(instrumentMasters))
	for i, instrumentMaster := range instrumentMasters {
		instrumentMasterDTOs[i] = masterModelToDTO(instrumentMaster)
	}

	// Calculate total pages
	totalPages := totalCount / filter.PageSize
	if totalCount%filter.PageSize > 0 {
		totalPages++
	}

	// Create paginated response
	paginatedResponse := dto.PaginatedResponse{
		Data:       instrumentMasterDTOs,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    paginatedResponse,
	})
}

// GetInstrumentMastersByExchangeSecurityIDs gets instrument masters by their exchange security IDs
func (c *InstrumentController) GetInstrumentMastersByExchangeSecurityIDs(ctx *gin.Context) {
	// Get exchange security IDs from query parameter - expected format: ?exchange_security_ids=1,2,3,4
	idsParam := ctx.Query("exchange_security_ids")
	if idsParam == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Missing required query parameter: exchange_security_ids",
		})
		return
	}

	// Split the comma-separated IDs and parse them to int64
	idStrings := strings.Split(idsParam, ",")
	exchangeSecurityIDs := make([]int64, 0, len(idStrings))

	for _, idStr := range idStrings {
		exchangeSecurityID, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
				Success: false,
				Error:   "Invalid exchange security ID format: " + idStr,
			})
			return
		}
		exchangeSecurityIDs = append(exchangeSecurityIDs, exchangeSecurityID)
	}

	// Get instrument masters
	instrumentMasters, err := c.instrumentService.GetInstrumentMastersByExchangeSecurityIDs(ctx.Request.Context(), exchangeSecurityIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get instrument masters: " + err.Error(),
		})
		return
	}

	// Convert instrument master models to DTOs
	instrumentMasterDTOs := make([]*dto.InstrumentMasterDTO, len(instrumentMasters))
	for i, instrumentMaster := range instrumentMasters {
		instrumentMasterDTOs[i] = masterModelToDTO(instrumentMaster)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    instrumentMasterDTOs,
	})
}

// CreateInstrument creates a new instrument
func (c *InstrumentController) CreateInstrument(ctx *gin.Context) {
	var req dto.CreateInstrumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create instrument from request
	instrument := &model.Instrument{
		ExchangeSecurityID: req.ExchangeSecurityID,
		Open:               req.Open,
		High:               req.High,
		Low:                req.Low,
		Close:              req.Close,
		Volume:             req.Volume,
		Timestamp:          req.Timestamp,
	}

	// Create the instrument
	err := c.instrumentService.CreateInstrument(ctx.Request.Context(), instrument)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create instrument: " + err.Error(),
		})
		return
	}

	// Convert instrument model to DTO
	instrumentDTO := modelToDTO(instrument)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Instrument created successfully",
		Data:    instrumentDTO,
	})
}

// UpdateInstrument updates an instrument
func (c *InstrumentController) UpdateInstrument(ctx *gin.Context) {
	// Get exchange security ID from URL parameter
	idParam := ctx.Param("id")
	exchangeSecurityID, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid exchange security ID",
		})
		return
	}

	var req dto.UpdateInstrumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// We need the timestamp from the request to identify the specific instrument
	if req.Timestamp.IsZero() {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Timestamp is required for updates",
		})
		return
	}

	// Get existing instruments for this exchange security ID
	instruments, err := c.instrumentService.GetInstrumentsByExchangeSecurityID(ctx.Request.Context(), int32(exchangeSecurityID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Find the specific instrument with matching timestamp
	var instrument *model.Instrument
	for _, inst := range instruments {
		if inst.Timestamp.Equal(req.Timestamp) {
			instrument = inst
			break
		}
	}

	if instrument == nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   "No instrument found with the specified timestamp",
		})
		return
	}

	// Update instrument fields
	instrument.Open = req.Open
	instrument.High = req.High
	instrument.Low = req.Low
	instrument.Close = req.Close
	instrument.Volume = req.Volume

	// Update the instrument
	err = c.instrumentService.UpdateInstrument(ctx.Request.Context(), instrument)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update instrument: " + err.Error(),
		})
		return
	}

	// Convert instrument model to DTO
	instrumentDTO := modelToDTO(instrument)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Instrument updated successfully",
		Data:    instrumentDTO,
	})
}

// CreateInstrumentMaster creates a new instrument master
func (c *InstrumentController) CreateInstrumentMaster(ctx *gin.Context) {
	var req dto.CreateInstrumentMasterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create instrument master from request
	instrumentMaster := &model.InstrumentMaster{
		Exchange:             req.Exchange,
		Segment:              req.Segment,
		Symbol:               req.Symbol,
		Name:                 req.Name,
		StartDate:            req.StartDate,
		Expiry:               req.Expiry,
		ExmExchID:            req.ExmExchID,
		SmstSecurityID:       req.SmstSecurityID,
		InstrumentName:       req.InstrumentName,
		TradingSymbol:        req.TradingSymbol,
		LotUnits:             req.LotUnits,
		StrikePrice:          req.StrikePrice,
		OptionType:           req.OptionType,
		TickSize:             req.TickSize,
		ExpiryFlag:           req.ExpiryFlag,
		ExchInstrumentType:   req.ExchInstrumentType,
		SymbolName:           req.SymbolName,
		UnderlyingSecurityID: req.UnderlyingSecurityID,
		UnderlyingSymbol:     req.UnderlyingSymbol,
	}

	// Create the instrument master
	err := c.instrumentService.CreateInstrumentMaster(ctx.Request.Context(), instrumentMaster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create instrument master: " + err.Error(),
		})
		return
	}

	// Convert instrument master model to DTO
	instrumentMasterDTO := masterModelToDTO(instrumentMaster)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Instrument master created successfully",
		Data:    instrumentMasterDTO,
	})
}

// UpdateInstrumentMaster updates an instrument master
func (c *InstrumentController) UpdateInstrumentMaster(ctx *gin.Context) {
	// Get exchange security ID from URL parameter
	idParam := ctx.Param("exchange_security_id")
	exchangeSecurityID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid exchange security ID",
		})
		return
	}

	var req dto.UpdateInstrumentMasterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get existing instrument master
	instrumentMaster, err := c.instrumentService.GetInstrumentMasterByExchangeSecurityID(ctx.Request.Context(), exchangeSecurityID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Update instrument master fields
	if req.Exchange != "" {
		instrumentMaster.Exchange = req.Exchange
	}
	if req.Segment != "" {
		instrumentMaster.Segment = req.Segment
	}
	if req.Symbol != "" {
		instrumentMaster.Symbol = req.Symbol
	}
	if req.Name != "" {
		instrumentMaster.Name = req.Name
	}
	if !req.StartDate.IsZero() {
		instrumentMaster.StartDate = req.StartDate
	}
	if req.Expiry != nil {
		instrumentMaster.Expiry = req.Expiry
	}
	if req.ExmExchID != 0 {
		instrumentMaster.ExmExchID = req.ExmExchID
	}
	if req.SmstSecurityID != 0 {
		instrumentMaster.SmstSecurityID = req.SmstSecurityID
	}
	if req.InstrumentName != "" {
		instrumentMaster.InstrumentName = req.InstrumentName
	}
	if req.TradingSymbol != "" {
		instrumentMaster.TradingSymbol = req.TradingSymbol
	}
	if req.LotUnits != 0 {
		instrumentMaster.LotUnits = req.LotUnits
	}
	if req.StrikePrice != 0 {
		instrumentMaster.StrikePrice = req.StrikePrice
	}
	if req.OptionType != "" {
		instrumentMaster.OptionType = req.OptionType
	}
	if req.TickSize != 0 {
		instrumentMaster.TickSize = req.TickSize
	}
	if req.ExpiryFlag != "" {
		instrumentMaster.ExpiryFlag = req.ExpiryFlag
	}
	if req.ExchInstrumentType != "" {
		instrumentMaster.ExchInstrumentType = req.ExchInstrumentType
	}
	if req.SymbolName != "" {
		instrumentMaster.SymbolName = req.SymbolName
	}
	if req.UnderlyingSecurityID != 0 {
		instrumentMaster.UnderlyingSecurityID = req.UnderlyingSecurityID
	}
	if req.UnderlyingSymbol != "" {
		instrumentMaster.UnderlyingSymbol = req.UnderlyingSymbol
	}

	// Update the instrument master
	err = c.instrumentService.UpdateInstrumentMaster(ctx.Request.Context(), instrumentMaster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update instrument master: " + err.Error(),
		})
		return
	}

	// Convert instrument master model to DTO
	instrumentMasterDTO := masterModelToDTO(instrumentMaster)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Instrument master updated successfully",
		Data:    instrumentMasterDTO,
	})
}

// DeleteInstrument deletes instruments by exchange security ID
func (c *InstrumentController) DeleteInstrument(ctx *gin.Context) {
	// Get exchange security ID from URL parameter
	idParam := ctx.Param("id")
	exchangeSecurityID, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid exchange security ID",
		})
		return
	}

	// Delete the instruments with the specified exchange security ID
	err = c.instrumentService.DeleteInstrument(ctx.Request.Context(), int32(exchangeSecurityID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete instruments: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Instruments deleted successfully",
	})
}

// DeleteInstrumentMaster deletes an instrument master
func (c *InstrumentController) DeleteInstrumentMaster(ctx *gin.Context) {
	// Get exchange security ID from URL parameter
	idParam := ctx.Param("exchange_security_id")
	exchangeSecurityID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid exchange security ID",
		})
		return
	}

	// Delete the instrument master
	err = c.instrumentService.DeleteInstrumentMaster(ctx.Request.Context(), exchangeSecurityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete instrument master: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Instrument master deleted successfully",
	})
}

// Helper function to convert instrument model to DTO
func modelToDTO(instrument *model.Instrument) *dto.InstrumentDTO {
	if instrument == nil {
		return nil
	}

	return &dto.InstrumentDTO{
		ExchangeSecurityID: instrument.ExchangeSecurityID,
		Open:               instrument.Open,
		High:               instrument.High,
		Low:                instrument.Low,
		Close:              instrument.Close,
		Volume:             instrument.Volume,
		Timestamp:          instrument.Timestamp,
	}
}

// Helper function to convert instrument master model to DTO
func masterModelToDTO(master *model.InstrumentMaster) *dto.InstrumentMasterDTO {
	if master == nil {
		return nil
	}

	return &dto.InstrumentMasterDTO{
		ExchangeSecurityID:   master.ExchangeSecurityID,
		Exchange:             master.Exchange,
		Segment:              master.Segment,
		Symbol:               master.Symbol,
		Name:                 master.Name,
		StartDate:            master.StartDate,
		Expiry:               master.Expiry,
		ExmExchID:            master.ExmExchID,
		SmstSecurityID:       master.SmstSecurityID,
		InstrumentName:       master.InstrumentName,
		TradingSymbol:        master.TradingSymbol,
		LotUnits:             master.LotUnits,
		StrikePrice:          master.StrikePrice,
		OptionType:           master.OptionType,
		TickSize:             master.TickSize,
		ExpiryFlag:           master.ExpiryFlag,
		ExchInstrumentType:   master.ExchInstrumentType,
		SymbolName:           master.SymbolName,
		UnderlyingSecurityID: master.UnderlyingSecurityID,
		UnderlyingSymbol:     master.UnderlyingSymbol,
	}
}
