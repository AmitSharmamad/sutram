package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthService defines the interface for authentication service
// This allows us to decouple the middleware from the specific auth implementation
type AuthService interface {
	ValidateToken(ctx interface{}, tokenString string) (interface{}, error)
	GetUserIDFromToken(tokenString string) (interface{}, error)
	HasRole(tokenString string, role string) (bool, error)
}

// ExtractToken extracts the JWT token from the request
// Looks for token in Authorization header (Bearer token) or in cookie
func ExtractToken(ctx *gin.Context) string {
	// First, try to get token from Authorization header
	bearerToken := ctx.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(strings.ToUpper(bearerToken), "BEARER ") {
		return bearerToken[7:]
	}

	// If not in header, try to get from cookie
	cookie, err := ctx.Request.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

// AuthMiddleware implements authentication middleware
// This middleware validates the JWT token and sets the userID in the context
// It can be reused across different modules to ensure authentication
func AuthMiddleware(authService AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header or cookie
		tokenString := ExtractToken(ctx)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "No token provided",
			})
			return
		}

		// Validate the token - pass the request context
		claims, err := authService.ValidateToken(ctx.Request.Context(), tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		// Set user ID in context for other handlers
		userID, err := authService.GetUserIDFromToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: user ID extraction failed",
			})
			return
		}

		// Set userID in context for use in handler functions
		ctx.Set("userID", userID)
		ctx.Set("claims", claims)
		ctx.Set("token", tokenString) // Add the token to context for use by other handlers if needed

		ctx.Next()
	}
}

// AdminMiddleware implements admin role middleware
// This middleware checks if the user has the admin role
func AdminMiddleware(authService AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header or cookie
		tokenString := ExtractToken(ctx)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "No token provided",
			})
			return
		}

		// Check if user has admin role
		hasAdminRole, err := authService.HasRole(tokenString, "admin")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		if !hasAdminRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Permission denied: Admin role required",
			})
			return
		}

		ctx.Next()
	}
}
