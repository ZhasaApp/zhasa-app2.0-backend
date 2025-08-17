package notify

import "context"

type Client interface {
	SendPushNotification(ctx context.Context, msg Message, topic string) error
	SubscribeToTopic(ctx context.Context, token, topic string) error
}
