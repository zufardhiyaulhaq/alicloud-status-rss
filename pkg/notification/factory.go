package notification

import (
	"fmt"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/notification/lark"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/settings"
)

func NotificationFactory(settings settings.Settings) (NotificationInterface, error) {
	if settings.NotificationType == "lark" {
		return lark.NewLarkClient(settings.NotificationLarkWebhookURL), nil
	}

	return nil, fmt.Errorf("notification is not supported")
}
