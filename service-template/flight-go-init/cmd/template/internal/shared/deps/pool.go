package deps

import (
	"{{MODULE_NAME}}/config"

	"github.com/sourcegraph/conc/pool"
)

type (
	WorkerPoolWrapper struct {
		SearchRequestPool *pool.Pool
		KafkaGeneralPool  *pool.Pool
	}
)

func NewWorkerPoolWrapper(cfg *config.Config) (*WorkerPoolWrapper, error) {
	var (
		p1 = pool.New().WithMaxGoroutines(cfg.KafkaConfig.FlightConfig.CorePool)
		p2 = pool.New().WithMaxGoroutines(cfg.KafkaConfig.GeneralConfig.CorePool)
	)

	return &WorkerPoolWrapper{
		SearchRequestPool: p1,
		KafkaGeneralPool:  p2,
	}, nil
}
