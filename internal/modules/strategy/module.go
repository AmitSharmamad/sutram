package strategy

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
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
	db             *sql.DB
	controller     *strategyHttp.StrategyHandler
	userController *userHttp.UserController
	service        service.StrategyService
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

	// Create service
	m.service = ProvideStrategyService(strategyRepo)

	// Create controller
	m.controller = ProvideStrategyHandler(m.service)
}

// RegisterRoutes registers the strategy module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
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
func ProvideStrategyService(repo repository.StrategyRepository) service.StrategyService {
	return service.NewStrategyService(repo)
}

// ProvideStrategyHandler provides a strategy handler
func ProvideStrategyHandler(svc service.StrategyService) *strategyHttp.StrategyHandler {
	return strategyHttp.NewStrategyHandler(svc)
}
