package settings

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	NotificationType                    string                   `envconfig:"NOTIFICATION_TYPE" default:"lark"`
	NotificationLarkWebhookURLs         []string                 `envconfig:"NOTIFICATION_LARK_WEBHOOK_URLS"`
	NotificationPoolIntervalMinutes     int                      `envconfig:"NOTIFICATION_POOL_INTERVAL_MINUTES" default:"5"`
	RSSConfigurations                   RSSConfigurationsDecoder `envconfig:"RSS_CONFIGURATIONS" required:"true"`
	RSSGUIDExclusionList                []string                 `envconfig:"RSS_GUID_EXCLUSION_LIST" default:""`
	RSSGUIDCacheEnabled                 bool                     `envconfig:"RSS_GUID_CACHE_ENABLED" default:"false"`
	RSSGUIDCacheType                    string                   `envconfig:"RSS_GUID_CACHE_TYPE" default:"redis"`
	RSSGUIDCacheKey                     string                   `envconfig:"RSS_GUID_CACHE_KEY" default:"rss-guids"`
	RSSGUIDCacheRedisType               string                   `envconfig:"RSS_GUID_CACHE_REDIS_TYPE" default:"standalone"`
	RSSGUIDCacheRedisAddress            string                   `envconfig:"RSS_GUID_CACHE_REDIS_ADDRESS"`
	RSSGUIDCacheRedisPassword           string                   `envconfig:"RSS_GUID_CACHE_REDIS_PASSWORD"`
	RSSGUIDCacheRedisSentinelAddress    []string                 `envconfig:"RSS_GUID_CACHE_REDIS_SENTINEL_ADDRESS"`
	RSSGUIDCacheRedisSentinelMasterName string                   `envconfig:"RSS_GUID_CACHE_REDIS_SENTINEL_MASTER_NAME"`
	RSSGUIDCacheRedisSentinelPassword   string                   `envconfig:"RSS_GUID_CACHE_REDIS_SENTINEL_PASSWORD"`
}

type RSSConfiguration struct {
	Type   string
	URL    string
	Region string
}

type RSSConfigurationsDecoder []RSSConfiguration

func (rcd *RSSConfigurationsDecoder) Decode(value string) error {
	rssConfigs := []RSSConfiguration{}
	pairs := strings.Split(value, ";")
	for _, pair := range pairs {
		kvpair := strings.Split(pair, ",")
		if len(kvpair) != 3 {
			return fmt.Errorf("invalid map item: %q", pair)
		}
		config := RSSConfiguration{
			Type:   kvpair[0],
			URL:    kvpair[1],
			Region: kvpair[2],
		}
		rssConfigs = append(rssConfigs, config)
	}
	*rcd = RSSConfigurationsDecoder(rssConfigs)
	return nil
}

func (s Settings) Validator(ctx context.Context) error {
	if s.NotificationType == "lark" {
		if len(s.NotificationLarkWebhookURLs) == 0 {
			return errors.New("NOTIFICATION_LARK_WEBHOOK_URLS is required when using lark notification type")
		}

		return nil
	}

	if s.RSSGUIDCacheEnabled {
		if s.RSSGUIDCacheType == "redis" {
			if s.RSSGUIDCacheRedisType == "sentinel" {
				if len(s.RSSGUIDCacheRedisSentinelAddress) == 0 {
					return errors.New("RSS_GUID_CACHE_REDIS_SENTINEL_ADDRESS is required when using redis sentinel cache type")
				}
				if s.RSSGUIDCacheRedisSentinelMasterName == "" {
					return errors.New("RSS_GUID_CACHE_REDIS_SENTINEL_MASTER_NAME is required when using redis sentinel cache type")
				}
			}
		}
	}

	return nil
}

func NewSettings(ctx context.Context) (Settings, error) {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		return settings, err
	}

	if err := settings.Validator(ctx); err != nil {
		return settings, err
	}

	return settings, nil
}
