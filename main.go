package main

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/cache"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/data"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/model"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/notification"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/settings"
)

func main() {
	context := context.Background()

	settings, err := settings.NewSettings(context)
	if err != nil {
		log.Fatalf("Error loading settings: %v", err)
	}

	notificationClient, err := notification.NotificationFactory(settings)
	if err != nil {
		log.Fatalf("Error creating notification client: %v", err)
	}

	var cacheClient cache.CacheInterface
	if settings.RSSGUIDCacheEnabled {
		cacheClient, err = cache.CacheFactory(context, settings)
		if err != nil {
			log.Fatalf("Error creating cache client: %v", err)
		}
	}

	var RSSData []model.RSS
	for _, rssConfiguration := range settings.RSSConfigurations {
		rss, err := data.ParseRSS(context, rssConfiguration.URL)
		if err != nil {
			log.Fatalf("Error parsing RSS feed: %v", err)
		}

		rss.Type = rssConfiguration.Type
		RSSData = append(RSSData, *rss)
	}

	processedRSS, err := data.ProcessRSS(RSSData)
	if err != nil {
		log.Fatalf("Error processing RSS data: %v", err)
	}

	var seenGuids []string
	seenGuids = append(seenGuids, settings.RSSGUIDExclusionList...)

	for {
		if settings.RSSGUIDCacheEnabled {
			rawSeenGuids, err := cacheClient.Get(context, settings.RSSGUIDCacheKey)
			if err != nil {
				log.Printf("Error retrieving seen GUIDs from cache: %v", err)
			}
			seenGuids = append(seenGuids, strings.Split(rawSeenGuids, ",")...)
		}

		for _, item := range processedRSS {
			if !slices.Contains(seenGuids, item.GUID) {
				message := item.ToMessage()

				err := notificationClient.SendNotification(context, message)
				if err != nil {
					log.Printf("Error sending Lark notification: %v", err)
				} else {
					fmt.Printf("Notified: %s\n", item.Title)
					seenGuids = append(seenGuids, item.GUID)
				}
			} else {
				fmt.Printf("Already seen: %s\n", item.Title)
			}
		}

		if settings.RSSGUIDCacheEnabled {
			err = cacheClient.Set(context, settings.RSSGUIDCacheKey, strings.Join(seenGuids, ","))
			if err != nil {
				log.Printf("Error saving seen GUIDs to cache: %v", err)
			}
		}

		time.Sleep(time.Duration(settings.NotificationPoolIntervalMinutes) * time.Minute)
	}
}
