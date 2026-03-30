package scheduler

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type Scheduler struct {
	dig.In
	SchedulerCredential      SchedulerCredential
	SchedulerSystemParameter SchedulerSystemParameter
}

type SchedulerCredential interface {
	CredentialPredefine() error
}

type SchedulerSystemParameter interface {
	RefreshSystemParameterPredifine() error
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewNoopSchedulerCredential); err != nil {
		return errors.Wrap(err, "failed to provide credential scheduler")
	}
	if err := container.Provide(NewNoopSchedulerSystemParameter); err != nil {
		return errors.Wrap(err, "failed to provide system parameter scheduler")
	}
	return nil
}
