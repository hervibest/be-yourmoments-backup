package adapter

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type CloudMessagingAdapter interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
	SendEachForMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
}

func NewCloudMessagingAdapter(app *firebase.App) CloudMessagingAdapter {
	ctx := context.Background()
	cloudMessagingAdapter, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase messaging: %v", err)
	}

	return cloudMessagingAdapter
}
