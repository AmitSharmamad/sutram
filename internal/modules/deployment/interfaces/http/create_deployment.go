// CreateDeployment handles the request to create a new deployment
// @Summary Create a deployment
// @Description Create a new deployment for a strategy
// @Tags Deployment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param deployment body dto.CreateDeploymentRequest true "Deployment details"
// @Success 201 {object} map[string]interface{} "Created deployment"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Strategy not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /deployment [post]
func (c *DeploymentController) CreateDeployment(ctx *gin.Context) {
	var req dto.CreateDeploymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	deployment, err := c.deploymentService.CreateDeployment(ctx, &req, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to create deployment"

		// Check for specific errors to provide better responses
		if err == service.ErrStrategyNotFound {
			status = http.StatusNotFound
			message = "Strategy not found"
		} else if err == service.ErrInvalidRequest {
			status = http.StatusBadRequest
			message = "Invalid request"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Deployment created successfully",
		"data":    deployment,
	})
}