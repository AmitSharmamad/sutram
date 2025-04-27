package http

import (
	"algo/internal/modules/backtest/application/service"

	"github.com/gin-gonic/gin"
)

// BacktestController handles HTTP requests for backtests
type BacktestController struct {
	backtestService service.BacktestService
}

// NewBacktestController creates a new backtest controller
func NewBacktestController(backtestService service.BacktestService) *BacktestController {
	return &BacktestController{
		backtestService: backtestService,
	}
}

// HealthCheck handles the health check endpoint
func (c *BacktestController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "Backtest module is running",
	})
}
