// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/interfaces/http/marketplace_controller.go
package http

import (
	"errors"
	"net/http"

	"algo/internal/modules/marketplace/application/dto"
	"algo/internal/modules/marketplace/application/service"
	"algo/internal/modules/marketplace/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MarketplaceController handles marketplace HTTP requests
type MarketplaceController struct {
	listingService      *service.ListingService
	categoryService     *service.CategoryService
	subscriptionService *service.SubscriptionService
	reviewService       *service.ReviewService
	transactionService  *service.TransactionService
	sellerService       *service.SellerService
}

// NewMarketplaceController creates a new marketplace controller
func NewMarketplaceController(
	listingService *service.ListingService,
	categoryService *service.CategoryService,
	subscriptionService *service.SubscriptionService,
	reviewService *service.ReviewService,
	transactionService *service.TransactionService,
	sellerService *service.SellerService,
) *MarketplaceController {
	return &MarketplaceController{
		listingService:      listingService,
		categoryService:     categoryService,
		subscriptionService: subscriptionService,
		reviewService:       reviewService,
		transactionService:  transactionService,
		sellerService:       sellerService,
	}
}

// GetListings handles the request to get all marketplace listings
func (c *MarketplaceController) GetListings(ctx *gin.Context) {
	var filters dto.ListingFilterRequest
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	// Get listings with filters
	listings, total, err := c.listingService.GetListings(
		ctx.Request.Context(),
		filters.Category,
		filters.Sort,
		filters.Search,
		filters.PriceMin,
		filters.PriceMax,
		filters.Page,
		filters.PerPage,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get listings: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	listingResponses := make([]dto.ListingResponse, len(listings))
	for i, listing := range listings {
		listingResponses[i] = dto.ToListingResponse(listing)
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: dto.NewPaginatedResponse(
			listingResponses,
			total,
			filters.Page,
			filters.PerPage,
		),
	})
}

// GetListingByID handles the request to get a specific marketplace listing
func (c *MarketplaceController) GetListingByID(ctx *gin.Context) {
	idStr := ctx.Param("listing_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	// Get listing
	listing, err := c.listingService.GetListingByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get listing: " + err.Error(),
		})
		return
	}

	// Increment view count
	go c.listingService.IncrementViews(ctx.Request.Context(), id)

	// Convert to response DTO
	response := dto.ToListingResponse(listing)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    response,
	})
}

// CreateListing handles the request to create a new marketplace listing
func (c *MarketplaceController) CreateListing(ctx *gin.Context) {
	var req dto.CreateListingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Parse category IDs
	var categoryIDs []uuid.UUID
	for _, idStr := range req.CategoryIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
				Success: false,
				Error:   "Invalid category ID: " + idStr,
			})
			return
		}
		categoryIDs = append(categoryIDs, id)
	}

	// Create listing
	listing, err := c.listingService.CreateListing(
		ctx.Request.Context(),
		userID,
		req.StrategyID,
		req.Title,
		req.Description,
		req.Price,
		req.SubscriptionType,
		categoryIDs,
		req.Features,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create listing: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToListingResponse(listing)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Listing created successfully",
		Data:    response,
	})
}

// UpdateListing handles the request to update an existing marketplace listing
func (c *MarketplaceController) UpdateListing(ctx *gin.Context) {
	idStr := ctx.Param("listing_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	var req dto.UpdateListingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.SubscriptionType != "" {
		updates["subscription_type"] = req.SubscriptionType
	}
	if req.IsActive != nil {
		updates["is_active"] = req.IsActive
	}

	// Parse category IDs if provided
	if req.CategoryIDs != nil && len(req.CategoryIDs) > 0 {
		var categoryIDs []uuid.UUID
		for _, idStr := range req.CategoryIDs {
			catID, err := uuid.Parse(idStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
					Success: false,
					Error:   "Invalid category ID: " + idStr,
				})
				return
			}
			categoryIDs = append(categoryIDs, catID)
		}
		updates["category_ids"] = categoryIDs
	}

	// Add features if provided
	if req.Features != nil {
		updates["features"] = req.Features
	}

	// Update listing
	updatedListing, err := c.listingService.UpdateListing(
		ctx.Request.Context(),
		id,
		userID,
		updates,
	)

	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to update this listing",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update listing: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToListingResponse(updatedListing)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Listing updated successfully",
		Data:    response,
	})
}

// DeleteListing handles the request to delete a marketplace listing
func (c *MarketplaceController) DeleteListing(ctx *gin.Context) {
	idStr := ctx.Param("listing_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Delete listing
	err = c.listingService.DeleteListing(ctx.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to delete this listing",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete listing: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Listing deleted successfully",
	})
}

// UpdatePerformanceMetrics handles the request to update performance metrics for a listing
func (c *MarketplaceController) UpdatePerformanceMetrics(ctx *gin.Context) {
	idStr := ctx.Param("listing_id")
	listingID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	var req dto.UpdatePerformanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Create performance model
	performance := &model.Performance{
		ListingID:          listingID,
		SharpeRatio:        req.SharpeRatio,
		SortinoRatio:       req.SortinoRatio,
		TotalReturnPct:     req.TotalReturnPct,
		AnnualizedReturn:   req.AnnualizedReturn,
		MaxDrawdown:        req.MaxDrawdown,
		WinRate:            req.WinRate,
		ProfitFactor:       req.ProfitFactor,
		BacktestPeriodDays: req.BacktestPeriodDays,
	}

	// Update performance metrics
	err = c.listingService.UpdatePerformanceMetrics(
		ctx.Request.Context(),
		listingID,
		userID,
		performance,
	)

	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to update this listing's performance",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update performance metrics: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Performance metrics updated successfully",
		Data: dto.PerformanceResponse{
			SharpeRatio:        req.SharpeRatio,
			SortinoRatio:       req.SortinoRatio,
			TotalReturnPct:     req.TotalReturnPct,
			AnnualizedReturn:   req.AnnualizedReturn,
			MaxDrawdown:        req.MaxDrawdown,
			WinRate:            req.WinRate,
			ProfitFactor:       req.ProfitFactor,
			BacktestPeriodDays: req.BacktestPeriodDays,
		},
	})
}

// GetCategories handles the request to get all marketplace categories
func (c *MarketplaceController) GetCategories(ctx *gin.Context) {
	categories, err := c.categoryService.GetAllCategories(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get categories: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			Icon:        category.Icon,
		}
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    categoryResponses,
	})
}

// CreateSubscription handles the request to create a new subscription
func (c *MarketplaceController) CreateSubscription(ctx *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Create subscription
	subscription, err := c.subscriptionService.CreateSubscription(
		ctx.Request.Context(),
		userID,
		req.ListingID,
		req.PaymentMethod,
		req.PaymentDetails,
	)

	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		if errors.Is(err, service.ErrAlreadySubscribed) {
			ctx.JSON(http.StatusConflict, dto.ApiResponse{
				Success: false,
				Error:   "You are already subscribed to this strategy",
			})
			return
		}
		if errors.Is(err, service.ErrPaymentFailed) {
			ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
				Success: false,
				Error:   "Payment processing failed",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create subscription: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToSubscriptionResponse(subscription)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Subscription created successfully",
		Data:    response,
	})
}

// GetUserSubscriptions handles the request to get all subscriptions for the user
func (c *MarketplaceController) GetUserSubscriptions(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get subscriptions
	subscriptions, err := c.subscriptionService.GetUserSubscriptions(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get subscriptions: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	subscriptionResponses := make([]dto.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		subscriptionResponses[i] = dto.ToSubscriptionResponse(subscription)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    subscriptionResponses,
	})
}

// CancelSubscription handles the request to cancel a subscription
func (c *MarketplaceController) CancelSubscription(ctx *gin.Context) {
	idStr := ctx.Param("subscription_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid subscription ID",
		})
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Only allow cancellation
	if req.Status != "cancelled" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Only cancellation is supported",
		})
		return
	}

	// Cancel subscription
	subscription, err := c.subscriptionService.CancelSubscription(ctx.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Subscription not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to cancel this subscription",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to cancel subscription: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToSubscriptionResponse(subscription)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Subscription cancelled successfully",
		Data:    response,
	})
}

// GetReviews handles the request to get reviews for a listing
func (c *MarketplaceController) GetReviews(ctx *gin.Context) {
	idStr := ctx.Param("listing_id")
	listingID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	// Get reviews
	reviews, err := c.reviewService.GetReviewsByListingID(ctx.Request.Context(), listingID)
	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get reviews: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	reviewResponses := make([]dto.ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = dto.ToReviewResponse(review)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    reviewResponses,
	})
}

// CreateReview handles the request to create a new review
func (c *MarketplaceController) CreateReview(ctx *gin.Context) {
	idStr := ctx.Param("listing_id")
	listingID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	var req dto.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Create review
	review, err := c.reviewService.CreateReview(
		ctx.Request.Context(),
		userID,
		listingID,
		req.Rating,
		req.Comment,
	)

	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Listing not found",
			})
			return
		}
		if errors.Is(err, service.ErrAlreadyReviewed) {
			ctx.JSON(http.StatusConflict, dto.ApiResponse{
				Success: false,
				Error:   "You have already reviewed this listing",
			})
			return
		}
		if errors.Is(err, service.ErrReviewNotPermitted) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You must purchase this strategy before reviewing it",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create review: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToReviewResponse(review)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Review created successfully",
		Data:    response,
	})
}

// UpdateReview handles the request to update a review
func (c *MarketplaceController) UpdateReview(ctx *gin.Context) {
	idStr := ctx.Param("review_id")
	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid review ID",
		})
		return
	}

	var req dto.UpdateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Update review
	review, err := c.reviewService.UpdateReview(
		ctx.Request.Context(),
		reviewID,
		userID,
		req.Rating,
		req.Comment,
	)

	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Review not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to update this review",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update review: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToReviewResponse(review)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Review updated successfully",
		Data:    response,
	})
}

// DeleteReview handles the request to delete a review
func (c *MarketplaceController) DeleteReview(ctx *gin.Context) {
	idStr := ctx.Param("review_id")
	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid review ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Delete review
	err = c.reviewService.DeleteReview(ctx.Request.Context(), reviewID, userID)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Review not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to delete this review",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete review: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Review deleted successfully",
	})
}

// GetTransactions handles the request to get user transactions
func (c *MarketplaceController) GetTransactions(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get transactions
	transactions, err := c.transactionService.GetTransactionsByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get transactions: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	transactionResponses := make([]dto.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = dto.ToTransactionResponse(transaction)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    transactionResponses,
	})
}

// GetTransactionByID handles the request to get a specific transaction
func (c *MarketplaceController) GetTransactionByID(ctx *gin.Context) {
	idStr := ctx.Param("transaction_id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid transaction ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get transaction
	transaction, err := c.transactionService.GetTransactionByID(ctx.Request.Context(), transactionID, userID)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Transaction not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to view this transaction",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get transaction: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToTransactionResponse(transaction)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    response,
	})
}

// GetSellerListings handles the request to get all listings for the seller
func (c *MarketplaceController) GetSellerListings(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get listings
	listings, err := c.sellerService.GetSellerListings(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get seller listings: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	listingResponses := make([]dto.ListingResponse, len(listings))
	for i, listing := range listings {
		listingResponses[i] = dto.ToListingResponse(listing)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    listingResponses,
	})
}

// GetSellerAnalytics handles the request to get analytics data for the seller
func (c *MarketplaceController) GetSellerAnalytics(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get analytics
	analytics, err := c.sellerService.GetSellerAnalytics(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get seller analytics: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    analytics,
	})
}
