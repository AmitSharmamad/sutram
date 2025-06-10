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
	controller           *http.BrokerController
	dhanBroker     *service.DhanBroker
	zerodhaBroker  *service.ZerodhaBroker
	hdfcsecBroker   *service.HDFCSecBroker
	sharekhanBroker *service.SharekhanBroker
	angeloneBroker  *service.AngelOneBroker
	growwBroker     *service.GrowwBroker
	brokerService        *service.BrokerService
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

	// Get HDFC Securities API configuration from environment
	hdfcsecBaseURL := getEnv("HDFCSEC_API_BASE_URL", "https://api.hdfcsec.com/v1")
	hdfcsecAccessToken := getEnv("HDFCSEC_ACCESS_TOKEN", "")

	// Get Sharekhan API configuration from environment
	sharekhanBaseURL := getEnv("SHAREKHAN_API_BASE_URL", "https://api.sharekhan.com")
	sharekhanAPIKey := getEnv("SHAREKHAN_API_KEY", "")
	sharekhanAccessToken := getEnv("SHAREKHAN_ACCESS_TOKEN", "")

	// Get Groww API configuration from environment
	growwBaseURL := getEnv("GROWW_API_BASE_URL", "https://api.groww.in")
	growwAPIKey := getEnv("GROWW_API_KEY", "")
	growwAccessToken := getEnv("GROWW_ACCESS_TOKEN", "")
	// Create Dhan order service
	m.dhanBroker = service.NewDhanBroker(dhanBaseURL, dhanAccessToken)

	// Create Zerodha order service
	m.zerodhaBroker = service.NewZerodhaBroker(zerodhaBaseURL, zerodhaAPIKey, zerodhaAccessToken)

	// Create AngelOne order service
	m.angeloneBroker = service.NewAngelOneBroker(angeloneBaseURL, angeloneAPIKey, angeloneAccessToken)

	// Create HDFC Securities order service
	m.hdfcsecBroker = service.NewHDFCSecBroker(hdfcsecBaseURL, hdfcsecAccessToken)

	// Create Sharekhan order service
	m.sharekhanBroker = service.NewSharekhanBroker(sharekhanBaseURL, sharekhanAPIKey, sharekhanAccessToken)

	// Create Groww order service
	m.growwBroker = service.NewGrowwBroker(growwBaseURL, growwAPIKey, growwAccessToken)

	// Create broker service (without repository dependencies for now since they're not implemented)
	// m.brokerService = service.NewBrokerService(brokerRepo, brokerUserRepo, orderRepo, m.dhanBroker)

	// Create controller
	m.controller = http.NewBrokerController(m.dhanBroker)
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

// GetDhanBroker returns the Dhan order service
func (m *Module) GetDhanBroker() *service.DhanBroker {
	return m.dhanBroker
}

// GetZerodhaBroker returns the Zerodha order service
func (m *Module) GetZerodhaBroker() *service.ZerodhaBroker {
	return m.zerodhaBroker
}

// GetAngelOneBroker returns the AngelOne order service
func (m *Module) GetAngelOneBroker() *service.AngelOneBroker {
	return m.angeloneBroker
}

// GetHDFCSecBroker returns the HDFC Securities order service
func (m *Module) GetHDFCSecBroker() *service.HDFCSecBroker {
	return m.hdfcsecBroker
}

// GetSharekhanBroker returns the Sharekhan order service
func (m *Module) GetSharekhanBroker() *service.SharekhanBroker {
	return m.sharekhanBroker
}

// GetGrowwBroker returns the Groww order service
func (m *Module) GetGrowwBroker() *service.GrowwBroker {
	return m.growwBroker
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