package deps

import (
	"context"
	"github.com/sirupsen/logrus"
	redisuniversal "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/cache/redis-universal-with-metrics"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/errormapper/deps"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/commondeps"
	"go.uber.org/dig"
	"time"
	"{{MODULE_NAME}}/config"
)

type (
	Deps struct {
		commondeps.BaseDeps
		dig.In
		Config *config.Config
	}
)

func (impl *Deps) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := impl.RedisWrapper.Close()
	if err != nil {
		return err
	}
	impl.Kafka.Close()
	err = impl.Mongo.Client().Disconnect(ctx)
	if err != nil {
		return err
	}
	err = impl.Echo.Close()
	if err != nil {
		return err
	}

	return nil
}

func Register(container *dig.Container) error {
	if err := container.Provide(commondeps.NewRedisWrapper); err != nil {
		return err
	}
	if err := container.Provide(func(rw *commondeps.RedisWrapper) *redisuniversal.RedisUniversalClient {
		return rw.RedisIntegrator
	}); err != nil {
		return err
	}
	if err := container.Provide(commondeps.NewWorkerPoolWrapper); err != nil {
		return err
	}
	if err := container.Provide(commondeps.NewKafka); err != nil {
		return err
	}
	if err := container.Provide(commondeps.NewPredefined); err != nil {
		return err
	}
	if err := container.Provide(deps.NewIntegratorErrorMappingConfig); err != nil {
		return err
	}
	if err := container.Provide(deps.NewIntegratorErrorMappingDeps); err != nil {
		return err
	}

	return nil
}

func (impl *Deps) GetLogger(ctx context.Context) *logrus.Entry {
	return impl.BaseDeps.GetLogger(ctx)
}
