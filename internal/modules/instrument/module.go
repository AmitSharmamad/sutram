package instrument

import (
	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/instrument/application/service"
	"algo/internal/modules/instrument/infrastructure/persistence/clickhouse"
	"algo/internal/modules/instrument/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

// Module represents the instrument module
type Module struct {
	db                *sqlx.DB
	controller        *http.InstrumentController
	instrumentService *service.InstrumentService
}

// ProvideModule provides a wire provider set for the instrument module
var ProvideModule = wire.NewSet(
	ProvideInstrumentModule,
)

// ProvideInstrumentModule provides an instrument module instance
func ProvideInstrumentModule() (*Module, error) {
	return New(), nil
}

// New creates a new instrument module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvideClickHouseDB()
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

// initialize sets up the instrument module
func (m *Module) initialize() {
	// Create repositories
	instrumentRepo := clickhouse.NewInstrumentRepository(m.db)
	instrumentMasterRepo := clickhouse.NewInstrumentMasterRepository(m.db)

	// Create services
	m.instrumentService = service.NewInstrumentService(instrumentRepo, instrumentMasterRepo)

	// Create controller
	m.controller = http.NewInstrumentController(m.instrumentService)
}

// RegisterRoutes registers the instrument module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	http.RegisterRoutes(router, m.controller)
}

// GetController returns the instrument controller
func (m *Module) GetController() *http.InstrumentController {
	return m.controller
}

// GetInstrumentService returns the instrument service
func (m *Module) GetInstrumentService() *service.InstrumentService {
	return m.instrumentService
}
