package deps

import (
	"context"
	"time"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/shared"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

type (
	Deps struct {
		dig.In
		Config            *config.Config
		RedisWrapper      *RedisWrapper
		Kafka             *KafkaWrapper
		WorkerPoolWrapper *WorkerPoolWrapper
		Metric            metrics.MonitorStatsd
		Echo              *echo.Echo
		Mongo             *mongo.Database
		PredefinedCache   *PredefinedCache
		Logger            *logrus.Logger
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
	if err := container.Provide(NewRedisWrapper); err != nil {
		return err
	}
	if err := container.Provide(NewWorkerPoolWrapper); err != nil {
		return err
	}
	if err := container.Provide(NewKafka); err != nil {
		return err
	}
	if err := container.Provide(NewPredefined); err != nil {
		return err
	}

	return nil
}

func (impl *Deps) GetLogger(ctx context.Context) *logrus.Entry {

	if ctx == nil {
		return impl.Logger.WithContext(context.Background())
	}
	if logger, ok := ctx.Value(shared.LOGGER).(*logrus.Entry); ok {
		return logger
	}
	return impl.Logger.WithContext(context.Background())
}
