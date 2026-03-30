package predefine

import (
	"github.com/pkg/errors"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/di"
	"go.uber.org/dig"
)

type (
	Scheduler struct {
		dig.In
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(di.NewPredefined); err != nil {
		return errors.Wrap(err, "failed to provide PredefinedImpl")
	}

	return nil
}
