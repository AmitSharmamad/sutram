package shared

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContext extracts the user ID from the gin context
func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	// Handle different possible types
	switch id := userID.(type) {
	case uuid.UUID:
		return id, nil
	case string:
		// Try to parse the string as UUID
		return uuid.Parse(id)
	default:
		return uuid.Nil, errors.New("invalid user ID format in context")
	}
}

// ParseCommaSeparatedList splits a comma-separated string into a slice of strings
func ParseCommaSeparatedList(input string) []string {
	if input == "" {
		return []string{}
	}

	// Split by comma and trim spaces
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
