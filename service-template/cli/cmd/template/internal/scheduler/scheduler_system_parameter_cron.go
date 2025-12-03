package scheduler

import (
	"context"
	"time"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/robfig/cron/v3"
)

type (
	SchedulerSystemParameter interface {
		RefreshSystemParameterPredifine() error
	}

	SchedulerSystemParameterImpl struct {
		deps        deps.Deps
		repo        repository.SystemParameterInterface
		predefinded predefine.Predefined
	}
)

func NewSystemParameterScheduler(deps deps.Deps, repo repository.SystemParameterInterface, predefinded predefine.Predefined) SchedulerSystemParameter {
	return &SchedulerSystemParameterImpl{deps: deps, repo: repo, predefinded: predefinded}
}

func (impl *SchedulerSystemParameterImpl) RefreshSystemParameterPredifine() error {
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime), cron.WithSeconds())

	// stop scheduler before action is terminate
	defer scheduler.Stop()

	secs := impl.deps.Config.SchedulerConfig.SystemParameterIntervalSecs

	impl.deps.GetLogger(nil).Infof("[scheduler_predefined][system_parameter] - will run every %s secs", secs)

	_, err := scheduler.AddFunc("@every "+secs+"s", impl.SystemParameter)
	if err != nil {
		return err
	}
	impl.SystemParameter()

	// start scheduler
	go scheduler.Start()

	return nil
}

func (impl *SchedulerSystemParameterImpl) SystemParameter() {
	impl.deps.GetLogger(nil).Info("[scheduler_predefined][system_parameter] - refresh system parameter predefine")
	// - we set run max of create predefined only 5 second
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, _ := impl.repo.FindAll(ctx)

	for _, sysParam := range data {
		impl.predefinded.SetData(ctx, sysParam.Variable, sysParam, time.Duration(impl.deps.Config.CacheConfig.DefaultTtlDuration))
	}
}
