package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/settings"
)

type RedisSentinelCache struct {
	ctx    context.Context
	client *redis.Client
}

func (r RedisSentinelCache) Set(ctx context.Context, key string, value string) error {
	status := r.client.Set(ctx, key, value, 24*time.Hour)

	return status.Err()
}

func (r RedisSentinelCache) Get(ctx context.Context, key string) (string, error) {
	out := r.client.Get(ctx, key)

	if out.Err() != nil {
		return "", out.Err()
	}

	return out.Val(), nil
}

func RedisSentinelClient(ctx context.Context, settings settings.Settings) (RedisSentinelCache, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: settings.RSSGUIDCacheRedisSentinelAddress,
		MasterName:    settings.RSSGUIDCacheRedisSentinelMasterName,
		Password:      settings.RSSGUIDCacheRedisPassword,
	})

	return RedisSentinelCache{
		ctx:    ctx,
		client: client,
	}, nil
}
