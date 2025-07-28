package http

import (
	"net/http"
	"strconv"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OptionLegHandler handles HTTP requests for option leg operations
type OptionLegHandler struct {
	optionLegService service.OptionLegService
}

// NewOptionLegHandler creates a new option leg handler
func NewOptionLegHandler(optionLegService service.OptionLegService) *OptionLegHandler {
	return &OptionLegHandler{
		optionLegService: optionLegService,
	}
}

// CreateOptionLeg creates a new option leg
// @Summary Create a new option leg
// @Description Create a new option leg for a strategy
// @Tags option-legs
// @Accept json
// @Produce json
// @Param strategy_id path string true "Strategy ID"
// @Param option_leg body dto.CreateOptionLegRequest true "Option leg data"
// @Success 201 {object} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/strategies/{strategy_id}/option-legs [post]
func (h *OptionLegHandler) CreateOptionLeg(c *gin.Context) {
	strategyIDParam := c.Param("strategy_id")
	strategyID, err := uuid.Parse(strategyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy ID"})
		return
	}

	var req dto.CreateOptionLegRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.optionLegService.CreateOptionLeg(c.Request.Context(), strategyID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOptionLegByID retrieves an option leg by its ID
// @Summary Get option leg by ID
// @Description Get option leg by ID
// @Tags option-legs
// @Produce json
// @Param id path string true "Option leg ID"
// @Success 200 {object} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/{id} [get]
func (h *OptionLegHandler) GetOptionLegByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option leg ID"})
		return
	}

	response, err := h.optionLegService.GetOptionLegByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetOptionLegsByStrategy retrieves all option legs for a strategy
// @Summary Get option legs by strategy
// @Description Get all option legs for a specific strategy
// @Tags option-legs
// @Produce json
// @Param strategy_id path string true "Strategy ID"
// @Success 200 {array} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/strategies/{strategy_id}/option-legs [get]
func (h *OptionLegHandler) GetOptionLegsByStrategy(c *gin.Context) {
	strategyIDParam := c.Param("strategy_id")
	strategyID, err := uuid.Parse(strategyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy ID"})
		return
	}

	responses, err := h.optionLegService.GetOptionLegsByStrategyID(c.Request.Context(), strategyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateOptionLeg updates an option leg
// @Summary Update option leg
// @Description Update an existing option leg
// @Tags option-legs
// @Accept json
// @Produce json
// @Param id path string true "Option leg ID"
// @Param option_leg body dto.UpdateOptionLegRequest true "Updated option leg data"
// @Success 200 {object} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/{id} [put]
func (h *OptionLegHandler) UpdateOptionLeg(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option leg ID"})
		return
	}

	var req dto.UpdateOptionLegRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.optionLegService.UpdateOptionLeg(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteOptionLeg deletes an option leg
// @Summary Delete option leg
// @Description Delete an option leg by ID
// @Tags option-legs
// @Produce json
// @Param id path string true "Option leg ID"
// @Success 204 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/{id} [delete]
func (h *OptionLegHandler) DeleteOptionLeg(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option leg ID"})
		return
	}

	if err := h.optionLegService.DeleteOptionLeg(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Option leg deleted successfully"})
}

// DeleteAllOptionLegsByStrategy deletes all option legs for a strategy
// @Summary Delete all option legs by strategy
// @Description Delete all option legs for a specific strategy
// @Tags option-legs
// @Produce json
// @Param strategy_id path string true "Strategy ID"
// @Success 204 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/strategies/{strategy_id}/option-legs [delete]
func (h *OptionLegHandler) DeleteAllOptionLegsByStrategy(c *gin.Context) {
	strategyIDParam := c.Param("strategy_id")
	strategyID, err := uuid.Parse(strategyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy ID"})
		return
	}

	if err := h.optionLegService.DeleteOptionLegsByStrategyID(c.Request.Context(), strategyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "All option legs deleted successfully"})
}

// GetOptionLegsByIndex retrieves option legs by index
// @Summary Get option legs by index
// @Description Get all option legs for a specific index
// @Tags option-legs
// @Produce json
// @Param index path string true "Index name"
// @Success 200 {array} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/index/{index} [get]
func (h *OptionLegHandler) GetOptionLegsByIndex(c *gin.Context) {
	index := c.Param("index")
	if index == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index parameter is required"})
		return
	}

	responses, err := h.optionLegService.GetOptionLegsByIndex(c.Request.Context(), index)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// SearchOptionLegs searches option legs with filters
// @Summary Search option legs
// @Description Search option legs with various filters
// @Tags option-legs
// @Produce json
// @Param strategy_id query string false "Strategy ID"
// @Param select_index query string false "Select index"
// @Param segment query string false "Segment"
// @Param option_type query string false "Option type"
// @Param action_type query string false "Action type"
// @Param expiry_type query string false "Expiry type"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {array} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/search [get]
func (h *OptionLegHandler) SearchOptionLegs(c *gin.Context) {
	var filter dto.OptionLegFilter

	// Parse query parameters
	if strategyID := c.Query("strategy_id"); strategyID != "" {
		id, err := uuid.Parse(strategyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy ID"})
			return
		}
		filter.StrategyID = &id
	}

	if selectIndex := c.Query("select_index"); selectIndex != "" {
		filter.SelectIndex = &selectIndex
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	filter.Page = page
	filter.PerPage = perPage

	responses, total, err := h.optionLegService.SearchOptionLegs(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     responses,
		"total":    total,
		"page":     filter.Page,
		"per_page": filter.PerPage,
	})
}
