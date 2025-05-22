// GetDefaultStrategies handles the retrieval of default strategies
// @Summary Get default strategies
// @Description Retrieve all default/system strategies
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.StrategyResponse "List of default strategies"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/default [get]
func (h *StrategyHandler) GetDefaultStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Get default strategies
	strategies, err := h.strategyService.GetDefaultStrategies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get default strategies: " + err.Error(),
		})
		return
	}

	// Return strategies
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategies,
	})
}