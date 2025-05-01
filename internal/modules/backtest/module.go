package backtest

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/backtest/application/service"
	"algo/internal/modules/backtest/infrastructure/persistence/postgres"
	backtestHttp "algo/internal/modules/backtest/interfaces/http"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the backtest module
type Module struct {
	db             *sql.DB
	controller     *backtestHttp.BacktestController
	userController interface{}
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

	module := &Module{
		db: db,
	}

	// Initialize module components
	module.initialize()

	return module
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

	// Create controller
	m.controller = backtestHttp.NewBacktestController(backtestService)
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

	backtestHttp.RegisterRoutes(router, m.controller, userController)
}

// GetController returns the backtest controller
func (m *Module) GetController() *backtestHttp.BacktestController {
	return m.controller
}

// SetUserController sets the user controller for auth middleware
func (m *Module) SetUserController(userController interface{}) {
	m.userController = userController
}
