// GetListingByID handles the request to get a specific marketplace listing
// @Summary Get a marketplace listing by ID
// @Description Retrieve a marketplace listing by its unique ID
// @Tags Marketplace
// @Accept json
// @Produce json
// @Param id path string true "Listing ID" format(uuid)
// @Success 200 {object} dto.ApiResponse{data=dto.ListingResponse} "Marketplace listing details"
// @Failure 400 {object} dto.ApiResponse "Invalid ID"
// @Failure 404 {object} dto.ApiResponse "Listing not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /marketplace/listings/{id} [get]
func (c *MarketplaceController) GetListingByID(ctx *gin.Context) {
	// Parse listing ID from URL
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid listing ID format",
		})
		return
	}

	// Get listing details
	listing, err := c.listingService.GetListingByID(ctx.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Failed to fetch listing: " + err.Error()

		if errors.Is(err, model.ErrListingNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "Listing not found"
		}

		ctx.JSON(statusCode, dto.ApiResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	// Return success response with listing details
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    listing,
	})
}