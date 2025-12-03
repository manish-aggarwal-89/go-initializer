package inbound

import (
	"{{MODULE_NAME}}/internal/shared/deps"

	"go.uber.org/dig"
)

type (
	Listeners struct {
		dig.In
		Deps deps.Deps
	}
)

func (impl *Listeners) Listen() {
	// - register listener here.
	// example:
	// impl.Deps.Kafka.KafkaFlight.AddTopicListener(cfg.KafkaConfig.Topics.IntegratorSearchRequest, impl.IntegratorListener.ListenToIntegratorRequest)
	// impl.Deps.Kafka.KafkaFlight.Listen()
}

func Register(container *dig.Container) error {

	return nil
}
