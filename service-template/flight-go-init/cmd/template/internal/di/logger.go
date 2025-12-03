package di

import (
	"{{MODULE_NAME}}/config"

	"github.com/sirupsen/logrus"
	commonlogrus "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
)

func NewLogger(cfg *config.Config) (logs.Logger, error) {
	lg, err := logs.New(&logs.Option{
		Level:     logs.Level(cfg.LogConfig.Level),
		Formatter: logs.Formatter(cfg.LogConfig.Formatter),
	})

	if err != nil {
		return nil, err
	}

	return lg, nil
}

func NewLoggerV2(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(commonlogrus.ParseLogLevel(cfg.LogConfig.Level))
	logger.SetFormatter(commonlogrus.GetFormatter(cfg.LogConfig.Formatter))
	logger.AddHook(&commonlogrus.ContextHook{})
	return logger
}
