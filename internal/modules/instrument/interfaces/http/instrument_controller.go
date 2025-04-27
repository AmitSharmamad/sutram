package http

import (
	"net/http"
	"strconv"

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

// GetInstrumentByID gets an instrument by ID
func (c *InstrumentController) GetInstrumentByID(ctx *gin.Context) {
	// Get instrument ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid instrument ID",
		})
		return
	}

	// Get the instrument
	instrument, err := c.instrumentService.GetInstrumentByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert instrument model to DTO
	instrumentDTO := modelToDTO(instrument)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    instrumentDTO,
	})
}

// GetInstrumentMasterByID gets an instrument master by ID
func (c *InstrumentController) GetInstrumentMasterByID(ctx *gin.Context) {
	// Get instrument master ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid instrument master ID",
		})
		return
	}

	// Get the instrument master
	instrumentMaster, err := c.instrumentService.GetInstrumentMasterByID(ctx.Request.Context(), id)
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

// GetInstrumentsBySymbol gets instruments by symbol
func (c *InstrumentController) GetInstrumentsBySymbol(ctx *gin.Context) {
	// Get symbol from URL parameter
	symbol := ctx.Param("symbol")

	// Get instruments by symbol
	instruments, err := c.instrumentService.GetInstrumentsBySymbol(ctx.Request.Context(), symbol)
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

// GetInstrumentsByTimeRange gets instruments by symbol within a time range
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
	instruments, err := c.instrumentService.GetInstrumentsByTimeRange(
		ctx.Request.Context(),
		req.Symbol,
		req.StartTime,
		req.EndTime,
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

// GetAllInstrumentMasters gets all instrument masters with filtering and pagination
func (c *InstrumentController) GetAllInstrumentMasters(ctx *gin.Context) {
	// Parse query parameters for filtering and pagination
	var filter dto.InstrumentMasterFilterDTO
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Set default values if not provided
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}

	// Get instrument masters with filters
	instrumentMasters, totalCount, err := c.instrumentService.GetInstrumentMastersWithFilters(
		ctx.Request.Context(),
		filter.Exchange,
		filter.Segment,
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
		Symbol:    req.Symbol,
		Open:      req.Open,
		High:      req.High,
		Low:       req.Low,
		Close:     req.Close,
		Volume:    req.Volume,
		Timestamp: req.Timestamp,
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
	// Get instrument ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid instrument ID",
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

	// Get existing instrument
	instrument, err := c.instrumentService.GetInstrumentByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Update instrument fields
	instrument.Open = req.Open
	instrument.High = req.High
	instrument.Low = req.Low
	instrument.Close = req.Close
	instrument.Volume = req.Volume
	instrument.Timestamp = req.Timestamp

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
		Exchange:  req.Exchange,
		Segment:   req.Segment,
		Symbol:    req.Symbol,
		Name:      req.Name,
		StartDate: req.StartDate,
		Expiry:    req.Expiry,
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
	// Get instrument master ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid instrument master ID",
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
	instrumentMaster, err := c.instrumentService.GetInstrumentMasterByID(ctx.Request.Context(), id)
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
	if !req.Expiry.IsZero() {
		instrumentMaster.Expiry = req.Expiry
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

// DeleteInstrument deletes an instrument
func (c *InstrumentController) DeleteInstrument(ctx *gin.Context) {
	// Get instrument ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid instrument ID",
		})
		return
	}

	// Delete the instrument
	err = c.instrumentService.DeleteInstrument(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete instrument: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Instrument deleted successfully",
	})
}

// DeleteInstrumentMaster deletes an instrument master
func (c *InstrumentController) DeleteInstrumentMaster(ctx *gin.Context) {
	// Get instrument master ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid instrument master ID",
		})
		return
	}

	// Delete the instrument master
	err = c.instrumentService.DeleteInstrumentMaster(ctx.Request.Context(), id)
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
		ID:        instrument.ID,
		Symbol:    instrument.Symbol,
		Open:      instrument.Open,
		High:      instrument.High,
		Low:       instrument.Low,
		Close:     instrument.Close,
		Volume:    instrument.Volume,
		Timestamp: instrument.Timestamp,
	}
}

// Helper function to convert instrument master model to DTO
func masterModelToDTO(master *model.InstrumentMaster) *dto.InstrumentMasterDTO {
	if master == nil {
		return nil
	}

	return &dto.InstrumentMasterDTO{
		ID:        master.ID,
		Exchange:  master.Exchange,
		Segment:   master.Segment,
		Symbol:    master.Symbol,
		Name:      master.Name,
		StartDate: master.StartDate,
		Expiry:    master.Expiry,
	}
}
