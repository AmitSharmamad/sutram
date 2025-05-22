package http

import (
	"algo/internal/infra/http/middleware"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the deployment module routes
func RegisterRoutes(router *gin.RouterGroup, controller *DeploymentController, userController *userHttp.UserController) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	// All deployment routes require authentication
	deploymentAuth := router.Group("/deployment")
	deploymentAuth.Use(middleware.AuthMiddleware(authServiceAdapter))
	{
		// Deployment management
		deploymentAuth.POST("", controller.CreateDeployment)
		deploymentAuth.GET("", controller.GetUserDeployments)
		deploymentAuth.GET("/search", controller.SearchDeployments)
		deploymentAuth.GET("/:deployment_id", controller.GetDeploymentByID)
		deploymentAuth.PUT("/:deployment_id", controller.UpdateDeployment)
		deploymentAuth.DELETE("/:deployment_id", controller.DeleteDeployment)
		deploymentAuth.PUT("/:deployment_id/activate", controller.ActivateDeployment)
		deploymentAuth.PUT("/:deployment_id/deactivate", controller.DeactivateDeployment)
		deploymentAuth.POST("/stop-all", controller.StopAllDeployments)
		deploymentAuth.POST("/start-all", controller.StartAllDeployments)

		// Strategy-specific deployments
		deploymentAuth.GET("/strategy/:strategy_id", controller.GetDeploymentsByStrategy)
		deploymentAuth.POST("/all", controller.DeployAllInstruments)
	}
}
