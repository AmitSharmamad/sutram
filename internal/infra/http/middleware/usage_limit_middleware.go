package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"algo/internal/domain/entity"
	"algo/internal/modules/user/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UsageLimitMiddleware creates a middleware that checks usage limits before allowing requests
func UsageLimitMiddleware(usageLimitService service.UsageLimitService, limitType entity.UsageLimitType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Check if user can perform the action
		canPerform, err := usageLimitService.CanPerformAction(c.Request.Context(), userID, limitType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to check usage limits: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !canPerform {
			// Get remaining usage for better error message
			remaining, _ := usageLimitService.GetRemainingUsage(c.Request.Context(), userID, limitType)

			var message string
			switch limitType {
			case entity.BacktestPerMonth:
				message = "Monthly backtest limit exceeded"
			case entity.DeploymentCount:
				message = "Deployment limit exceeded"
			case entity.LiveTradesPerDay:
				message = "Daily live trades limit exceeded"
			case entity.APICallsPerHour:
				message = "Hourly API calls limit exceeded"
			default:
				message = "Usage limit exceeded"
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":    false,
				"error":      message,
				"limit_type": string(limitType),
				"remaining":  remaining,
			})
			c.Abort()
			return
		}

		// Store the limit type and user ID in context for post-processing
		c.Set("usageLimitType", limitType)
		c.Set("usageLimitService", usageLimitService)

		c.Next()

		// After successful request completion, increment usage
		// Only increment if the response was successful (2xx status code)
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			metadata := map[string]interface{}{
				"endpoint":   c.Request.URL.Path,
				"method":     c.Request.Method,
				"user_agent": c.Request.UserAgent(),
				"timestamp":  time.Now().UTC(),
			}

			// Increment usage in a goroutine to not block the response
			go func() {
				ctx := context.Background()
				if err := usageLimitService.IncrementUsage(ctx, userID, limitType, metadata); err != nil {
					// Log the error but don't fail the request
					fmt.Printf("Failed to increment usage for user %s, limit type %s: %v\n", userID, limitType, err)
				}
			}()
		}
	}
}

// BacktestLimitMiddleware specifically for backtest endpoints
func BacktestLimitMiddleware(usageLimitService service.UsageLimitService) gin.HandlerFunc {
	return UsageLimitMiddleware(usageLimitService, entity.BacktestPerMonth)
}

// DeploymentLimitMiddleware specifically for deployment creation endpoints
func DeploymentLimitMiddleware(usageLimitService service.UsageLimitService) gin.HandlerFunc {
	return UsageLimitMiddleware(usageLimitService, entity.DeploymentCount)
}
