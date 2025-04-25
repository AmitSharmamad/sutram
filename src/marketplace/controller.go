package marketplace

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"algo/src/shared"
	"algo/src/web"
)

// Controller handles marketplace HTTP requests
type Controller struct {
	service Service
}

// NewController creates a new marketplace controller
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

// GetListings handles the request to get marketplace listings
func (c *Controller) GetListings(ctx *gin.Context) {
	// Parse query parameters
	category := ctx.Query("category")
	sort := ctx.Query("sort")
	search := ctx.Query("search")

	// Parse pagination parameters
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Create query params
	params := ListingQueryParams{
		Category: category,
		Sort:     sort,
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	// Get listings from service
	listings, err := c.service.GetListings(ctx, params)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve listings", "MARKETPLACE_LISTINGS_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, listings)
}

// GetListing handles the request to get a specific marketplace listing
func (c *Controller) GetListing(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid listing ID", "INVALID_ID")
		return
	}

	listing, err := c.service.GetListingByID(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("listing not found")) {
			web.BadReqResp(ctx, "Listing not found", "LISTING_NOT_FOUND")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to retrieve listing", "LISTING_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, listing)
}

// CreateListing handles the request to create a new marketplace listing
func (c *Controller) CreateListing(ctx *gin.Context) {
	var req shared.CreateListingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	listing, err := c.service.CreateListing(ctx, userUUID, &req)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to create listing", "CREATE_LISTING_ERROR")
		return
	}

	ctx.JSON(http.StatusCreated, listing)
}

// UpdateListing handles the request to update an existing marketplace listing
func (c *Controller) UpdateListing(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid listing ID", "INVALID_ID")
		return
	}

	var req shared.UpdateListingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	listing, err := c.service.UpdateListing(ctx, id, userUUID, &req)
	if err != nil {
		if errors.Is(err, errors.New("access denied: you are not the owner of this listing")) {
			web.ForbiddenResp(ctx, "Access denied", "ACCESS_DENIED")
		} else if errors.Is(err, errors.New("listing not found")) {
			web.BadReqResp(ctx, "Listing not found", "LISTING_NOT_FOUND")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to update listing", "UPDATE_LISTING_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, listing)
}

// DeleteListing handles the request to delete a marketplace listing
func (c *Controller) DeleteListing(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid listing ID", "INVALID_ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	err = c.service.DeleteListing(ctx, id, userUUID)
	if err != nil {
		if errors.Is(err, errors.New("access denied: you are not the owner of this listing")) {
			web.ForbiddenResp(ctx, "Access denied", "ACCESS_DENIED")
		} else if errors.Is(err, errors.New("listing not found")) {
			web.BadReqResp(ctx, "Listing not found", "LISTING_NOT_FOUND")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to delete listing", "DELETE_LISTING_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

// GetCategories handles the request to get all marketplace categories
func (c *Controller) GetCategories(ctx *gin.Context) {
	categories, err := c.service.GetCategories(ctx)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve categories", "CATEGORIES_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// UpdatePerformance handles the request to update performance metrics for a listing
func (c *Controller) UpdatePerformance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid listing ID", "INVALID_ID")
		return
	}

	var req shared.UpdatePerformanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	err = c.service.UpdatePerformance(ctx, id, userUUID, &req)
	if err != nil {
		if errors.Is(err, errors.New("access denied: you are not the owner of this listing")) {
			web.ForbiddenResp(ctx, "Access denied", "ACCESS_DENIED")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to update performance metrics", "PERFORMANCE_UPDATE_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Performance metrics updated successfully"})
}

// CreateSubscription handles the request to subscribe to a marketplace listing
func (c *Controller) CreateSubscription(ctx *gin.Context) {
	var req shared.CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	subscription, err := c.service.CreateSubscription(ctx, userUUID, &req)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to create subscription", "SUBSCRIPTION_ERROR")
		return
	}

	ctx.JSON(http.StatusCreated, subscription)
}

// GetUserSubscriptions handles the request to get a user's subscriptions
func (c *Controller) GetUserSubscriptions(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	subscriptions, err := c.service.GetUserSubscriptions(ctx, userUUID)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve subscriptions", "SUBSCRIPTIONS_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}

// CancelSubscription handles the request to cancel a subscription
func (c *Controller) CancelSubscription(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid subscription ID", "INVALID_ID")
		return
	}

	var req shared.CancelSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	err = c.service.CancelSubscription(ctx, id, userUUID, &req)
	if err != nil {
		if errors.Is(err, errors.New("subscription not found or access denied")) {
			web.ForbiddenResp(ctx, "Access denied", "ACCESS_DENIED")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to cancel subscription", "CANCEL_SUBSCRIPTION_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

// GetReviews handles the request to get reviews for a listing
func (c *Controller) GetReviews(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid listing ID", "INVALID_ID")
		return
	}

	reviews, err := c.service.GetReviews(ctx, id)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve reviews", "REVIEWS_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

// CreateReview handles the request to create a review for a listing
func (c *Controller) CreateReview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid listing ID", "INVALID_ID")
		return
	}

	var req shared.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	review, err := c.service.CreateReview(ctx, id, userUUID, &req)
	if err != nil {
		if errors.Is(err, errors.New("user already submitted a review for this listing")) {
			web.BadReqResp(ctx, "User already submitted a review", "DUPLICATE_REVIEW")
		} else if errors.Is(err, errors.New("you must be a subscriber to leave a review")) {
			web.ForbiddenResp(ctx, "Must be a subscriber to leave a review", "NOT_SUBSCRIBER")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to create review", "CREATE_REVIEW_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusCreated, review)
}

// UpdateReview handles the request to update a review
func (c *Controller) UpdateReview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid review ID", "INVALID_ID")
		return
	}

	var req shared.UpdateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadReqResp(ctx, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	review, err := c.service.UpdateReview(ctx, id, userUUID, &req)
	if err != nil {
		if errors.Is(err, errors.New("review not found or access denied")) {
			web.ForbiddenResp(ctx, "Review not found or access denied", "ACCESS_DENIED")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to update review", "UPDATE_REVIEW_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, review)
}

// DeleteReview handles the request to delete a review
func (c *Controller) DeleteReview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid review ID", "INVALID_ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	err = c.service.DeleteReview(ctx, id, userUUID)
	if err != nil {
		if errors.Is(err, errors.New("review not found or access denied")) {
			web.ForbiddenResp(ctx, "Review not found or access denied", "ACCESS_DENIED")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to delete review", "DELETE_REVIEW_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// GetTransactions handles the request to get a user's transaction history
func (c *Controller) GetTransactions(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	transactions, err := c.service.GetTransactions(ctx, userUUID)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve transactions", "TRANSACTIONS_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// GetTransaction handles the request to get a specific transaction
func (c *Controller) GetTransaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.BadReqResp(ctx, "Invalid transaction ID", "INVALID_ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	transaction, err := c.service.GetTransactionByID(ctx, id, userUUID)
	if err != nil {
		if errors.Is(err, errors.New("transaction not found or access denied")) {
			web.ForbiddenResp(ctx, "Transaction not found or access denied", "ACCESS_DENIED")
		} else {
			web.InternalServerErrorResp(ctx, "Failed to retrieve transaction", "TRANSACTION_ERROR")
		}
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

// GetSellerAnalytics handles the request to get analytics for a seller
func (c *Controller) GetSellerAnalytics(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	analytics, err := c.service.GetSellerAnalytics(ctx, userUUID)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve seller analytics", "ANALYTICS_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, analytics)
}

// GetSellerListings handles the request to get listings created by a seller
func (c *Controller) GetSellerListings(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		web.UnauthorizedResp(ctx, "Unauthorized", "UNAUTHORIZED")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		web.InternalServerErrorResp(ctx, "Invalid user ID format", "INVALID_USER_ID")
		return
	}

	listings, err := c.service.GetSellerListings(ctx, userUUID)
	if err != nil {
		web.InternalServerErrorResp(ctx, "Failed to retrieve seller listings", "LISTINGS_ERROR")
		return
	}

	ctx.JSON(http.StatusOK, listings)
}
