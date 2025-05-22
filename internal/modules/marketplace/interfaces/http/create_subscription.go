// CreateSubscription handles the request to create a new subscription
// @Summary Purchase/Subscribe to a strategy
// @Description Subscribe to a marketplace strategy listing
// @Tags Marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subscription body dto.CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} dto.ApiResponse{data=dto.SubscriptionResponse} "Created subscription"
// @Failure 400 {object} dto.ApiResponse "Invalid request"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Listing not found"
// @Failure 409 {object} dto.ApiResponse "Subscription already exists"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /marketplace/subscriptions [post]
func (c *MarketplaceController) CreateSubscription(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}
	
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req dto.CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create subscription
	subscription, err := c.subscriptionService.CreateSubscription(ctx.Request.Context(), userID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Failed to create subscription: " + err.Error()

		switch {
		case errors.Is(err, model.ErrListingNotFound):
			statusCode = http.StatusNotFound
			errorMsg = "Listing not found"
		case errors.Is(err, model.ErrSubscriptionAlreadyExists):
			statusCode = http.StatusConflict
			errorMsg = "You are already subscribed to this listing"
		case errors.Is(err, model.ErrInvalidPaymentMethod):
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid payment method"
		}

		ctx.JSON(statusCode, dto.ApiResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	// Create transaction record
	_, err = c.transactionService.CreateTransaction(ctx.Request.Context(), userID, req.ListingID, req.PaymentMethod, subscription.ID)
	if err != nil {
		// Log error but don't fail the request
		// TODO: Implement proper transaction rollback
		// For now, we'll just log the error
		// logger.Error("Failed to create transaction record", err)
	}

	// Return success response with created subscription
	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Data:    subscription,
	})
}