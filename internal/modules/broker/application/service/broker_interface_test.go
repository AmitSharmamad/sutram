package service

import (
	"testing"

	"algo/internal/modules/broker/domain"
)

// TestAllBrokersImplementInterface verifies that all broker services implement the Broker interface
func TestAllBrokersImplementInterface(t *testing.T) {
	// This test will compile only if all services implement the Broker interface
	var brokers []domain.Broker
	
	// Test that each service can be assigned to the interface
	brokers = append(brokers, NewDhanBroker("", ""))
	brokers = append(brokers, NewZerodhaBroker("", "", ""))
	brokers = append(brokers, NewSharekhanBroker("", "", ""))
	brokers = append(brokers, NewHDFCSecBroker("", ""))
	brokers = append(brokers, NewAngelOneBroker("", "", ""))
	brokers = append(brokers, NewGrowwBroker("", "", ""))
	
	// Verify we have all expected brokers
	expectedCount := 6
	if len(brokers) != expectedCount {
		t.Errorf("Expected %d brokers, got %d", expectedCount, len(brokers))
	}
	
	// Verify each broker implements the interface
	for i, broker := range brokers {
		if broker == nil {
			t.Errorf("Broker %d is nil", i)
		}
	}
}