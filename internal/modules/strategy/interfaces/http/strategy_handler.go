package http

import (
	"net/http"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StrategyHandler handles HTTP requests for strategies
type StrategyHandler struct {
	strategyService service.StrategyService
}

// NewStrategyHandler creates a new strategy handler
func NewStrategyHandler(strategyService service.StrategyService) *StrategyHandler {
	return &StrategyHandler{
		strategyService: strategyService,
	}
}

// RegisterRoutes registers the routes for the strategy handler
func (h *StrategyHandler) RegisterRoutes(r *gin.RouterGroup) {
	strategies := r.Group("/strategy")
	{
		strategies.POST("/", h.CreateStrategy)
		strategies.GET("/", h.GetUserStrategies)
		strategies.GET("/:id", h.GetStrategyByID)
		strategies.PUT("/:id", h.UpdateStrategy)
		strategies.DELETE("/:id", h.DeleteStrategy)
		strategies.POST("/:id/activate", h.ActivateStrategy)
		strategies.POST("/:id/deactivate", h.DeactivateStrategy)
		strategies.POST("/:id/clone", h.DuplicateStrategy)
		strategies.POST("/search", h.SearchStrategies)
	}
}

// CreateStrategy handles the creation of a new strategy
func (h *StrategyHandler) CreateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req dto.CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Create strategy
	strategy, err := h.strategyService.CreateStrategy(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		case service.ErrInvalidExpression:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expression: " + err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return created strategy
	c.JSON(http.StatusCreated, strategy)
}

// GetStrategyByID handles the retrieval of a strategy by ID
func (h *StrategyHandler) GetStrategyByID(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get strategy by ID
	strategy, err := h.strategyService.GetStrategyByID(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to strategy"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return strategy
	c.JSON(http.StatusOK, strategy)
}

// GetUserStrategies handles the retrieval of all strategies for a user
func (h *StrategyHandler) GetUserStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get strategies for user
	strategies, err := h.strategyService.GetUserStrategies(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	// Return strategies
	c.JSON(http.StatusOK, strategies)
}

// UpdateStrategy handles the update of an existing strategy
func (h *StrategyHandler) UpdateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Parse request body
	var req dto.UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Ensure ID in body matches URL
	req.ID = id

	// Update strategy
	strategy, err := h.strategyService.UpdateStrategy(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to strategy"})
		case service.ErrInvalidExpression:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expression: " + err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return updated strategy
	c.JSON(http.StatusOK, strategy)
}

// DeleteStrategy handles the deletion of a strategy
func (h *StrategyHandler) DeleteStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Delete strategy
	err = h.strategyService.DeleteStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to strategy"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return success
	c.Status(http.StatusNoContent)
}

// ActivateStrategy handles the activation of a strategy
func (h *StrategyHandler) ActivateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Activate strategy
	err = h.strategyService.ActivateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to strategy"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return success
	c.Status(http.StatusOK)
}

// DeactivateStrategy handles the deactivation of a strategy
func (h *StrategyHandler) DeactivateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Deactivate strategy
	err = h.strategyService.DeactivateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to strategy"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return success
	c.Status(http.StatusOK)
}

// DuplicateStrategy handles the duplication of a strategy
func (h *StrategyHandler) DuplicateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Duplicate strategy
	strategy, err := h.strategyService.DuplicateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to strategy"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	// Return duplicated strategy
	c.JSON(http.StatusCreated, strategy)
}

// SearchStrategies handles the search of strategies with filter criteria
func (h *StrategyHandler) SearchStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var filter dto.StrategyFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Search strategies
	strategies, err := h.strategyService.SearchStrategies(c.Request.Context(), &filter, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	// Return strategies
	c.JSON(http.StatusOK, strategies)
}
