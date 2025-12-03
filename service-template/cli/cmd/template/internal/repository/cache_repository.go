package repository

import (
	"context"
	"time"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type CacheRepositoryImpl struct {
	client *deps.RedisWrapper
}

type CacheRepository interface {
	GetWithContext(ctx context.Context, key string, data interface{}) error
	IsExistsWithContext(ctx context.Context, key string) error
	SetWithContext(ctx context.Context, key string, data interface{}) error
	SetWithExpirationWithContext(ctx context.Context, key string, data interface{}, ttl time.Duration) error
	RemoveWithContext(ctx context.Context, key string) error
}

func NewCacheRepository(client *deps.RedisWrapper) *CacheRepositoryImpl {
	return &CacheRepositoryImpl{
		client: client,
	}
}

func (r *CacheRepositoryImpl) GetWithContext(ctx context.Context, key string, data interface{}) error {
	return r.client.RedisIntegrator.GetWithContext(ctx, key, data)
}

func (r *CacheRepositoryImpl) IsExistsWithContext(ctx context.Context, key string) error {
	return r.client.RedisIntegrator.Exists(ctx, key)
}

func (r *CacheRepositoryImpl) SetWithContext(ctx context.Context, key string, data interface{}) error {
	return r.client.RedisIntegrator.SetWithContext(ctx, key, data)
}

func (r *CacheRepositoryImpl) SetWithExpirationWithContext(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	return r.client.RedisIntegrator.SetWithExpirationWithContext(ctx, key, data, ttl)
}

func (r *CacheRepositoryImpl) RemoveWithContext(ctx context.Context, key string) error {
	return r.client.RedisIntegrator.RemoveWithContext(ctx, key)
}
