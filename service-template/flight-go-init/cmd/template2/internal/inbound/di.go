package inbound

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type Listeners struct {
	dig.In
	Deps               deps.Deps
	IntegratorListener IntegratorListener
}

func (impl *Listeners) Listen() {
	cfg := impl.Deps.Config
	kafkaFlight := impl.Deps.Kafka.KafkaFlight
	topic := cfg.KafkaConfig.Topics.IntegratorSearchRequest
	if topic != "" {
		kafkaFlight.AddTopicListener(topic, impl.IntegratorListener.ListenToIntegratorRequest)
	}
	kafkaFlight.Listen()
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewIntegratorListener); err != nil {
		return errors.Wrap(err, "failed to provide IntegratorListener")
	}
	return nil
}
