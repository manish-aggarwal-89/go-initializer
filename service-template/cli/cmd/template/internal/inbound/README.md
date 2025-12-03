**{{SERVICE_TITLE}}**

# Table Of Content

- [Table Of Content](#table-of-content)
    - [Inbound](#inbound)
    - [Sample DI](#sample-di)

## Inbound

What is inbound? inbound is our function on system will handle incomming inbound layer on oursystem.

Example we need to add inbound process to listen data on kafka. We need to declare function on this package

## Sample DI

We need add di.go to register all inboud on our process, sample :

```go
package inbound

import (
  "go.uber.org/dig"
)

type (
  Listeners struct {
    dig.In
    Deps               deps.Deps
    IntegratorListener IntegratorListener
    CurrencyListener   CurrencyV1Listener
  }
)

func (impl *Listeners) Listen() {
  var (
    cfg          = impl.Deps.Config
    kafkaFlight  = impl.Deps.Kafka.KafkaFlight
    kafkaGeneral = impl.Deps.Kafka.KafkaGeneral
  )

  // - register listener to topic in kafka search request
  kafkaFlight.AddTopicListener(cfg.TopicIntegratorSearchRequest, impl.IntegratorListener.ListenToIntegratorRequest)
  kafkaFlight.Listen()

  // - register listener to topic in kafka search request
  kafkaGeneral.AddTopicListener(cfg.TopicIntegratorMultiCurrencyV1, impl.CurrencyListener.ListenMultiCurrencyRequest)
  kafkaGeneral.Listen()

}

func Register(container *dig.Container) error {

  if err := container.Provide(NewIntegratorListener); err != nil {
    return errors.Wrap(err, "failed to provide IntegratorListener")
  }

  if err := container.Provide(NewCurrencyV1Listener); err != nil {
    return errors.Wrap(err, "failed to provide KafkaGeneralListener")
  }

  return nil
}
```

we can inject inbound on DI to actived the listener. On the sample we have 2 kind of inbound (Search Listener and
Currency Listener)
