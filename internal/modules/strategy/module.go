package strategy

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	deploymentRepository "algo/internal/modules/deployment/domain/repository"
	"algo/internal/modules/strategy/application/service"
	"algo/internal/modules/strategy/domain/repository"
	"algo/internal/modules/strategy/infrastructure/persistence/postgres"
	strategyHttp "algo/internal/modules/strategy/interfaces/http"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the strategy module
type Module struct {
	db               *sql.DB
	controller       *strategyHttp.StrategyHandler
	userController   *userHttp.UserController
	service          service.StrategyService
	optionLegService service.OptionLegService
	deploymentRepo   deploymentRepository.DeploymentRepository
}

// ProvideModule provides a wire provider set for the strategy module
var ProvideModule = wire.NewSet(
	ProvideStrategyModule,
)

// ProvideStrategyModule provides a strategy module instance
func ProvideStrategyModule() (*Module, error) {
	return New(), nil
}

// New creates a new strategy module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	module := &Module{
		db: db,
	}

	// Initialize module components
	module.initialize()

	return module
}

// initialize sets up the strategy module
func (m *Module) initialize() {
	// Create repository
	strategyRepo := ProvideStrategyRepository(m.db)

	// Create option leg service early
	optionLegRepo := ProvideOptionLegRepository(m.db)
	m.optionLegService = ProvideOptionLegService(optionLegRepo)

	// Note: deployment repository will be set later to avoid circular dependency
	// For now, create service with nil deployment repo - this will be updated later
	m.service = ProvideStrategyService(strategyRepo, nil, m.optionLegService)

	// Do NOT create controller here - wait until deployment repository is set
	// m.controller = ProvideStrategyHandler(m.service)
}

// SetDeploymentRepository sets the deployment repository and initializes the service and controller
func (m *Module) SetDeploymentRepository(deploymentRepo deploymentRepository.DeploymentRepository) {
	m.deploymentRepo = deploymentRepo

	// Recreate the service with the deployment repository
	strategyRepo := ProvideStrategyRepository(m.db)

	// Option leg service should already be created in initialize
	if m.optionLegService == nil {
		optionLegRepo := ProvideOptionLegRepository(m.db)
		m.optionLegService = ProvideOptionLegService(optionLegRepo)
	}

	m.service = ProvideStrategyService(strategyRepo, deploymentRepo, m.optionLegService)

	// Create the controller with the updated service (first time creation)
	m.controller = ProvideStrategyHandler(m.service, m.optionLegService)
}

// RegisterRoutes registers the strategy module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	if m.controller == nil {
		panic("Strategy controller not initialized. Call SetDeploymentRepository first.")
	}
	m.controller.RegisterRoutes(router, m.userController)
}

// SetUserController sets the user controller instance
func (m *Module) SetUserController(userController *userHttp.UserController) {
	m.userController = userController
}

// GetController returns the strategy controller
func (m *Module) GetController() *strategyHttp.StrategyHandler {
	return m.controller
}

// GetStrategyService returns the strategy service
func (m *Module) GetStrategyService() service.StrategyService {
	return m.service
}

// ProvideStrategyRepository provides a strategy repository
func ProvideStrategyRepository(db *sql.DB) repository.StrategyRepository {
	return postgres.NewStrategyRepository(db)
}

// ProvideStrategyService provides a strategy service
func ProvideStrategyService(repo repository.StrategyRepository, deploymentRepo deploymentRepository.DeploymentRepository, optionLegService service.OptionLegService) service.StrategyService {
	return service.NewStrategyService(repo, deploymentRepo, optionLegService)
}

// ProvideOptionLegRepository provides an option leg repository
func ProvideOptionLegRepository(db *sql.DB) repository.OptionLegRepository {
	return postgres.NewOptionLegRepository(db)
}

// ProvideOptionLegService provides an option leg service
func ProvideOptionLegService(repo repository.OptionLegRepository) service.OptionLegService {
	return service.NewOptionLegService(repo)
}

// ProvideStrategyHandler provides a strategy controller
func ProvideStrategyHandler(service service.StrategyService, optionLegService service.OptionLegService) *strategyHttp.StrategyHandler {
	return strategyHttp.NewStrategyHandler(service, optionLegService)
}
