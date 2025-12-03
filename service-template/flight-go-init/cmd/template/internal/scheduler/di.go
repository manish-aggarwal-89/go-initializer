package scheduler

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Scheduler struct {
		dig.In
		SchedulerCredential      SchedulerCredential
		SchedulerSystemParameter SchedulerSystemParameter
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewSchedulerCredential); err != nil {
		return errors.Wrap(err, "failed to provide credential scheduler cron job")
	}
	if err := container.Provide(NewSystemParameterScheduler, dig.As(new(SchedulerSystemParameter))); err != nil {
		return errors.Wrap(err, "failed to provide system parameter scheduler cron job")
	}
	return nil
}
