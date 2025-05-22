// CreateReview handles the request to create a new review
// @Summary Create a review
// @Description Create a new review for a marketplace listing
// @Tags Marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param listing_id path string true "Listing ID" format(uuid)
// @Param review body dto.CreateReviewRequest true "Review details"
// @Success 201 {object} dto.ApiResponse{data=dto.ReviewResponse} "Created review"
// @Failure 400 {object} dto.ApiResponse "Invalid request"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Listing not found"
// @Failure 409 {object} dto.ApiResponse "Review already exists"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /marketplace/listings/{listing_id}/reviews [post]
func (c *MarketplaceController) CreateReview(ctx *gin.Context) {
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

	// Parse listing ID from URL
	listingIDStr := ctx.Param("listing_id")
	listingID, err := uuid.Parse(listingIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID format",
		})
		return
	}

	// Parse request body
	var req dto.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create review
	review, err := c.reviewService.CreateReview(ctx.Request.Context(), listingID, userID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Failed to create review: " + err.Error()

		switch {
		case errors.Is(err, model.ErrListingNotFound):
			statusCode = http.StatusNotFound
			errorMsg = "Listing not found"
		case errors.Is(err, model.ErrReviewAlreadyExists):
			statusCode = http.StatusConflict
			errorMsg = "You have already reviewed this listing"
		case errors.Is(err, model.ErrInvalidRating):
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid rating value (must be between 1 and 5)"
		}

		ctx.JSON(statusCode, dto.ApiResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	// Return success response with created review
	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Data:    review,
	})
}