package di

import (
	"{{MODULE_NAME}}/config"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

func NewMetric(cfg *config.Config) (metrics.MonitorStatsd, error) {
	return metrics.NewMonitor(cfg.StatsDConfig.Addr, cfg.StatsDConfig.Port, cfg.ServiceName)
}
