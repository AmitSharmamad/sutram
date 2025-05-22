// CloneDefaultStrategy handles the cloning of a default strategy to a user's account
// @Summary Clone a default strategy
// @Description Clone a default/system strategy to the user's account
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Default Strategy ID to clone" format(uuid)
// @Success 201 {object} dto.StrategyResponse "Cloned strategy"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Strategy not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/default/{id}/clone [post]
func (h *StrategyHandler) CloneDefaultStrategy(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	// Clone default strategy
	clonedStrategy, err := h.strategyService.CloneDefaultStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Default strategy not found",
			})
		case service.ErrNotDefaultStrategy:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Specified strategy is not a default strategy",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to clone default strategy: " + err.Error(),
			})
		}
		return
	}

	// Return cloned strategy
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    clonedStrategy,
	})
}