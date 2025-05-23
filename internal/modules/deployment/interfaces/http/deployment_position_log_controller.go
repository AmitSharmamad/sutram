package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"algo/internal/modules/deployment/application/service"
	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// DeploymentPositionLogController handles HTTP requests related to deployment position logs
type DeploymentPositionLogController struct {
	service *service.DeploymentPositionLogService
}

// NewDeploymentPositionLogController creates a new deployment position log controller
func NewDeploymentPositionLogController(service *service.DeploymentPositionLogService) *DeploymentPositionLogController {
	return &DeploymentPositionLogController{
		service: service,
	}
}

// CreateLog creates a new position log
func (c *DeploymentPositionLogController) CreateLog(ctx *gin.Context) {
	var input struct {
		DeploymentPositionID uuid.UUID          `json:"deployment_position_id" binding:"required"`
		Results              map[string]any `json:"results" binding:"required"`
		Timestamp            time.Time      `json:"timestamp" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert results to json.RawMessage
	resultsJSON, err := json.Marshal(input.Results)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process results data"})
		return
	}

	positionLog := model.NewDeploymentPositionLog(
		input.DeploymentPositionID,
		resultsJSON,
		input.Timestamp,
	)

	if err := c.service.CreateLog(ctx, positionLog); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, positionLog)
}

// GetLogByID retrieves a position log by its ID
func (c *DeploymentPositionLogController) GetLogByID(ctx *gin.Context) {
	id := ctx.Param("id")

	positionLog, err := c.service.GetLogByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if positionLog == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Position log not found"})
		return
	}

	ctx.JSON(http.StatusOK, positionLog)
}

// GetLogsByPositionID retrieves all logs for a specific position
func (c *DeploymentPositionLogController) GetLogsByPositionID(ctx *gin.Context) {
	positionIDStr := ctx.Param("positionID")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}

	logs, err := c.service.GetLogsByPositionID(ctx, positionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

// GetLogsByPositionIDAndTimeRange retrieves logs for a position within a specific time range
func (c *DeploymentPositionLogController) GetLogsByPositionIDAndTimeRange(ctx *gin.Context) {
	positionIDStr := ctx.Param("positionID")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}

	startTimeStr := ctx.Query("startTime")
	endTimeStr := ctx.Query("endTime")

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
			return
		}
	} else {
		// Default to 24 hours ago if not specified
		startTime = time.Now().Add(-24 * time.Hour)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
			return
		}
	} else {
		// Default to current time if not specified
		endTime = time.Now()
	}

	logs, err := c.service.GetLogsByPositionIDAndTimeRange(ctx, positionID, startTime, endTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

// GetLogsByDeploymentID retrieves all logs for all positions in a deployment
func (c *DeploymentPositionLogController) GetLogsByDeploymentID(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deployment_id")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid deployment ID: " + err.Error(),
		})
		return
	}

	logs, err := c.service.GetLogsByDeploymentID(ctx, deploymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get position logs: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
	})
}
