package backtest

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/backtest/application/service"
	"algo/internal/modules/backtest/infrastructure/persistence/postgres"
	backtestHttp "algo/internal/modules/backtest/interfaces/http"
	instrumentModule "algo/internal/modules/instrument"
	strategyModule "algo/internal/modules/strategy"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the backtest module
type Module struct {
	db               *sql.DB
	controller       *backtestHttp.BacktestController
	userController   interface{}
	strategyModule   *strategyModule.Module
	instrumentModule *instrumentModule.Module
	backtestService  service.BacktestService
}

// ProvideModule provides a wire provider set for the backtest module
var ProvideModule = wire.NewSet(
	ProvideBacktestModule,
)

// ProvideBacktestModule provides a backtest module instance
func ProvideBacktestModule() (*Module, error) {
	return New(), nil
}

// New creates a new backtest module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	// Create instrument module
	instrumentMod, err := instrumentModule.ProvideInstrumentModule()
	if err != nil {
		panic(err)
	}

	module := &Module{
		db:               db,
		instrumentModule: instrumentMod,
	}

	return module
}

// SetStrategyModule sets the strategy module (to be called from main.go)
func (m *Module) SetStrategyModule(strategyMod *strategyModule.Module) {
	m.strategyModule = strategyMod
	// Initialize module components after strategy module is set
	m.initialize()
}

// initialize sets up the backtest module
func (m *Module) initialize() {
	// Create repositories
	runGroupRepo := postgres.NewRunGroupRepository(m.db)
	runRepo := postgres.NewRunRepository(m.db)
	positionRepo := postgres.NewPositionRepository(m.db)
	runMetricsRepo := postgres.NewRunMetricsRepository(m.db)

	// Create services
	backtestService := service.NewBacktestService(
		runGroupRepo,
		runRepo,
		positionRepo,
		runMetricsRepo,
	)

	// Store the service
	m.backtestService = backtestService

	// Get strategy service
	strategyService := m.strategyModule.GetStrategyService()

	// Get instrument service
	instrumentService := m.instrumentModule.GetInstrumentService()

	// Create backtest engine service - fix the type mismatch
	// Either by changing how we pass the backtestService or adjusting the signature of NewBacktestEngine
	engineService := service.NewBacktestEngine(
		&backtestService, // Change to pointer if NewBacktestEngine expects a pointer
		strategyService,
		instrumentService,
		runGroupRepo,
		runRepo,
		positionRepo,
		runMetricsRepo,
	)

	// Create controller with both services
	m.controller = backtestHttp.NewBacktestController(backtestService, engineService)
}

// RegisterRoutes registers the backtest module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Cast userController to the expected type
	userController, ok := m.userController.(*userHttp.UserController)
	if !ok {
		// If userController is not set or is of wrong type, we can't proceed with auth
		// In a production app, you'd want to log this error properly
		panic("Backtest module: UserController not set or invalid type")
	}

	// Get usage limit service from user controller
	usageLimitService := userController.GetUsageLimitService()
	if usageLimitService == nil {
		panic("Backtest module: Usage limit service not available from user controller")
	}

	backtestHttp.RegisterRoutes(router, m.controller, userController, usageLimitService)
}

// GetController returns the backtest controller
func (m *Module) GetController() *backtestHttp.BacktestController {
	return m.controller
}

// SetUserController sets the user controller for auth middleware
func (m *Module) SetUserController(userController interface{}) {
	m.userController = userController
}

// GetBacktestService returns the backtest service
func (m *Module) GetBacktestService() service.BacktestService {
	return m.backtestService
}
