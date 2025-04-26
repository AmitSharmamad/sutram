package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"algo/src/shared"
)

// Common errors
var (
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden access")
)

// AuthMiddleware ensures the request is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := shared.GetUserIDFromContext(c)
		if err != nil {
			RespondError(c, http.StatusUnauthorized, "Authentication required", ErrUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RoleRequiredMiddleware ensures the authenticated user has the required role
func RoleRequiredMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure the user is authenticated
		_, err := shared.GetUserIDFromContext(c)
		if err != nil {
			RespondError(c, http.StatusUnauthorized, "Authentication required", ErrUnauthorized)
			c.Abort()
			return
		}

		// Get user role from context or perform a role check
		// Note: This implementation assumes you have a way to get the user's role
		userRole, exists := c.Get("userRole")
		if !exists {
			// If role is not in context, you might need to retrieve it from a service/database
			// For example: userRole = userService.GetUserRole(userID)
			// This is a placeholder - replace with actual implementation
			userRole = "user" // Default role
		}

		// Check if the user has the required role
		if userRole != requiredRole {
			RespondError(c, http.StatusForbidden, "Insufficient permissions", ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
