package http

import (
	"algo/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the user module routes
func RegisterRoutes(router *gin.RouterGroup, controller *UserController) {
	// Create auth service adapter
	authServiceAdapter := NewAuthServiceAdapter(controller.authService)

	// Public routes (no authentication required)
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.POST("/logout", controller.Logout)         // Add logout endpoint to clear token cookie
	router.GET("/user/me", controller.CheckMe)        // Public endpoint to check login status
	router.POST("/waitlist", controller.JoinWaitlist) // Public endpoint to join waitlist

	// Authenticated routes
	auth := router.Group("/user")
	auth.Use(middleware.AuthMiddleware(authServiceAdapter))
	{
		auth.POST("/refresh-token", controller.RefreshToken)
		auth.GET("/profile", controller.GetProfile)
		auth.PUT("/profile", controller.UpdateProfile)
		auth.POST("/change-password", controller.ChangePassword)
		auth.GET("/dashboard/statistics", controller.GetDashboardStatistics)

		// Get user roles (can be own roles or admin access)
		auth.GET("/:id/roles", controller.GetUserRoles)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.AdminMiddleware(authServiceAdapter))
	{
		admin.GET("/:id", controller.GetUserByID)
		admin.POST("/:id/roles", controller.AssignRole)
		admin.DELETE("/:id/roles", controller.RemoveRole)
		admin.GET("/roles", controller.GetAllRoles)
		admin.PUT("/:id/activate", controller.ActivateUser)
		admin.PUT("/:id/deactivate", controller.DeactivateUser)
		admin.GET("/waitlist", controller.GetAllWaitlistEntries)
		admin.GET("/waitlist/count", controller.GetWaitlistCount)
	}
}
