package deployment

import (
	"database/sql"
	"fmt"

	"algo/internal/infra/http/middleware"
	infraRepo "algo/internal/infra/persistance/repository"
	backtestEngine "algo/internal/modules/backtest/application/service"
	"algo/internal/modules/deployment/application/service"
	"algo/internal/modules/deployment/infrastructure/persistence/postgres"
	"algo/internal/modules/deployment/interfaces/http"
	instrumentRepo "algo/internal/modules/instrument/infrastructure/persistence/clickhouse"
	strategyModule "algo/internal/modules/strategy"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the deployment module
type Module struct {
	db                 *sql.DB
	controller         *http.DeploymentController
	logController      *http.DeploymentLogController
	positionController *http.DeploymentPositionController
	metricsController  *http.DeploymentMetricsController
	userController     *userHttp.UserController
	strategyModule     *strategyModule.Module
}

// ProvideModule provides a wire provider set for the deployment module
var ProvideModule = wire.NewSet(
	ProvideDeploymentModule,
)

// ProvideDeploymentModule provides a deployment module instance
func ProvideDeploymentModule() (*Module, error) {
	return New(), nil
}

// New creates a new deployment module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	module := &Module{
		db: db,
	}

	return module
}

// Initialize initializes the deployment module with dependencies
func (m *Module) Initialize(userController *userHttp.UserController, strategyModule *strategyModule.Module) {
	m.userController = userController
	m.strategyModule = strategyModule

	// Initialize module components
	m.initializeComponents()
}

// initializeComponents sets up the deployment module components
func (m *Module) initializeComponents() {
	// Get sqlx DB instance
	dbx := infraRepo.ProvideSqlxDB()

	// Create repositories
	deploymentRepo := postgres.NewDeploymentRepository(m.db)
	logRepo := postgres.NewPgDeploymentLogRepository(dbx)
	positionRepo := postgres.NewPgDeploymentPositionRepository(dbx)
	metricsRepo := postgres.NewPgDeploymentMetricsRepository(dbx)

	// Get strategy repository from strategy module
	strategyService := m.strategyModule.GetStrategyService()

	// Get Clickhouse DB for instrument repository
	chDB, err := infraRepo.ProvideClickHouseDB()
	if err != nil {
		// In a real application, handle this error properly
		fmt.Printf("Error connecting to ClickHouse: %v\n", err)
		// Create a nil instrument repository as a fallback
		// In a real app, you would handle this better
	}
	instrumentRepository := instrumentRepo.NewInstrumentRepository(chDB)

	// For now, create a nil backtest engine
	// In a real implementation, we would obtain this from the backtest module
	var engine *backtestEngine.BacktestEngine = nil

	// Create services
	deploymentService := service.NewDeploymentService(
		deploymentRepo,
		strategyService.GetRepository(),
		logRepo,
		positionRepo,
		metricsRepo,
		instrumentRepository,
		engine,
	)

	// Create controllers
	m.controller = http.NewDeploymentController(deploymentService)
	logController := http.NewDeploymentLogController(deploymentService)
	positionController := http.NewDeploymentPositionController(deploymentService)
	metricsController := http.NewDeploymentMetricsController(deploymentService)

	// Store controllers in the module
	m.logController = logController
	m.positionController = positionController
	m.metricsController = metricsController
}

// RegisterRoutes registers the deployment module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Register main deployment routes using the router.go implementation
	http.RegisterRoutes(router, m.controller, m.userController)

	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(m.userController.GetAuthService())

	// Create middleware
	authMiddleware := middleware.AuthMiddleware(authServiceAdapter)

	// Register log routes
	logRoutes := router.Group("/deployment")
	{
		logRoutes.POST("/:deployment_id/logs", authMiddleware, m.logController.CreateDeploymentLog)
		logRoutes.GET("/:deployment_id/logs", authMiddleware, m.logController.GetDeploymentLogs)
		logRoutes.GET("/:deployment_id/logs/latest", authMiddleware, m.logController.GetLatestDeploymentLog)
	}
	// Register position routes
	positionRoutes := router.Group("/deployment")
	{
		positionRoutes.POST("/:deployment_id/positions", authMiddleware, m.positionController.CreatePosition)
		positionRoutes.GET("/:deployment_id/positions", authMiddleware, m.positionController.GetPositionsByDeployment)
		positionRoutes.GET("/:deployment_id/positions/open", authMiddleware, m.positionController.GetOpenPositionsByDeployment)
		positionRoutes.GET("/positions/:positionID", authMiddleware, m.positionController.GetPositionByID)
		positionRoutes.PUT("/positions/:positionID/close", authMiddleware, m.positionController.ClosePosition)
		positionRoutes.PUT("/positions/:positionID/cancel", authMiddleware, m.positionController.CancelPosition)

		// Add new route format as specified in the manual
		positionRoutes.PUT("/:deployment_id/positions/:positionID/close", authMiddleware, m.positionController.ClosePosition)
		positionRoutes.PUT("/:deployment_id/positions/:positionID/cancel", authMiddleware, m.positionController.CancelPosition)
	}

	// Register position routes with new format (deployment/:deployment_id/positions/:positionID/...)
	// These routes have been implemented above

	// Register metrics routes
	metricsRoutes := router.Group("/deployment")
	{
		metricsRoutes.POST("/:deployment_id/metrics", authMiddleware, m.metricsController.CreateOrUpdateMetrics)
		metricsRoutes.GET("/:deployment_id/metrics", authMiddleware, m.metricsController.GetMetricsByDeployment)
		metricsRoutes.GET("/:deployment_id/metrics/latest", authMiddleware, m.metricsController.GetLatestMetrics)
		metricsRoutes.GET("/:deployment_id/metrics/:date", authMiddleware, m.metricsController.GetMetricsByDate)
		metricsRoutes.GET("/:deployment_id/metrics/range", authMiddleware, m.metricsController.GetMetricsTimeRange)
	}
}

// GetController returns the deployment controller
func (m *Module) GetController() *http.DeploymentController {
	return m.controller
}

// GetLogController returns the deployment log controller
func (m *Module) GetLogController() *http.DeploymentLogController {
	return m.logController
}

// GetPositionController returns the deployment position controller
func (m *Module) GetPositionController() *http.DeploymentPositionController {
	return m.positionController
}

// GetMetricsController returns the deployment metrics controller
func (m *Module) GetMetricsController() *http.DeploymentMetricsController {
	return m.metricsController
}
