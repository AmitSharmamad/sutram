package dvendor

import (
	env "algo/src/core/config"

	"github.com/google/wire"
)

// VendorService is responsible for providing configured data vendors
type VendorService struct {
	factory *DataVendorFactory
	envProv *env.EnvProv
}

// NewVendorService creates a new vendor service with the given factory and environment provider
func NewVendorService(factory *DataVendorFactory, envProv *env.EnvProv) *VendorService {
	return &VendorService{
		factory: factory,
		envProv: envProv,
	}
}

// GetDhanVendor returns a configured Dhan data vendor
func (s *VendorService) GetDhanVendor() (DataVendor, error) {
	config := VendorConfig{
		Type:     DhanVendor,
		ClientID: s.envProv.GetProp(env.DHAN_CLIENT_ID),
		Token:    s.envProv.GetProp(env.DHAN_ACCESS_TOKEN),
		BaseURL:  s.envProv.GetProp(env.DHAN_BASE_URL),
		Parameters: map[string]string{
			"dhanClientID": s.envProv.GetProp(env.DHAN_DHAN_CLIENT_ID),
		},
	}

	vendor, err := s.factory.CreateVendor(config)
	if err != nil {
		return nil, err
	}

	// We don't connect here to avoid network operations during dependency injection
	// The consumer should call Connect() when ready
	return vendor, nil
}

// GetZerodhaVendor returns a configured Zerodha data vendor
func (s *VendorService) GetZerodhaVendor() (DataVendor, error) {
	config := VendorConfig{
		Type:       ZerodhaVendor,
		ClientID:   s.envProv.GetProp(env.ZERODHA_CLIENT_ID),
		Token:      s.envProv.GetProp(env.ZERODHA_ACCESS_TOKEN),
		BaseURL:    s.envProv.GetProp(env.ZERODHA_BASE_URL),
		Parameters: map[string]string{},
	}

	vendor, err := s.factory.CreateVendor(config)
	if err != nil {
		return nil, err
	}

	// We don't connect here to avoid network operations during dependency injection
	// The consumer should call Connect() when ready
	return vendor, nil
}

// GetDefaultVendor returns the default configured data vendor (currently Dhan)
func (s *VendorService) GetDefaultVendor() (DataVendor, error) {
	return s.GetDhanVendor()
}

// ProvideDataVendorFactory is a Wire provider for DataVendorFactory
func ProvideDataVendorFactory() *DataVendorFactory {
	return NewDataVendorFactory()
}

// ProvideVendorService is a Wire provider for VendorService
func ProvideVendorService(factory *DataVendorFactory, envProv *env.EnvProv) *VendorService {
	return NewVendorService(factory, envProv)
}

// ProvideDefaultDataVendor is a Wire provider for the default DataVendor
func ProvideDefaultDataVendor(service *VendorService) (DataVendor, error) {
	return service.GetDefaultVendor()
}

// Set provides the Wire provider set for d-vendor components
var Set = wire.NewSet(
	ProvideDataVendorFactory,
	ProvideVendorService,
	ProvideDefaultDataVendor,
)
