package predefine

import (
	"context"
	"errors"
	"time"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type (
	Predefined interface {
		SetData(ctx context.Context, key string, data interface{}, ttl time.Duration) error
		GetDataByID(ctx context.Context, key string) (interface{}, error)
		GetAll(ctx context.Context) map[string]interface{}
		FlushData(ctx context.Context) error
	}

	PredefinedImpl struct {
		dep deps.Deps
	}
)

func NewPredefined(deps deps.Deps) Predefined {
	return &PredefinedImpl{dep: deps}
}

// FlushData implements Predefined.
func (impl *PredefinedImpl) FlushData(_ context.Context) error {
	impl.dep.PredefinedCache.Predefined.Flush()
	return nil
}

// GetDataByID implements Predefined.
func (impl *PredefinedImpl) GetDataByID(_ context.Context, key string) (interface{}, error) {
	data, ok := impl.dep.PredefinedCache.Predefined.Get(key)
	if !ok {
		return nil, errors.New("key " + key + " is not exist!")
	}
	return data, nil
}

func (impl *PredefinedImpl) GetAll(_ context.Context) map[string]interface{} {
	items := impl.dep.PredefinedCache.Predefined.Items()

	result := make(map[string]interface{})
	for key, item := range items {
		result[key] = item.Object
	}

	return result
}

// SetData implements Predefined.
func (impl *PredefinedImpl) SetData(_ context.Context, key string, data interface{}, ttl time.Duration) error {
	impl.dep.PredefinedCache.Predefined.Set(key, data, ttl*time.Second)
	return nil
}
