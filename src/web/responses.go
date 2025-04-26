package web

import (
	"github.com/gin-gonic/gin"
)

// RespondError sends an error response with the given status code, message, and error
func RespondError(ctx *gin.Context, statusCode int, message string, err error) {
	errorResponse := gin.H{
		"success": false,
		"message": message,
	}

	// Include error details if available
	if err != nil {
		errorResponse["error"] = err.Error()
	}

	ctx.JSON(statusCode, errorResponse)
}

// RespondSuccess sends a success response with the given status code and data
func RespondSuccess(ctx *gin.Context, statusCode int, data gin.H) {
	// Add success flag to the response
	data["success"] = true
	ctx.JSON(statusCode, data)
}
