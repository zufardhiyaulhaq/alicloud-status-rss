package lark

import (
	"context"
	"log"

	golark "github.com/go-lark/lark"
	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/model"
)

type LarkNotification struct {
	WebhookURLs []string
}

func (ln *LarkNotification) SendNotification(ctx context.Context, message model.Message) error {
	for _, url := range ln.WebhookURLs {
		bot := golark.NewNotificationBot(url)

		contentStr := "[🚨 " + message.Type + " 🚨] " + message.Title
		content := golark.NewPostBuilder().
			Title(contentStr).
			TextTag(message.Content, 3, true).
			LinkTag("read more", message.Link).
			Render()
		buffer := golark.NewMsgBuffer(golark.MsgPost).Post(content).Build()

		_, err := bot.PostNotificationV2(buffer)
		if err != nil {
			log.Printf("failed to send lark notification: %v\n", err.Error())
			return err
		}
	}

	return nil
}

func NewLarkClient(webhookURLs []string) *LarkNotification {
	return &LarkNotification{
		WebhookURLs: webhookURLs,
	}
}
