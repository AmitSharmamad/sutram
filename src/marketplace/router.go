package marketplace

import (
	"github.com/gin-gonic/gin"
	// Removing auth import for now as we're bypassing authentication
)

// RegisterRoutes registers the marketplace routes with the gin engine
func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	marketplace := r.Group("/marketplace")

	// Public routes
	marketplace.GET("/listings", controller.GetListings)
	marketplace.GET("/listings/:id", controller.GetListing)
	marketplace.GET("/categories", controller.GetCategories)
	marketplace.GET("/listings/:id/reviews", controller.GetReviews)

	// Routes that would normally require authentication
	// Using the same group for all routes to bypass authentication temporarily
	{
		// Listing operations
		marketplace.POST("/listings", controller.CreateListing)
		marketplace.PUT("/listings/:id", controller.UpdateListing)
		marketplace.DELETE("/listings/:id", controller.DeleteListing)
		marketplace.PUT("/listings/:id/performance", controller.UpdatePerformance)

		// Subscription operations
		marketplace.POST("/subscriptions", controller.CreateSubscription)
		marketplace.GET("/subscriptions/user", controller.GetUserSubscriptions)
		marketplace.PUT("/subscriptions/:id", controller.CancelSubscription)

		// Review operations
		marketplace.POST("/listings/:id/reviews", controller.CreateReview)
		marketplace.PUT("/reviews/:id", controller.UpdateReview)
		marketplace.DELETE("/reviews/:id", controller.DeleteReview)

		// Transaction operations
		marketplace.GET("/transactions", controller.GetTransactions)
		marketplace.GET("/transactions/:id", controller.GetTransaction)

		// Seller operations
		marketplace.GET("/seller/analytics", controller.GetSellerAnalytics)
		marketplace.GET("/seller/listings", controller.GetSellerListings)
	}
}
