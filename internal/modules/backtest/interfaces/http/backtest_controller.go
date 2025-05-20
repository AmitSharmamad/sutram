package http

import (
	"algo/internal/modules/backtest/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BacktestController handles HTTP requests for backtests
type BacktestController struct {
	backtestService service.BacktestService
	engineService   service.EngineService
}

// NewBacktestController creates a new backtest controller
func NewBacktestController(backtestService service.BacktestService, engineService service.EngineService) *BacktestController {
	return &BacktestController{
		backtestService: backtestService,
		engineService:   engineService,
	}
}

// HealthCheck handles the health check endpoint
func (c *BacktestController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "Backtest module is running",
	})
}

// RunBacktest handles the run backtest endpoint
func (c *BacktestController) RunBacktest(ctx *gin.Context) {
	// Get strategy ID from query parameter
	strategyIDStr := ctx.Query("strategy_id")
	if strategyIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "strategy_id query parameter is required",
		})
		return
	}

	// Parse the strategy ID
	strategyID, err := uuid.Parse(strategyIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid strategy ID format",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Call the engine service to start the backtest
	response, err := c.engineService.Start(ctx, strategyID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response with the run group ID
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Backtest started successfully",
		"backtest_details": response,
	})
}
