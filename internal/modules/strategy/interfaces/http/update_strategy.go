// UpdateStrategy handles the update of an existing strategy
// @Summary Update a strategy
// @Description Update an existing strategy with the provided parameters
// @Tags Strategy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Strategy ID" format(uuid)
// @Param strategy body dto.UpdateStrategyRequest true "Updated strategy details"
// @Success 200 {object} dto.StrategyResponse "Updated strategy"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Strategy not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /strategy/{id} [put]
func (h *StrategyHandler) UpdateStrategy(c *gin.Context) {
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

	// Parse request body
	var req dto.UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Ensure ID in body matches URL
	if req.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID in body does not match ID in URL",
		})
		return
	}

	// Check if user is authorized to update the strategy
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
			"error":   "You are not authorized to update this strategy",
		})
		return
	}

	// Update strategy
	updatedStrategy, err := h.strategyService.UpdateStrategy(c.Request.Context(), &req)
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
		case service.ErrStrategyNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Strategy not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update strategy: " + err.Error(),
			})
		}
		return
	}

	// Return updated strategy
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedStrategy,
	})
}