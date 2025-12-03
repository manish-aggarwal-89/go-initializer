package deps

import (
	"{{MODULE_NAME}}/config"

	"github.com/patrickmn/go-cache"
)

type (
	PredefinedCache struct {
		Predefined cache.Cache
	}
)

func (impl *PredefinedCache) Close() error {
	impl.Predefined.Flush()
	return nil
}

func NewPredefined(cfg *config.Config) (*PredefinedCache, error) {
	cacheClient := cache.New(cfg.CacheConfig.DefaultTtlDuration, cfg.CacheConfig.DefaultCleanupIntervalDuration)

	return &PredefinedCache{
		Predefined: *cacheClient,
	}, nil
}
