package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/model"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/notification"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/parser"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/settings"
)

var seenGuids = make(map[string]bool)

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

	var RSSData []model.RSS
	for _, rssConfiguration := range settings.RSSConfigurations {
		rss, err := parser.RSS(context, rssConfiguration.URL)
		if err != nil {
			log.Fatalf("Error parsing RSS feed: %v", err)
		}

		rss.Type = rssConfiguration.Type
		RSSData = append(RSSData, *rss)
	}

	for {
		for _, rss := range RSSData {
			for _, item := range rss.Channel.Items {
				if !seenGuids[item.GUID] {
					message := item.ToMessage()
					message.Type = rss.Type

					err := notificationClient.SendNotification(context, message)
					if err != nil {
						log.Printf("Error sending Lark notification: %v", err)
					} else {
						fmt.Printf("Notified: %s\n", item.Title)
						seenGuids[item.GUID] = true
					}
				}
			}
		}

		time.Sleep(time.Duration(settings.NotificationPoolIntervalMinutes) * time.Minute)
	}
}
