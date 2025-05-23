package http

import (
	"net/http"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeploymentLogController handles deployment log HTTP requests
type DeploymentLogController struct {
	deploymentService service.DeploymentService
}

// NewDeploymentLogController creates a new deployment log controller
func NewDeploymentLogController(deploymentService service.DeploymentService) *DeploymentLogController {
	return &DeploymentLogController{
		deploymentService: deploymentService,
	}
}

// CreateDeploymentLog handles the request to create a new deployment log entry
func (c *DeploymentLogController) CreateDeploymentLog(ctx *gin.Context) {
	var req dto.DeploymentLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Create the log entry
	response, err := c.deploymentService.CreateDeploymentLog(ctx, &req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create deployment log: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetDeploymentLogs handles the request to get all logs for a deployment
func (c *DeploymentLogController) GetDeploymentLogs(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deploymentID")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid deployment ID: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get the logs
	logs, err := c.deploymentService.GetDeploymentLogs(ctx, deploymentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get deployment logs: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
	})
}

// GetLatestDeploymentLog handles the request to get the latest log for a deployment
func (c *DeploymentLogController) GetLatestDeploymentLog(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deploymentID")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid deployment ID: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get the latest log
	log, err := c.deploymentService.GetLatestDeploymentLog(ctx, deploymentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get latest deployment log: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    log,
	})
}
