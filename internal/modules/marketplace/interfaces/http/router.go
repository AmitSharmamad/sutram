// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/interfaces/http/router.go
package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the marketplace module routes
func RegisterRoutes(router *gin.RouterGroup, controller *MarketplaceController) {
	// Public routes (no authentication required)
	router.GET("/marketplace/listings", controller.GetListings)
	router.GET("/marketplace/listings/:listing_id", controller.GetListingByID)
	router.GET("/marketplace/categories", controller.GetCategories)
	router.GET("/marketplace/listings/:listing_id/reviews", controller.GetReviews)

	// Authenticated routes
	marketplaceAuth := router.Group("/marketplace")
	marketplaceAuth.Use(AuthMiddleware())
	{
		// Listing management
		marketplaceAuth.POST("/listings", controller.CreateListing)
		marketplaceAuth.PUT("/listings/:listing_id", controller.UpdateListing)
		marketplaceAuth.DELETE("/listings/:listing_id", controller.DeleteListing)
		marketplaceAuth.PUT("/listings/:listing_id/performance", controller.UpdatePerformanceMetrics)

		// Subscriptions
		marketplaceAuth.POST("/subscriptions", controller.CreateSubscription)
		marketplaceAuth.GET("/subscriptions/user", controller.GetUserSubscriptions)
		marketplaceAuth.PUT("/subscriptions/:subscription_id", controller.CancelSubscription)

		// Reviews
		marketplaceAuth.POST("/listings/:listing_id/reviews", controller.CreateReview)
		marketplaceAuth.PUT("/reviews/:review_id", controller.UpdateReview)
		marketplaceAuth.DELETE("/reviews/:review_id", controller.DeleteReview)

		// Transactions
		marketplaceAuth.GET("/transactions", controller.GetTransactions)
		marketplaceAuth.GET("/transactions/:transaction_id", controller.GetTransactionByID)

		// Seller endpoints
		marketplaceAuth.GET("/seller/listings", controller.GetSellerListings)
		marketplaceAuth.GET("/seller/analytics", controller.GetSellerAnalytics)
	}
}

// AuthMiddleware implements authentication middleware
func AuthMiddleware() gin.HandlerFunc {
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

		// In a real implementation, validate the token and extract user ID and roles
		// For now, we're assuming this is handled by a common middleware
		// and we're just checking if a user ID exists in the context

		// If the user ID is not set in the context, abort
		if _, exists := ctx.Get("userID"); !exists {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			return
		}

		ctx.Next()
	}
}
