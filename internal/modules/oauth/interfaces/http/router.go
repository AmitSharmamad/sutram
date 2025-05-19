package http

import (
	"algo/internal/infra/http/middleware"

	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the OAuth module routes
func RegisterRoutes(router *gin.RouterGroup, controller *OAuthController, userController *userHttp.UserController) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	// Public routes (no authentication required)
	auth := router.Group("/auth")
	{
		// OAuth routes
		oauth := auth.Group("/oauth")
		{
			// Get available OAuth providers
			oauth.GET("/providers", controller.GetAvailableProviders)

			// Provider-specific routes
			provider := oauth.Group("/:provider")
			{
				// Initiate OAuth login
				provider.GET("/login", controller.InitiateLogin)

				// Handle OAuth callback
				provider.GET("/callback", controller.HandleCallback)
			}
		}

		// Google-specific route for backward compatibility
		auth.GET("/google/login", func(ctx *gin.Context) {
			ctx.Params = append(ctx.Params, gin.Param{Key: "provider", Value: "google"})
			controller.InitiateLogin(ctx)
		})
	}

	// Authenticated routes
	authUser := router.Group("/user/oauth")
	authUser.Use(middleware.AuthMiddleware(authServiceAdapter))
	{
		// Get user OAuth connections
		authUser.GET("/connections", controller.GetUserConnections)

		// Disconnect provider
		authUser.DELETE("/:provider", controller.DisconnectProvider)
	}
}
