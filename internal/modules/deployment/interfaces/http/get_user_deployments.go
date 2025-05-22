// GetUserDeployments handles the request to get all deployments for the current user
// @Summary Get user deployments
// @Description Retrieve all deployments for the current user
// @Tags Deployment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param active query boolean false "Filter by active status"
// @Param strategy_id query string false "Filter by strategy ID" format(uuid)
// @Success 200 {object} map[string]interface{} "List of deployments"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /deployment [get]
func (c *DeploymentController) GetUserDeployments(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Parse query parameters
	var filter dto.DeploymentFilter
	filter.UserID = &userID

	// Check for active filter
	activeStr := ctx.Query("active")
	if activeStr != "" {
		active := activeStr == "true"
		filter.Active = &active
	}

	// Check for strategy ID filter
	strategyIDStr := ctx.Query("strategy_id")
	if strategyIDStr != "" {
		strategyID, err := uuid.Parse(strategyIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid strategy ID",
				"error":   err.Error(),
			})
			return
		}
		filter.StrategyID = &strategyID
	}

	// Get deployments with filter
	deployments, err := c.deploymentService.GetDeployments(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get deployments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deployments,
	})
}