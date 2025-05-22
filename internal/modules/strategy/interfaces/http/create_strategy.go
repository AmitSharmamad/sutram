// CreateStrategy handles the creation of a new strategy
// @Summary Create a new trading strategy
// @Description Create a new trading strategy with the provided parameters
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param strategy body dto.CreateStrategyRequest true "Strategy details"
// @Success 201 {object} dto.StrategyResponse "Created strategy"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy [post]
func (h *StrategyHandler) CreateStrategy(c *gin.Context) {
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
	var req dto.CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create strategy
	strategy, err := h.strategyService.CreateStrategy(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid request: " + err.Error(),
			})
		case service.ErrInvalidExpression:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid expression: " + err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to create strategy: " + err.Error(),
			})
		}
		return
	}

	// Return created strategy
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    strategy,
	})
}