package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the user module routes
func RegisterRoutes(router *gin.RouterGroup, controller *UserController) {
	// Public routes (no authentication required)
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	// Authenticated routes
	auth := router.Group("/")
	auth.Use(AuthMiddleware(controller))
	{
		auth.POST("/refresh-token", controller.RefreshToken)
		auth.GET("/profile", controller.GetProfile)
		auth.PUT("/profile", controller.UpdateProfile)
		auth.POST("/change-password", controller.ChangePassword)

		// Get user roles (can be own roles or admin access)
		auth.GET("/users/:id/roles", controller.GetUserRoles)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(AdminMiddleware(controller))
	{
		admin.GET("/users/:id", controller.GetUserByID)
		admin.POST("/users/:id/roles", controller.AssignRole)
		admin.DELETE("/users/:id/roles", controller.RemoveRole)
		admin.GET("/roles", controller.GetAllRoles)
		admin.PUT("/users/:id/activate", controller.ActivateUser)
		admin.PUT("/users/:id/deactivate", controller.DeactivateUser)
	}
}

// AuthMiddleware implements authentication middleware
func AuthMiddleware(controller *UserController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header
		tokenString := extractToken(ctx)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "No token provided",
			})
			return
		}

		// Validate the token
		_, err := controller.authService.ValidateToken(ctx.Request.Context(), tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		ctx.Next()
	}
}

// AdminMiddleware implements admin role middleware
func AdminMiddleware(controller *UserController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header
		tokenString := extractToken(ctx)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "No token provided",
			})
			return
		}

		// Check if user has admin role
		hasAdminRole, err := controller.authService.HasRole(tokenString, "admin")
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		if !hasAdminRole {
			ctx.AbortWithStatusJSON(403, gin.H{
				"success": false,
				"error":   "Permission denied: Admin role required",
			})
			return
		}

		ctx.Next()
	}
}
