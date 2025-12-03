package scheduler

import (
	"context"
	"time"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/robfig/cron/v3"
)

type (
	SchedulerCredential interface {
		CredentialPredefine() error
	}

	SchedulerCredentialImpl struct {
		dep        deps.Deps
		repo       repository.CredentialRepositoryInterface
		predefined predefine.Predefined
	}
)

func NewSchedulerCredential(deps deps.Deps, repo repository.CredentialRepositoryInterface,
	predefined predefine.Predefined) SchedulerCredential {
	return &SchedulerCredentialImpl{dep: deps, repo: repo, predefined: predefined}
}

func (impl SchedulerCredentialImpl) CredentialPredefine() error {
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime), cron.WithSeconds())

	defer scheduler.Stop()

	secs := impl.dep.Config.SchedulerConfig.CredentialIntervalSecs

	impl.dep.GetLogger(nil).Infof("[scheduler_predefined][credential] - will run every %s secs", secs)

	_, err := scheduler.AddFunc("@every "+secs+"s", impl.PopulateCredentials)
	if err != nil {
		return err
	}
	impl.PopulateCredentials()

	go scheduler.Start()

	return nil
}

func (impl SchedulerCredentialImpl) PopulateCredentials() {
	impl.dep.GetLogger(nil).Info("[scheduler_predefined][credential] - running saving data")
	// - we set run max of create predefined only 5 second
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, _ := impl.repo.FindAll(ctx)

	for _, credential := range data {
		if err := impl.predefined.SetData(ctx, constructKey(credential.Supplier), credential,
			impl.dep.Config.CacheConfig.DefaultTtlDuration); err != nil {
			impl.dep.GetLogger(ctx).Errorf("[scheduler_predefined][credential] - %s", err.Error())
		}
	}
}

func constructKey(supplier string) string {
	return shared.PREDEFINED_CREDENTIAL_CACHE_KEY_PREFIX + supplier
}
