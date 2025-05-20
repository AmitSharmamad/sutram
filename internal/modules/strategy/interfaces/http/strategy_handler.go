package http

import (
	"net/http"

	"algo/internal/infra/http/middleware"
	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/application/service"
	userHttp "algo/internal/modules/user/interfaces/http"

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
// All routes are protected by authentication middleware
func (h *StrategyHandler) RegisterRoutes(r *gin.RouterGroup, userController *userHttp.UserController) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	strategies := r.Group("/strategy")
	// Apply authentication middleware to all strategy endpoints
	strategies.Use(middleware.AuthMiddleware(authServiceAdapter))
	{
		// Handle both with and without trailing slash for POST and GET
		strategies.POST("", h.CreateStrategy)
		strategies.POST("/", h.CreateStrategy)
		strategies.GET("", h.GetUserStrategies)
		strategies.GET("/", h.GetUserStrategies)

		// Other routes remain the same
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse request body
	var req dto.CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create strategy
	strategy, err := h.strategyService.CreateStrategy(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid request: " + err.Error(),
			})
		case service.ErrInvalidExpression:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid expression: " + err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return created strategy
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Strategy created successfully",
		"data":    strategy,
	})
}

// GetStrategyByID handles the retrieval of a strategy by ID
func (h *StrategyHandler) GetStrategyByID(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Get strategy by ID
	strategy, err := h.strategyService.GetStrategyByID(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Unauthorized access to strategy",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return strategy
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategy,
	})
}

// GetUserStrategies handles the retrieval of all strategies for a user
func (h *StrategyHandler) GetUserStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Get strategies for user
	strategies, err := h.strategyService.GetUserStrategies(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error: " + err.Error(),
		})
		return
	}

	// Ensure we return an empty array instead of null if strategies is nil
	if strategies == nil {
		strategies = []*dto.StrategyResponse{}
	}

	// Return strategies
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategies,
	})
}

// UpdateStrategy handles the update of an existing strategy
func (h *StrategyHandler) UpdateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Parse request body
	var req dto.UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Ensure ID in body matches URL
	req.ID = id

	// Update strategy
	strategy, err := h.strategyService.UpdateStrategy(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid request: " + err.Error(),
			})
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Unauthorized access to strategy",
			})
		case service.ErrInvalidExpression:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid expression: " + err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return updated strategy
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Strategy updated successfully",
		"data":    strategy,
	})
}

// DeleteStrategy handles the deletion of a strategy
func (h *StrategyHandler) DeleteStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Delete strategy
	err = h.strategyService.DeleteStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Unauthorized access to strategy",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Strategy deleted successfully",
	})
}

// ActivateStrategy handles the activation of a strategy
func (h *StrategyHandler) ActivateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Activate strategy
	err = h.strategyService.ActivateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Unauthorized access to strategy",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Strategy activated successfully",
	})
}

// DeactivateStrategy handles the deactivation of a strategy
func (h *StrategyHandler) DeactivateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Deactivate strategy
	err = h.strategyService.DeactivateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Unauthorized access to strategy",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Strategy deactivated successfully",
	})
}

// DuplicateStrategy handles the duplication of a strategy
func (h *StrategyHandler) DuplicateStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Duplicate strategy
	strategy, err := h.strategyService.DuplicateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Unauthorized access to strategy",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error: " + err.Error(),
			})
		}
		return
	}

	// Return duplicated strategy
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Strategy cloned successfully",
		"data":    strategy,
	})
}

// SearchStrategies handles the search of strategies with filter criteria
func (h *StrategyHandler) SearchStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse request body
	var filter dto.StrategyFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Search strategies
	strategies, err := h.strategyService.SearchStrategies(c.Request.Context(), &filter, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error: " + err.Error(),
		})
		return
	}

	// Ensure we return an empty array instead of null if strategies is nil
	if strategies == nil {
		strategies = []*dto.StrategyResponse{}
	}

	// Return strategies
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategies,
	})
}
