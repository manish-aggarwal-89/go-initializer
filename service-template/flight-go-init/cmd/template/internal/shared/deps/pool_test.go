package deps

import (
	"testing"
	"{{MODULE_NAME}}/config"
)

func TestNewWorkerPoolWrapper(t *testing.T) {
	// Test case 1: default configuration
	cfg := &config.Config{
		KafkaConfig: config.KafkaConfig{
			FlightConfig: config.KafkaServerConfig{
				CorePool: 1,
			},
			GeneralConfig: config.KafkaServerConfig{
				CorePool: 1,
			},
		},
	}
	wrapper, err := NewWorkerPoolWrapper(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if wrapper.SearchRequestPool.MaxGoroutines() != 1 {
		t.Errorf("unexpected search request pool max goroutines: %d", wrapper.SearchRequestPool.MaxGoroutines())
	}
	if wrapper.KafkaGeneralPool.MaxGoroutines() != 1 {
		t.Errorf("unexpected kafka general pool max goroutines: %d", wrapper.KafkaGeneralPool.MaxGoroutines())
	}

	// Test case 2: custom configuration
	cfg = &config.Config{
		KafkaConfig: config.KafkaConfig{
			FlightConfig: config.KafkaServerConfig{
				CorePool: 1,
			},
			GeneralConfig: config.KafkaServerConfig{
				CorePool: 1,
			},
		},
	}
	wrapper, err = NewWorkerPoolWrapper(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if wrapper.SearchRequestPool.MaxGoroutines() != 1 {
		t.Errorf("unexpected search request pool max goroutines: %d", wrapper.SearchRequestPool.MaxGoroutines())
	}
	if wrapper.KafkaGeneralPool.MaxGoroutines() != 1 {
		t.Errorf("unexpected kafka general pool max goroutines: %d", wrapper.KafkaGeneralPool.MaxGoroutines())
	}
}
