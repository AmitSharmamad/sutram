package broker

import (
	"os"

	"algo/internal/modules/broker/application/service"
	"algo/internal/modules/broker/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the broker module
type Module struct {
	controller            *http.BrokerController
	dhanOrderService      *service.DhanOrderService
	zerodhaOrderService   *service.ZerodhaOrderService
	angeloneOrderService  *service.AngelOneOrderService
	brokerService         *service.BrokerService
}

// ProvideModule provides a wire provider set for the broker module
var ProvideModule = wire.NewSet(
	ProvideBrokerModule,
)

// ProvideBrokerModule provides a broker module instance
func ProvideBrokerModule() (*Module, error) {
	return New(), nil
}

// New creates a new broker module
func New() *Module {
	module := &Module{}

	// Initialize module components
	module.initialize()

	return module
}

// initialize sets up the broker module
func (m *Module) initialize() {
	// Get Dhan API configuration from environment
	dhanBaseURL := getEnv("DHAN_API_BASE_URL", "https://api.dhan.co/v2")
	dhanAccessToken := getEnv("DHAN_ACCESS_TOKEN", "")

	// Get Zerodha API configuration from environment
	zerodhaBaseURL := getEnv("ZERODHA_API_BASE_URL", "https://api.kite.trade")
	zerodhaAPIKey := getEnv("ZERODHA_API_KEY", "")
	zerodhaAccessToken := getEnv("ZERODHA_ACCESS_TOKEN", "")

	// Get AngelOne API configuration from environment
	angeloneBaseURL := getEnv("ANGELONE_API_BASE_URL", "https://apiconnect.angelbroking.com")
	angeloneAPIKey := getEnv("ANGELONE_API_KEY", "")
	angeloneAccessToken := getEnv("ANGELONE_ACCESS_TOKEN", "")

	// Create Dhan order service
	m.dhanOrderService = service.NewDhanOrderService(dhanBaseURL, dhanAccessToken)

	// Create Zerodha order service
	m.zerodhaOrderService = service.NewZerodhaOrderService(zerodhaBaseURL, zerodhaAPIKey, zerodhaAccessToken)

	// Create AngelOne order service
	m.angeloneOrderService = service.NewAngelOneOrderService(angeloneBaseURL, angeloneAPIKey, angeloneAccessToken)

	// Create broker service (without repository dependencies for now since they're not implemented)
	// m.brokerService = service.NewBrokerService(brokerRepo, brokerUserRepo, orderRepo, m.dhanOrderService)

	// Create controller
	m.controller = http.NewBrokerController(m.dhanOrderService)
}

// RegisterRoutes registers the broker module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	brokerGroup := router.Group("/broker")
	http.RegisterRoutes(brokerGroup, m.controller)
}

// GetController returns the broker controller
func (m *Module) GetController() *http.BrokerController {
	return m.controller
}

// GetDhanOrderService returns the Dhan order service
func (m *Module) GetDhanOrderService() *service.DhanOrderService {
	return m.dhanOrderService
}

// GetZerodhaOrderService returns the Zerodha order service
func (m *Module) GetZerodhaOrderService() *service.ZerodhaOrderService {
	return m.zerodhaOrderService
}

// GetAngelOneOrderService returns the AngelOne order service
func (m *Module) GetAngelOneOrderService() *service.AngelOneOrderService {
	return m.angeloneOrderService
}

// GetBrokerService returns the broker service
func (m *Module) GetBrokerService() *service.BrokerService {
	return m.brokerService
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}