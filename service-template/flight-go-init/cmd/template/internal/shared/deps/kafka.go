package deps

import (
	"{{MODULE_NAME}}/config"

	"github.com/labstack/gommon/log"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
	ks "github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging/kafka-sarama"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

type (
	KafkaWrapper struct {
		KafkaFlight  messaging.QueueV2
		KafkaGeneral messaging.QueueV2
	}
)

func (impl KafkaWrapper) Close() {
	err := impl.KafkaFlight.Close()
	if err != nil {
		log.Errorf("kafka flight close err: %s", err.Error())
	}

	err = impl.KafkaGeneral.Close()
	if err != nil {
		log.Errorf("kafka general close err: %s", err.Error())
	}
}

func NewKafka(cfg *config.Config, logger logs.Logger, metric metrics.MonitorStatsd) (*KafkaWrapper, error) {
	kafkaFlight, err := newKafkaFlight(cfg, logger, metric)
	if err != nil {
		return nil, err
	}

	kafkaGeneral, err := newKafkaGeneral(cfg, logger, metric)
	if err != nil {
		return nil, err
	}

	return &KafkaWrapper{
		KafkaFlight:  kafkaFlight,
		KafkaGeneral: kafkaGeneral,
	}, nil
}

func newKafkaGeneral(cfg *config.Config, logger logs.Logger, metric metrics.MonitorStatsd) (messaging.QueueV2, error) {
	generalCfg := cfg.KafkaConfig.GeneralConfig
	client, err := ks.New(&ks.Option{
		KafkaVersion:  generalCfg.Version,
		Host:          generalCfg.Host,
		ConsumerGroup: generalCfg.ConsumerGroup,
		Log:           logger,
		CaptureMetric: generalCfg.CaptureMetric,
		MetricData: ks.MetricData{
			Monitor:     metric,
			ServiceName: cfg.ServiceName,
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newKafkaFlight(cfg *config.Config, logger logs.Logger, metric metrics.MonitorStatsd) (messaging.QueueV2, error) {
	flightCfg := cfg.KafkaConfig.FlightConfig
	client, err := ks.New(&ks.Option{
		KafkaVersion:  flightCfg.Version,
		Host:          flightCfg.Host,
		ConsumerGroup: flightCfg.ConsumerGroup,
		Log:           logger,
		CaptureMetric: flightCfg.CaptureMetric,
		MetricData: ks.MetricData{
			Monitor:     metric,
			ServiceName: cfg.ServiceName,
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
