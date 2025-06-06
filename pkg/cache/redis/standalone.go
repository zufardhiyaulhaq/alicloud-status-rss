package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/settings"
)

type RedisCache struct {
	ctx    context.Context
	client *redis.Client
}

func (r RedisCache) Set(ctx context.Context, key string, value string) error {
	status := r.client.Set(ctx, key, value, 24*time.Hour)

	return status.Err()
}

func (r RedisCache) Get(ctx context.Context, key string) (string, error) {
	out := r.client.Get(ctx, key)

	if out.Err() != nil {
		return "", out.Err()
	}

	return out.Val(), nil
}

func RedisClient(ctx context.Context, settings settings.Settings) (RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     settings.RSSGUIDCacheRedisAddress,
		Password: settings.RSSGUIDCacheRedisPassword,
	})

	return RedisCache{
		ctx:    ctx,
		client: client,
	}, nil
}
