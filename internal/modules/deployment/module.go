package deployment

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/deployment/application/service"
	"algo/internal/modules/deployment/infrastructure/persistence/postgres"
	"algo/internal/modules/deployment/interfaces/http"
	strategyModule "algo/internal/modules/strategy"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the deployment module
type Module struct {
	db             *sql.DB
	controller     *http.DeploymentController
	userController *userHttp.UserController
	strategyModule *strategyModule.Module
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
	// Create repositories
	deploymentRepo := postgres.NewDeploymentRepository(m.db)

	// Get strategy repository from strategy module
	strategyService := m.strategyModule.GetStrategyService()

	// Create services
	deploymentService := service.NewDeploymentService(deploymentRepo, strategyService.GetRepository())

	// Create controller
	m.controller = http.NewDeploymentController(deploymentService)
}

// RegisterRoutes registers the deployment module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	http.RegisterRoutes(router, m.controller, m.userController)
}

// GetController returns the deployment controller
func (m *Module) GetController() *http.DeploymentController {
	return m.controller
}
