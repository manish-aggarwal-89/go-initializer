package deps

import (
	"context"
	"fmt"
	"time"
	"{{MODULE_NAME}}/config"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache"
	redis "github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache/redis-universal-with-metrics"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

type (
	RedisWrapper struct {
		RedisIntegrator cache.CacheWithContext
	}
)

func (impl *RedisWrapper) Close() error {
	ctx, ccl := context.WithTimeout(context.Background(), 3*time.Second)
	defer ccl()

	err := impl.RedisIntegrator.CloseWithContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewRedisWrapper(cfg *config.Config, metric metrics.MonitorStatsd) (*RedisWrapper, error) {
	redisIntegrator, err := newRedisIntegrator(cfg, metric)
	if err != nil {
		return nil, err
	}

	return &RedisWrapper{
		RedisIntegrator: redisIntegrator,
	}, nil
}

func newRedisIntegrator(cfg *config.Config, metric metrics.MonitorStatsd) (cache.CacheWithContext, error) {
	redisCfg := cfg.RedisConfig.IntegratorRedisConfig
	opts := redis.Option{
		DB:           redisCfg.Database,
		Address:      redisCfg.Address,
		Password:     redisCfg.Password,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConn,
		DialTimeout:  redisCfg.DialTimeout,
		PoolTimeout:  redisCfg.PoolTimeout,
		ReadTimeout:  redisCfg.ReadTimeout,
		WriteTimeout: redisCfg.WriteTimeout,
		MaxConnAge:   redisCfg.MaxConnAge,
		ReadOnly:     false,
	}

	ow, err := redis.NewWithMetric(&opts, metric, cfg.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to provide redis integrator %s", err)
	}

	return ow, nil
}
