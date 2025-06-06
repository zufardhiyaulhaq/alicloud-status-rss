package notification

import (
	"context"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/model"
)

type NotificationInterface interface {
	SendNotification(ctx context.Context, message model.Message) error
}
