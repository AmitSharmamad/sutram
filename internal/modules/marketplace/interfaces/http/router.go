package http

import (
	"algo/internal/infra/http/middleware"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the marketplace module routes
func RegisterRoutes(router *gin.RouterGroup, controller *MarketplaceController, userController *userHttp.UserController) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	// Public routes (no authentication required)
	router.GET("/marketplace/listings", controller.GetListings)
	router.GET("/marketplace/listings/:listing_id", controller.GetListingByID)
	router.GET("/marketplace/categories", controller.GetCategories)
	router.GET("/marketplace/listings/:listing_id/reviews", controller.GetReviews)

	// Authenticated routes
	marketplaceAuth := router.Group("/marketplace")
	marketplaceAuth.Use(middleware.AuthMiddleware(authServiceAdapter))
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
