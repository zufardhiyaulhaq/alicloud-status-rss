package cache

import (
	"context"
	"fmt"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/cache/redis"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/settings"
)

func CacheFactory(ctx context.Context, settings settings.Settings) (CacheInterface, error) {
	if settings.RSSGUIDCacheType == "redis" {
		if settings.RSSGUIDCacheRedisType == "sentinel" {
			client, err := redis.RedisSentinelClient(ctx, settings)
			if err != nil {
				return nil, err
			}

			return client, nil
		}

		if settings.RSSGUIDCacheRedisType == "standalone" {
			client, err := redis.RedisClient(ctx, settings)
			if err != nil {
				return nil, err
			}

			return client, nil
		}
	}

	return nil, fmt.Errorf("cache configuration is not supported")
}
