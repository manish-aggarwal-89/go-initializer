package di

import (
	"fmt"
	config2 "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/config"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/commondeps"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/di"
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/controller"
	"{{MODULE_NAME}}/internal/inbound"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/predefine"
	"{{MODULE_NAME}}/internal/repository"
	"{{MODULE_NAME}}/internal/scheduler"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential_log"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential_master_data"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/promotion"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/system_param"

	"go.uber.org/dig"
)

var Container = dig.New()

// registerSystemParam registers system_param from common lib with project config.
func registerSystemParam(container *dig.Container) error {
	return container.Invoke(func(cfg *config.Config, d deps.Deps) error {
		sysParamTtl := cfg.CacheConfig.SystemParameterTtlDuration
		if sysParamTtl == 0 {
			sysParamTtl = cfg.CacheConfig.DefaultTtlDuration
		}
		return system_param.Register(container,
			system_param.WithLogger(d.Logger),
			system_param.WithRedis(d.RedisWrapper.RedisIntegrator),
			system_param.WithMongoDB(d.Mongo),
			system_param.WithCacheConfig(
				cfg.GetSystemParameterCachePrefix(shared.DISTRIBUTION_TYPE),
				sysParamTtl,
				cfg.CacheConfig.DefaultTtlDuration,
			),
		)
	})
}

// registerCredential registers credential from common lib.
func registerCredential(container *dig.Container) error {
	return container.Invoke(func(cfg *config.Config, d deps.Deps) error {
		credTtl := cfg.CacheConfig.CredentialTtlDuration
		if credTtl == 0 {
			credTtl = cfg.CacheConfig.DefaultTtlDuration
		}
		return credential.Register(container,
			credential.WithLogger(d.Logger),
			credential.WithRedis(d.RedisWrapper.RedisIntegrator),
			credential.WithMongoDB(d.Mongo),
			credential.WithCacheConfig(
				cfg.GetCredentialCachePrefix(shared.DISTRIBUTION_TYPE),
				credTtl,
				cfg.CacheConfig.DefaultTtlDuration,
			),
			credential.WithIsStaging(cfg.IsStaging),
		)
	})
}

// wireCredentialLogAdapter wires the credential service with credential log adapter.
func wireCredentialLogAdapter(credSvc credential.Service, logSvc credential_log.Service) {
	credSvc.SetLogCreator(credential_log.NewCredentialLogAdapter(logSvc))
}

// registerCredentialLog registers credential_log from common lib.
func registerCredentialLog(container *dig.Container) error {
	return container.Invoke(func(cfg *config.Config, d deps.Deps) error {
		return credential_log.Register(container,
			credential_log.WithLogger(d.Logger),
			credential_log.WithMongoDB(d.Mongo),
		)
	})
}

// registerCredentialMasterData registers credential_master_data from common lib.
func registerCredentialMasterData(container *dig.Container) error {
	return container.Invoke(func(cfg *config.Config, d deps.Deps) error {
		return credential_master_data.Register(container,
			credential_master_data.WithLogger(d.Logger),
			credential_master_data.WithMongoDB(d.Mongo),
		)
	})
}

// registerPromotion registers promotion from common lib.
func registerPromotion(container *dig.Container) error {
	return container.Invoke(func(cfg *config.Config, d deps.Deps) error {
		dist := cfg.DistributionType
		if dist == "" {
			dist = "lion"
		}
		return promotion.Register(
			container,
			promotion.WithCachePrefix(cfg.PromotionConfig.CachePrefix),
			promotion.WithDistribution(dist),
			promotion.WithLogger(d.Logger),
			promotion.WithMongoDB(d.Mongo),
			promotion.WithListener(cfg.PromotionConfig.KafkaTopic, cfg.KafkaConfig.GeneralConfig.ReadTimeout, cfg.PromotionConfig.ListenerEnable),
			promotion.WithMetricClient(d.Metric),
			promotion.WithWorkerPool(d.WorkerPoolWrapper.KafkaGeneralPool),
			promotion.WithScheduler(cfg.PromotionConfig.SchedulerInterval, cfg.PromotionConfig.SchedulerEnable),
		)
	})
}

func init() {
	// - config
	if err := Container.Provide(config.New); err != nil {
		panic(fmt.Sprintf("failed to provide config %s", err))
	}

	if err := Container.Provide(func(cfg *config.Config) *config2.BaseConfig {
		return &cfg.BaseConfig
	}); err != nil {
		panic(fmt.Sprintf("failed to provide base config: %v", err))
	}

	if err := Container.Provide(di.NewMetric, dig.As(new(commondeps.MonitorStatsd))); err != nil {
		panic(fmt.Sprintf("failed to provide metric %s", err))
	}

	// - logger
	if err := Container.Provide(di.NewLogrusLogger); err != nil {
		panic(fmt.Sprintf("failed to provide logger %s", err))
	}
	// - logger using logrus
	if err := Container.Provide(di.NewSlogLogger); err != nil {
		panic(fmt.Sprintf("failed to provide logger v2 %s", err))
	}

	// - mongoDB
	if err := Container.Provide(di.NewMongoDB); err != nil {
		panic(fmt.Sprintf("failed to provide mongoDB %s", err))
	}

	// - echo
	if err := Container.Provide(di.NewEcho); err != nil {
		panic(fmt.Sprintf("failed to provide echo %s", err))
	}

	// - deps
	if err := deps.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register deps %s", err))
	}

	// - inbound
	if err := inbound.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register inbound %s", err))
	}

	// - controllers
	if err := controller.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register controller %s", err))
	}

	// - service
	if err := service.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register service %s", err))
	}

	// - repository
	if err := repository.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register repository %s", err))
	}

	// - scheduler
	if err := scheduler.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register scheduler %s", err))
	}

	// - predefine
	if err := predefine.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register predefine %s", err))
	}

	// - outbound
	if err := outbound.Register(Container); err != nil {
		panic(fmt.Sprintf("failed to register predefine %s", err))
	}

	// Common lib controllers (CredentialController, CredentialMasterDataController, SystemParamController, PromotionController)
	// are provided by these Register() functions.
	RegisterCommonControllers(Container)
}

func RegisterCommonControllers(container *dig.Container) {
	if err := registerSystemParam(container); err != nil {
		panic(fmt.Sprintf("failed to register system parameter service (common lib): %s", err))
	}
	if err := registerCredential(container); err != nil {
		panic(fmt.Sprintf("failed to register credential service (common lib): %s", err))
	}
	if err := registerCredentialLog(container); err != nil {
		panic(fmt.Sprintf("failed to register credential log service (common lib): %s", err))
	}
	if err := registerCredentialMasterData(container); err != nil {
		panic(fmt.Sprintf("failed to register credential master data service (common lib): %s", err))
	}
	if err := registerPromotion(container); err != nil {
		panic(fmt.Sprintf("failed to register promotion service (common lib): %s", err))
	}
	if err := container.Invoke(wireCredentialLogAdapter); err != nil {
		panic(fmt.Sprintf("failed to wire credential log adapter: %s", err))
	}
}
