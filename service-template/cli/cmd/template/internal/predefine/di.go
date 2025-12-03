package predefine

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Scheduler struct {
		dig.In
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewPredefined); err != nil {
		return errors.Wrap(err, "failed to provide PredefinedImpl")
	}

	return nil
}
