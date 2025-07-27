package http

import (
	"net/http"
	"strconv"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/application/service"
	"algo/internal/modules/strategy/domain/model"

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
// @Summary Create option leg
// @Description Create a new option leg for a strategy
// @Tags option-legs
// @Accept json
// @Produce json
// @Param option_leg body dto.CreateOptionLegRequest true "Option leg data"
// @Success 201 {object} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs [post]
func (h *OptionLegHandler) CreateOptionLeg(c *gin.Context) {
	var req dto.CreateOptionLegRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.optionLegService.CreateOptionLeg(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOptionLeg retrieves an option leg by ID
// @Summary Get option leg
// @Description Get an option leg by its ID
// @Tags option-legs
// @Produce json
// @Param id path string true "Option leg ID"
// @Success 200 {object} dto.OptionLegResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/{id} [get]
func (h *OptionLegHandler) GetOptionLeg(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option leg ID"})
		return
	}

	response, err := h.optionLegService.GetOptionLegByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOptionLegNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Option leg not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// UpdateOptionLeg updates an existing option leg
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
		if err == service.ErrOptionLegNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Option leg not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteOptionLeg deletes an option leg
// @Summary Delete option leg
// @Description Delete an option leg by its ID
// @Tags option-legs
// @Param id path string true "Option leg ID"
// @Success 204
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

	err = h.optionLegService.DeleteOptionLeg(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOptionLegNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Option leg not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteOptionLegsByStrategy deletes all option legs for a strategy
// @Summary Delete option legs by strategy
// @Description Delete all option legs for a specific strategy
// @Tags option-legs
// @Param strategy_id path string true "Strategy ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/strategies/{strategy_id}/option-legs [delete]
func (h *OptionLegHandler) DeleteOptionLegsByStrategy(c *gin.Context) {
	strategyIDParam := c.Param("strategy_id")
	strategyID, err := uuid.Parse(strategyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy ID"})
		return
	}

	err = h.optionLegService.DeleteOptionLegsByStrategyID(c.Request.Context(), strategyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetOptionLegsByIndex retrieves option legs by index/underlying
// @Summary Get option legs by index
// @Description Get option legs for a specific index/underlying
// @Tags option-legs
// @Produce json
// @Param index path string true "Index/Underlying name"
// @Success 200 {array} dto.OptionLegResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/index/{index} [get]
func (h *OptionLegHandler) GetOptionLegsByIndex(c *gin.Context) {
	index := c.Param("index")

	responses, err := h.optionLegService.GetOptionLegsByIndex(c.Request.Context(), index)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// SearchOptionLegs searches for option legs with filtering
// @Summary Search option legs
// @Description Search option legs with various filter criteria
// @Tags option-legs
// @Produce json
// @Param strategy_id query string false "Strategy ID"
// @Param select_index query string false "Index/Underlying"
// @Param segment query string false "Segment" Enums(F&O, EQ)
// @Param option_type query string false "Option type" Enums(CALL, PUT)
// @Param action_type query string false "Action type" Enums(BUY, SELL)
// @Param expiry_type query string false "Expiry type" Enums(Weekly, Monthly, Quarterly)
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/option-legs/search [get]
func (h *OptionLegHandler) SearchOptionLegs(c *gin.Context) {
	filter := &dto.OptionLegFilter{
		Page:    1,
		PerPage: 20,
	}

	// Parse query parameters
	if strategyIDParam := c.Query("strategy_id"); strategyIDParam != "" {
		strategyID, err := uuid.Parse(strategyIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy ID"})
			return
		}
		filter.StrategyID = &strategyID
	}

	if selectIndex := c.Query("select_index"); selectIndex != "" {
		filter.SelectIndex = &selectIndex
	}

	if segmentParam := c.Query("segment"); segmentParam != "" {
		segment := model.Segment(segmentParam)
		filter.Segment = &segment
	}

	if optionTypeParam := c.Query("option_type"); optionTypeParam != "" {
		optionType := model.OptionType(optionTypeParam)
		filter.OptionType = &optionType
	}

	if actionTypeParam := c.Query("action_type"); actionTypeParam != "" {
		actionType := model.TransactionType(actionTypeParam)
		filter.ActionType = &actionType
	}

	if expiryTypeParam := c.Query("expiry_type"); expiryTypeParam != "" {
		expiryType := model.ExpiryType(expiryTypeParam)
		filter.ExpiryType = &expiryType
	}

	if pageParam := c.Query("page"); pageParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
		filter.Page = page
	}

	if perPageParam := c.Query("per_page"); perPageParam != "" {
		perPage, err := strconv.Atoi(perPageParam)
		if err != nil || perPage < 1 || perPage > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid per_page value (1-100)"})
			return
		}
		filter.PerPage = perPage
	}

	responses, total, err := h.optionLegService.SearchOptionLegs(c.Request.Context(), filter)
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
