package service

import (
	"context"
	"errors"
	"time"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

type (
	HealthCheckServiceItf interface {
		HealthCheck(ctx context.Context, mr common.MandatoryRequest) error
	}
	healthCheckServiceImpl struct {
		deps deps.Deps
	}
)

func NewHealthCheckService(deps deps.Deps) (HealthCheckServiceItf, error) {
	return &healthCheckServiceImpl{deps}, nil
}

func (impl *healthCheckServiceImpl) HealthCheck(ctx context.Context, mr common.MandatoryRequest) error {
	ctxOperation := "health_check"
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// - check DB
	err := impl.deps.Mongo.Client().Ping(ctx, nil)
	if err != nil {
		impl.deps.GetLogger(ctx).Errorf("%s [%s] [%s] %s", shared.SERVICE_IMPL, ctxOperation, err.Error(), "failed to connect db mongo")
		return errors.New("failed to connect db mongo")
	}

	// - check Redis
	err = impl.deps.RedisWrapper.RedisIntegrator.PingWithContext(ctx)
	if err != nil {
		impl.deps.GetLogger(ctx).Errorf("%s [%s] [%s] %s", shared.SERVICE_IMPL, ctxOperation, err.Error(), "failed to connect db redis")
		return errors.New("failed to connect db redis")
	}
	return nil
}
