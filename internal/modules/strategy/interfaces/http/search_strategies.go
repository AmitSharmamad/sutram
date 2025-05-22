// SearchStrategies handles the search for strategies
// @Summary Search strategies
// @Description Search for strategies with various filters
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param filter body dto.StrategyFilter true "Search filters"
// @Success 200 {array} dto.StrategyResponse "List of matching strategies"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/search [post]
func (h *StrategyHandler) SearchStrategies(c *gin.Context) {
	// Extract userID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse request body
	var filter dto.StrategyFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Always filter by the current user's ID for security
	userUUID := userID.(uuid.UUID)
	filter.UserID = &userUUID

	// Search strategies
	strategies, err := h.strategyService.GetStrategies(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to search strategies: " + err.Error(),
		})
		return
	}

	// Return strategies
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategies,
	})
}