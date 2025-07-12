package http

import (
	"algo/internal/infra/http/middleware"
	"algo/internal/modules/user/application/service"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the backtest module routes
func RegisterRoutes(router *gin.RouterGroup, controller *BacktestController, userController *userHttp.UserController, usageLimitService service.UsageLimitService) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	// Create a backtest router group for public endpoints
	backtest := router.Group("/backtest")
	{
		// Health check endpoint
		backtest.GET("/health", controller.HealthCheck)
	}

	// Create a backtest router group for authenticated endpoints
	backtestAuth := router.Group("/backtest")
	backtestAuth.Use(middleware.AuthMiddleware(authServiceAdapter))
	{
		// Add run backtest endpoint with usage limit middleware
		backtestAuth.POST("/run", middleware.BacktestLimitMiddleware(usageLimitService), controller.RunBacktest)
		// Add Excel export endpoint
		backtestAuth.GET("/export", controller.ExportToExcel)
		// Add indicator export endpoint
		backtestAuth.GET("/export/indicators", controller.ExportIndicatorsToExcel)
		// Add new endpoint for OHLCV + indicators Excel export
		backtestAuth.GET("/export/ohlcv-indicators", controller.ExportOHLCVIndicatorsToExcel)
	}
}
