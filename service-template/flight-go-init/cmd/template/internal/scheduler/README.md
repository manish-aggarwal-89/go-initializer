**{{SERVICE_TITLE}}**

This layer use for creating the scheduler of project. This layer need add di.go to register the dependency so can runing
by start-up system

Sample di.go

```go
package scheduler

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type (
	Scheduler struct {
		dig.In
		Scheduler scheduler
	}
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewScheduler); err != nil {
		return errors.Wrap(err, "failed to provide Scheduler cron job")
	}

	return nil
}
```

The scheduler we are using the library [robfig](github.com/robfig/cron/v3)

Sample of scheduler

```go
package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type (
	scheduler interface {
		SystemParameterPredifine() error
	}

	SchedulerImpl struct {
		deps        deps.Deps
		repo        repository.SystemParameterInterface
		predefinded predefine.Predefined
	}
)

func NewScheduler(deps deps.Deps, repo repository.SystemParameterInterface, predefinded predefine.Predefined) scheduler {
	return &SchedulerImpl{deps: deps, repo: repo, predefinded: predefinded}
}

func (impl *SchedulerImpl) SystemParameterPredifine() error {
	// ctx_o := "[system_parameter_predifine]"

	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime), cron.WithSeconds())

	// stop scheduler before action is terminate
	defer scheduler.Stop()

	scheduler.AddFunc("*/"+impl.deps.Config.PredefinedSchedulerRunning+" * * * * *", impl.SystemParameter)

	// start scheduler
	go scheduler.Start()

	return nil
}

func (impl *SchedulerImpl) SystemParameter() {
	impl.deps.GetLogger(nil).Info("[scheduler_predefined][system_parameter] - refresh system parameter predefine")
	// - we set run max of create predefined only 5 second
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, _ := impl.repo.FindAll(ctx)

	for _, sysParam := range data {
		impl.predefinded.SetData(ctx, sysParam.Variable, sysParam, time.Duration(impl.deps.Config.PredefinedSchedulerSystemParameterTtl))
	}
}

```