package http

import (
	"net/http"
	"time"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeploymentMetricsController handles deployment metrics HTTP requests
type DeploymentMetricsController struct {
	deploymentService service.DeploymentService
}

// NewDeploymentMetricsController creates a new deployment metrics controller
func NewDeploymentMetricsController(deploymentService service.DeploymentService) *DeploymentMetricsController {
	return &DeploymentMetricsController{
		deploymentService: deploymentService,
	}
}

// CreateOrUpdateMetrics handles the request to create or update metrics
func (c *DeploymentMetricsController) CreateOrUpdateMetrics(ctx *gin.Context) {
	var req dto.CreateMetricsRequest
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

	// Create or update the metrics
	response, err := c.deploymentService.CreateOrUpdateMetrics(ctx, &req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create or update metrics: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetMetricsByDeployment handles the request to get all metrics for a deployment
func (c *DeploymentMetricsController) GetMetricsByDeployment(ctx *gin.Context) {
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

	// Get the metrics
	metrics, err := c.deploymentService.GetMetricsByDeployment(ctx, deploymentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get metrics: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetMetricsByDate handles the request to get metrics for a deployment and date
func (c *DeploymentMetricsController) GetMetricsByDate(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deploymentID")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid deployment ID: " + err.Error(),
		})
		return
	}

	dateStr := ctx.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid date format (use YYYY-MM-DD): " + err.Error(),
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

	// Get the metrics
	metrics, err := c.deploymentService.GetMetricsByDeploymentAndDate(ctx, deploymentID, date, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get metrics: " + err.Error(),
		})
		return
	}

	if metrics == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Metrics not found for the given date",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetMetricsTimeRange handles the request to get metrics for a time range
func (c *DeploymentMetricsController) GetMetricsTimeRange(ctx *gin.Context) {
	deploymentIDStr := ctx.Param("deploymentID")
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid deployment ID: " + err.Error(),
		})
		return
	}

	// Parse query parameters
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	filter := &dto.MetricsFilter{
		DeploymentID: deploymentID,
	}

	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid start date format (use YYYY-MM-DD): " + err.Error(),
			})
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid end date format (use YYYY-MM-DD): " + err.Error(),
			})
			return
		}
		filter.EndDate = &endDate
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

	// Get the metrics
	metrics, err := c.deploymentService.GetMetricsTimeRange(ctx, filter, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get metrics for time range: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetLatestMetrics handles the request to get the latest metrics for a deployment
func (c *DeploymentMetricsController) GetLatestMetrics(ctx *gin.Context) {
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

	// Get the latest metrics
	metrics, err := c.deploymentService.GetLatestMetrics(ctx, deploymentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get latest metrics: " + err.Error(),
		})
		return
	}

	if metrics == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "No metrics found for this deployment",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}
