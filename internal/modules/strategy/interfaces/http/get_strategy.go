// GetStrategyByID handles the retrieval of a strategy by ID
// @Summary Get a strategy by ID
// @Description Retrieve a strategy by its unique ID
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Strategy ID" format(uuid)
// @Success 200 {object} dto.StrategyResponse "Strategy details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Strategy not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/{id} [get]
func (h *StrategyHandler) GetStrategyByID(c *gin.Context) {
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

	// Get strategy
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

	// Check if user is authorized to access the strategy
	if strategy.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "You are not authorized to access this strategy",
		})
		return
	}

	// Return strategy
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategy,
	})
}