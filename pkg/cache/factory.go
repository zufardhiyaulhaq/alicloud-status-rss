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
			return redis.RedisSentinelClient(ctx, settings), nil
		}
	}

	return nil, fmt.Errorf("notification is not supported")
}
