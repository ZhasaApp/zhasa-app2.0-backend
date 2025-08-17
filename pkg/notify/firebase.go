package notify

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type firebaseClient struct {
	client *messaging.Client
	Config string `json:"config"`
}

func NewFirebaseClient(config string) (Client, error) {
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(config))
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}
	return firebaseClient{
		client: client,
		Config: config,
	}, nil
}

func (f firebaseClient) SendPushNotification(ctx context.Context, msg Message, topic string) error {
	_, err := f.client.Send(ctx, &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: msg.Heading,
			Body:  msg.Message,
		},
		Data:    msg.Payload,
		Android: &messaging.AndroidConfig{Priority: "high"},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{"apns-priority": "10"},
		},
	})
	return err
}

func (f firebaseClient) SubscribeToTopic(ctx context.Context, token, topic string) error {
	_, err := f.client.SubscribeToTopic(ctx, []string{token}, topic)
	if err != nil {
		return err
	}
	return nil
}
