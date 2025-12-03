package outbound

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Outbound struct {
		dig.In
		Impl OutboundImpl
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewOutboundImpl, dig.As(new(OutboundItf))); err != nil {
		return errors.Wrap(err, "failed to provide OutboundImpl")
	}
	return nil
}
