package http

import (
	"net/http"
	"strconv"
	"time"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeploymentPositionController handles deployment position HTTP requests
type DeploymentPositionController struct {
	deploymentService service.DeploymentService
}

// NewDeploymentPositionController creates a new deployment position controller
func NewDeploymentPositionController(deploymentService service.DeploymentService) *DeploymentPositionController {
	return &DeploymentPositionController{
		deploymentService: deploymentService,
	}
}

// CreatePosition handles the request to create a new position
func (c *DeploymentPositionController) CreatePosition(ctx *gin.Context) {
	var req dto.CreatePositionRequest
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

	// Create the position
	response, err := c.deploymentService.CreatePosition(ctx, &req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create position: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetPositionByID handles the request to get a position by ID
func (c *DeploymentPositionController) GetPositionByID(ctx *gin.Context) {
	positionIDStr := ctx.Param("positionID")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid position ID: " + err.Error(),
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

	// Get the position
	position, err := c.deploymentService.GetPositionByID(ctx, positionID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get position: " + err.Error(),
		})
		return
	}

	if position == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Position not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    position,
	})
}

// GetPositionsByDeployment handles the request to get all positions for a deployment
func (c *DeploymentPositionController) GetPositionsByDeployment(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deployment_id")
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

	// Get the positions
	positions, err := c.deploymentService.GetPositionsByDeployment(ctx, deploymentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get positions: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    positions,
	})
}

// GetOpenPositionsByDeployment handles the request to get all open positions for a deployment
func (c *DeploymentPositionController) GetOpenPositionsByDeployment(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deployment_id")
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

	// Get the open positions
	positions, err := c.deploymentService.GetOpenPositionsByDeployment(ctx, deploymentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get open positions: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    positions,
	})
}

// ClosePosition handles the request to close a position
func (c *DeploymentPositionController) ClosePosition(ctx *gin.Context) {
	positionIDStr := ctx.Param("positionID")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid position ID: " + err.Error(),
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

	// If deployment_id is provided in the URL, validate it
	deploymentIDStr := ctx.Param("deployment_id")
	if deploymentIDStr == "" {
		// Try the deploymentID parameter format for backward compatibility
		deploymentIDStr = ctx.Param("deploymentID")
	}

	if deploymentIDStr != "" {
		// This is a sanity check to ensure the position belongs to the deployment
		deploymentID, err := uuid.Parse(deploymentIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid deployment ID: " + err.Error(),
			})
			return
		}

		// Get the position first to check its deployment ID
		position, err := c.deploymentService.GetPositionByID(ctx, positionID, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get position: " + err.Error(),
			})
			return
		}

		if position == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Position not found",
			})
			return
		}

		// Ensure the position belongs to the specified deployment
		if position.DeploymentID != deploymentID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Position does not belong to the specified deployment",
			})
			return
		}
	}

	// Parse request body if available
	var req dto.ClosePositionRequest
	if ctx.Request.ContentLength > 0 {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid request: " + err.Error(),
			})
			return
		}
	} else {
		// If no request body, use current time and last price
		req.ExitTime = time.Now()
		// Use the entry price as the exit price if we don't have another value
		// This is a fallback and might need to be adjusted based on business logic
		position, err := c.deploymentService.GetPositionByID(ctx, positionID, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get position: " + err.Error(),
			})
			return
		}
		req.ExitPrice = position.EntryPrice
	}

	// Close the position
	response, err := c.deploymentService.ClosePosition(ctx, positionID, &req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to close position: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// CancelPosition handles the request to cancel a position
func (c *DeploymentPositionController) CancelPosition(ctx *gin.Context) {
	positionIDStr := ctx.Param("positionID")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid position ID: " + err.Error(),
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

	// If deployment_id is provided in the URL, validate it
	deploymentIDStr := ctx.Param("deployment_id")
	if deploymentIDStr == "" {
		// Try the deploymentID parameter format for backward compatibility
		deploymentIDStr = ctx.Param("deploymentID")
	}

	if deploymentIDStr != "" {
		// This is a sanity check to ensure the position belongs to the deployment
		deploymentID, err := uuid.Parse(deploymentIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid deployment ID: " + err.Error(),
			})
			return
		}

		// Get the position first to check its deployment ID
		position, err := c.deploymentService.GetPositionByID(ctx, positionID, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get position: " + err.Error(),
			})
			return
		}

		if position == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Position not found",
			})
			return
		}

		// Ensure the position belongs to the specified deployment
		if position.DeploymentID != deploymentID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Position does not belong to the specified deployment",
			})
			return
		}
	}

	// Cancel the position
	if err := c.deploymentService.CancelPosition(ctx, positionID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to cancel position: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Position cancelled successfully",
	})
}
