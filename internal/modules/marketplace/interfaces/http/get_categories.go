// GetCategories handles the request to get all marketplace categories
// @Summary Get marketplace categories
// @Description Retrieve all available marketplace categories
// @Tags Marketplace
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]dto.CategoryResponse} "List of marketplace categories"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /marketplace/categories [get]
func (c *MarketplaceController) GetCategories(ctx *gin.Context) {
	// Get all categories
	categories, err := c.categoryService.GetAllCategories(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to fetch categories: " + err.Error(),
		})
		return
	}

	// Return success response with categories
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    categories,
	})
}