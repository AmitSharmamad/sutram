package http

import (
	"algo/internal/infra/http/middleware"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the backtest module routes
func RegisterRoutes(router *gin.RouterGroup, controller *BacktestController, userController *userHttp.UserController) {
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
		// Add authenticated endpoints here when needed
		// This is a placeholder for future authenticated endpoints
	}
}
