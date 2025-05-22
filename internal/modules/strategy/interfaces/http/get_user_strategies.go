// GetUserStrategies handles the retrieval of all strategies for the current user
// @Summary Get user strategies
// @Description Retrieve all strategies for the current authenticated user
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.StrategyResponse "List of strategies"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy [get]
func (h *StrategyHandler) GetUserStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Get strategies for user
	filter := dto.StrategyFilter{
		UserID: &userID.(uuid.UUID),
	}

	strategies, err := h.strategyService.GetStrategies(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get strategies: " + err.Error(),
		})
		return
	}

	// Return strategies
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategies,
	})
}