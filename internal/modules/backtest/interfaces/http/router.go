package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the backtest module routes
func RegisterRoutes(router *gin.RouterGroup, controller *BacktestController) {
	// Create a backtest router group
	backtest := router.Group("/backtest")
	{
		// Health check endpoint
		backtest.GET("/health", controller.HealthCheck)

		// Note: As per requirements, we're not implementing CRUD APIs
		// This is just a placeholder for the backtest router
	}
}
