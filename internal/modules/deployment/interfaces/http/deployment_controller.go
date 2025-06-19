package http

import (
	"fmt"
	"net/http"
	"strings"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeploymentController handles deployment HTTP requests
type DeploymentController struct {
	deploymentService service.DeploymentService
}

// NewDeploymentController creates a new deployment controller
func NewDeploymentController(deploymentService service.DeploymentService) *DeploymentController {
	return &DeploymentController{
		deploymentService: deploymentService,
	}
}

// CreateDeployment handles the request to create a new deployment
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

		switch err {
		case service.ErrStrategyNotFound:
			status = http.StatusNotFound
			message = "Strategy not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to deploy this strategy"
		case service.ErrDeploymentAlreadyExists:
			status = http.StatusConflict
			message = "Deployment already exists for this strategy and instrument"
		case service.ErrInvalidRequest:
			status = http.StatusBadRequest
			message = "Invalid request data"
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

// GetDeploymentByID handles the request to get a deployment by ID
func (c *DeploymentController) GetDeploymentByID(ctx *gin.Context) {
	// Parse deployment ID from URL
	idParam := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(idParam)
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

	deployment, err := c.deploymentService.GetDeploymentByID(ctx, deploymentID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to get deployment"

		switch err {
		case service.ErrDeploymentNotFound:
			status = http.StatusNotFound
			message = "Deployment not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to access this deployment"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deployment,
	})
}

// GetUserDeployments handles the request to get all deployments for the current user
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

	// Check for status filter
	status := ctx.Query("status")
	var filter *dto.DeploymentFilter
	if status != "" {
		if status != "waiting" && status != "entered" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid status filter. Use 'waiting' or 'entered'",
			})
			return
		}
		filter = &dto.DeploymentFilter{
			Status: &status,
		}
	}

	var deployments []*dto.DeploymentResponse
	var err error

	if filter != nil {
		deployments, err = c.deploymentService.GetUserDeploymentsWithFilter(ctx, userID, filter)
	} else {
		deployments, err = c.deploymentService.GetUserDeployments(ctx, userID)
	}

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

// GetDeploymentsByStrategy handles the request to get all deployments for a strategy
func (c *DeploymentController) GetDeploymentsByStrategy(ctx *gin.Context) {
	// Parse strategy ID from URL
	idParam := ctx.Param("strategy_id")
	strategyID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid strategy ID",
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

	deployments, err := c.deploymentService.GetDeploymentsByStrategyID(ctx, strategyID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to get deployments"

		switch err {
		case service.ErrStrategyNotFound:
			status = http.StatusNotFound
			message = "Strategy not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to access this strategy"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deployments,
	})
}

// GetDeploymentStatus handles the request to get deployment status for specified instruments against a strategy
func (c *DeploymentController) GetDeploymentStatus(ctx *gin.Context) {
	// Parse strategy ID from URL
	idParam := ctx.Param("strategy_id")
	strategyID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid strategy ID",
			"error":   err.Error(),
		})
		return
	}

	// Parse instruments query parameter
	instrumentsParam := ctx.Query("instruments")
	if instrumentsParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "instruments query parameter is required",
		})
		return
	}

	// Parse comma-separated instrument IDs
	var instrumentIDs []int
	instrumentStrings := strings.Split(instrumentsParam, ",")
	for _, idStr := range instrumentStrings {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		var instrumentID int
		_, err := fmt.Sscanf(idStr, "%d", &instrumentID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Invalid instrument ID: %s", idStr),
				"error":   err.Error(),
			})
			return
		}
		instrumentIDs = append(instrumentIDs, instrumentID)
	}

	if len(instrumentIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No valid instrument IDs provided",
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

	deploymentStatus, err := c.deploymentService.GetDeploymentStatus(ctx, strategyID, instrumentIDs, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to get deployment status"

		switch err {
		case service.ErrStrategyNotFound:
			status = http.StatusNotFound
			message = "Strategy not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to access this strategy"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deploymentStatus,
	})
}

// UpdateDeployment handles the request to update a deployment
func (c *DeploymentController) UpdateDeployment(ctx *gin.Context) {
	var req dto.UpdateDeploymentRequest
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

	deployment, err := c.deploymentService.UpdateDeployment(ctx, &req, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update deployment"

		switch err {
		case service.ErrDeploymentNotFound:
			status = http.StatusNotFound
			message = "Deployment not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to update this deployment"
		case service.ErrInvalidRequest:
			status = http.StatusBadRequest
			message = "Invalid request data"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deployment updated successfully",
		"data":    deployment,
	})
}

// DeleteDeployment handles the request to delete a deployment
func (c *DeploymentController) DeleteDeployment(ctx *gin.Context) {
	// Parse deployment ID from URL
	idParam := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(idParam)
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

	err = c.deploymentService.DeleteDeployment(ctx, deploymentID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to delete deployment"

		switch err {
		case service.ErrDeploymentNotFound:
			status = http.StatusNotFound
			message = "Deployment not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to delete this deployment"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deployment deleted successfully",
	})
}

// ActivateDeployment handles the request to activate a deployment
func (c *DeploymentController) ActivateDeployment(ctx *gin.Context) {
	// Parse deployment ID from URL
	idParam := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(idParam)
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

	err = c.deploymentService.ActivateDeployment(ctx, deploymentID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to activate deployment"

		switch err {
		case service.ErrDeploymentNotFound:
			status = http.StatusNotFound
			message = "Deployment not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to activate this deployment"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deployment activated successfully",
	})
}

// DeactivateDeployment handles the request to deactivate a deployment
func (c *DeploymentController) DeactivateDeployment(ctx *gin.Context) {
	// Parse deployment ID from URL
	idParam := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(idParam)
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

	err = c.deploymentService.DeactivateDeployment(ctx, deploymentID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to deactivate deployment"

		switch err {
		case service.ErrDeploymentNotFound:
			status = http.StatusNotFound
			message = "Deployment not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to deactivate this deployment"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deployment deactivated successfully",
	})
}

// SearchDeployments handles the request to search deployments with filters
func (c *DeploymentController) SearchDeployments(ctx *gin.Context) {
	var filter dto.DeploymentFilter

	// Parse query parameters
	if strategyIDStr := ctx.Query("strategy_id"); strategyIDStr != "" {
		strategyID, err := uuid.Parse(strategyIDStr)
		if err == nil {
			filter.StrategyID = &strategyID
		}
	}

	if instrumentIDStr := ctx.Query("instrument_id"); instrumentIDStr != "" {
		var instrumentID int
		_, err := fmt.Sscanf(instrumentIDStr, "%d", &instrumentID)
		if err == nil {
			filter.InstrumentID = &instrumentID
		}
	}

	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	if name := ctx.Query("name"); name != "" {
		filter.Name = &name
	}

	if status := ctx.Query("status"); status != "" {
		if status != "waiting" && status != "entered" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid status filter. Use 'waiting' or 'entered'",
			})
			return
		}
		filter.Status = &status
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

	deployments, err := c.deploymentService.SearchDeployments(ctx, &filter, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to search deployments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deployments,
	})
}

// StopAllDeployments handles the request to stop all deployments for the current user
func (c *DeploymentController) StopAllDeployments(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Stop all deployments for the user
	err := c.deploymentService.StopAllDeployments(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to stop deployments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All deployments stopped successfully",
	})
}

// DeployAllInstruments handles the request to deploy all instruments for a specific strategy
func (c *DeploymentController) DeployAllInstruments(ctx *gin.Context) {
	var strategyID uuid.UUID
	var err error

	// Check if strategy_id is provided as a query parameter
	strategyIDParam := ctx.Query("strategy_id")
	if strategyIDParam != "" {
		strategyID, err = uuid.Parse(strategyIDParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid strategy ID in query parameter",
				"error":   err.Error(),
			})
			return
		}
	} else {
		// If not in query, try to get from JSON body
		var req struct {
			StrategyID uuid.UUID `json:"strategy_id" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Strategy ID required in request body or as query parameter",
				"error":   err.Error(),
			})
			return
		}
		strategyID = req.StrategyID
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

	deployments, err := c.deploymentService.DeployAllInstruments(ctx.Request.Context(), strategyID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to deploy instruments"

		switch err {
		case service.ErrStrategyNotFound:
			status = http.StatusNotFound
			message = "Strategy not found"
		case service.ErrUnauthorized:
			status = http.StatusForbidden
			message = "You don't have permission to deploy this strategy"
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Instruments deployed successfully",
		"data":    deployments,
	})
}

// StartAllDeployments handles the request to start all paused deployments for the current user
func (c *DeploymentController) StartAllDeployments(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Start all paused deployments for the user
	err := c.deploymentService.StartAllDeployments(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start deployments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All paused deployments started successfully",
	})
}

// Helper function to get user ID from context
func getUserIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	// Get the user ID from the context (set by middleware)
	userIDVal, exists := ctx.Get("userID")
	if exists {
		userID, ok := userIDVal.(uuid.UUID)
		if ok {
			return userID, true
		}
	}

	// Fallback to user_id for backward compatibility
	userIDVal, exists = ctx.Get("user_id")
	if exists {
		userID, ok := userIDVal.(uuid.UUID)
		if ok {
			return userID, true
		}
	}

	return uuid.UUID{}, false
}
