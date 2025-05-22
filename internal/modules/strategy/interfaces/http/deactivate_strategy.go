// DeactivateStrategy handles the deactivation of a strategy
// @Summary Deactivate a strategy
// @Description Deactivate a strategy by its ID
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Strategy ID" format(uuid)
// @Success 200 {object} dto.StrategyResponse "Deactivated strategy"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Strategy not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/{id}/deactivate [post]
func (h *StrategyHandler) DeactivateStrategy(c *gin.Context) {
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

	// Check if user is authorized to deactivate the strategy
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
			"error":   "You are not authorized to deactivate this strategy",
		})
		return
	}

	// Deactivate strategy
	deactivatedStrategy, err := h.strategyService.DeactivateStrategy(c.Request.Context(), id)
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
				"error":   "Failed to deactivate strategy: " + err.Error(),
			})
		}
		return
	}

	// Return deactivated strategy
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deactivatedStrategy,
	})
}