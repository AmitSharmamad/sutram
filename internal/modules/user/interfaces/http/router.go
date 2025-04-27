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
	auth := router.Group("/user")
	auth.Use(AuthMiddleware(controller))
	{
		auth.POST("/refresh-token", controller.RefreshToken)
		auth.GET("/profile", controller.GetProfile)
		auth.PUT("/profile", controller.UpdateProfile)
		auth.POST("/change-password", controller.ChangePassword)

		// Get user roles (can be own roles or admin access)
		auth.GET("/:id/roles", controller.GetUserRoles)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(AdminMiddleware(controller))
	{
		admin.GET("/:id", controller.GetUserByID)
		admin.POST("/:id/roles", controller.AssignRole)
		admin.DELETE("/:id/roles", controller.RemoveRole)
		admin.GET("/roles", controller.GetAllRoles)
		admin.PUT("/:id/activate", controller.ActivateUser)
		admin.PUT("/:id/deactivate", controller.DeactivateUser)
	}
}

// AuthMiddleware implements authentication middleware
// This middleware validates the JWT token and sets the userID in the context
// It can be reused across different modules to ensure authentication
func AuthMiddleware(controller *UserController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header
		tokenString := ExtractToken(ctx)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "No token provided",
			})
			return
		}

		// Validate the token
		claims, err := controller.authService.ValidateToken(ctx.Request.Context(), tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		// Set user ID in context for other handlers
		userID, err := controller.GetUserIDFromToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "Invalid token: user ID extraction failed",
			})
			return
		}

		// Set userID in context for use in handler functions
		ctx.Set("userID", userID)
		ctx.Set("claims", claims)
		// Add the token to context so it can be used by other handlers if needed
		ctx.Set("token", tokenString)

		ctx.Next()
	}
}

// AdminMiddleware implements admin role middleware
func AdminMiddleware(controller *UserController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header
		tokenString := ExtractToken(ctx)
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
