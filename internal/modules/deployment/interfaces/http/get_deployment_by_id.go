// GetDeploymentByID handles the request to get a deployment by ID
// @Summary Get a deployment by ID
// @Description Retrieve a deployment by its ID
// @Tags Deployment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deployment ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Deployment details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Deployment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /deployment/{id} [get]
func (c *DeploymentController) GetDeploymentByID(ctx *gin.Context) {
	// Get deployment ID from URL
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid deployment ID",
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

	// Get deployment by ID
	deployment, err := c.deploymentService.GetDeploymentByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to get deployment"

		if err == service.ErrDeploymentNotFound {
			status = http.StatusNotFound
			message = "Deployment not found"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	// Check if user is authorized to access this deployment
	if deployment.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You are not authorized to access this deployment",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deployment,
	})
}