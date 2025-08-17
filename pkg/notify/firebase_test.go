package notify_test

import (
	"context"
	"testing"
	"zhasa2.0/pkg/notify"

	"github.com/stretchr/testify/require"
)

func TestSendPushNotification_Integration(t *testing.T) {
	credFile := "../../config/firebase.json"
	token := ""

	client, err := notify.NewFirebaseClient(credFile)
	require.NoError(t, err)

	msg := notify.Message{
		Heading: "Test notification",
		Message: "Hello from integration test",
		Payload: map[string]string{"deeplink": "doschamp://news?id=1"},
	}

	topic := "news-topic"

	err = client.SubscribeToTopic(context.Background(), token, topic)
	require.NoError(t, err)

	err = client.SendPushNotification(context.Background(), msg, topic)
	require.NoError(t, err)
}
