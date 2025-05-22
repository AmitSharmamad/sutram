// DuplicateStrategy handles the duplication of a strategy
// @Summary Duplicate a strategy
// @Description Create a copy of an existing strategy
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Strategy ID to duplicate" format(uuid)
// @Success 201 {object} dto.StrategyResponse "Duplicated strategy"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Strategy not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/{id}/clone [post]
func (h *StrategyHandler) DuplicateStrategy(c *gin.Context) {
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

	// Check if user is authorized to duplicate the strategy
	strategy, err := h.strategyService.GetStrategyByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get strategy: " + err.Error(),
			})
		}
		return
	}

	if strategy.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "You are not authorized to duplicate this strategy",
		})
		return
	}

	// Duplicate strategy
	duplicatedStrategy, err := h.strategyService.DuplicateStrategy(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to duplicate strategy: " + err.Error(),
			})
		}
		return
	}

	// Return duplicated strategy
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    duplicatedStrategy,
	})
}