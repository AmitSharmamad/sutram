// GetListings handles the request to get all marketplace listings
// @Summary Get marketplace listings
// @Description Retrieve all marketplace listings with optional filtering
// @Tags Marketplace
// @Accept json
// @Produce json
// @Param category query string false "Filter by category ID"
// @Param sort query string false "Sort field (price, rating, popularity, newest)"
// @Param order query string false "Sort order (asc, desc)"
// @Param page query integer false "Page number (default: 1)"
// @Param limit query integer false "Results per page (default: 10)"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param min_rating query number false "Minimum rating filter"
// @Success 200 {object} dto.ApiResponse{data=[]dto.ListingResponse} "List of marketplace listings"
// @Failure 400 {object} dto.ApiResponse "Invalid parameters"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /marketplace/listings [get]
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
		filters.Order,
		filters.Page,
		filters.Limit,
		filters.MinPrice,
		filters.MaxPrice,
		filters.MinRating,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to fetch listings: " + err.Error(),
		})
		return
	}

	// Return success response with listings and pagination info
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: gin.H{
			"listings": listings,
			"total":    total,
			"page":     filters.Page,
			"limit":    filters.Limit,
			"pages":    (total + filters.Limit - 1) / filters.Limit,
		},
	})
}