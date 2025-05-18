// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/interfaces/http/utils.go
package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExtractToken extracts the token from the Authorization header
func ExtractToken(ctx *gin.Context) string {
	bearerToken := ctx.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	// Get the user ID from the context (set by middleware)
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		return uuid.Nil, false
	}

	// Convert to UUID
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		// Try to convert from string
		userIDStr, ok := userIDVal.(string)
		if !ok {
			return uuid.Nil, false
		}

		// Parse UUID
		parsedID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, false
		}

		return parsedID, true
	}

	return userID, true
}
