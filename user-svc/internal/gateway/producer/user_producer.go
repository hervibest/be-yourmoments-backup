package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model/event"
)

type UserProducer interface {
	ProduceUserCreated(ctx context.Context, userEvent *event.UserEvent) error
	ProduceUserDeviceCreated(ctx context.Context, userDeviceEvent *event.UserDeviceEvent) error
}

type userProducer struct {
	messagingAdapter adapter.MessagingAdapter
	logs             logger.Log
}

func NewUserProducer(messagingAdapter adapter.MessagingAdapter, logs logger.Log) UserProducer {
	return &userProducer{
		messagingAdapter: messagingAdapter,
		logs:             logs,
	}
}

func (s *userProducer) ProduceUserCreated(ctx context.Context, userEvent *event.UserEvent) error {
	subject := "user.created"

	err := s.messagingAdapter.Publish(ctx, subject, userEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create user event: %w", err)
	}

	log.Printf("Published create user event for user id %s", userEvent.Id)
	return nil
}

func (s *userProducer) ProduceUserDeviceCreated(ctx context.Context, userDeviceEvent *event.UserDeviceEvent) error {
	subject := "user.device.created"

	err := s.messagingAdapter.Publish(ctx, subject, userDeviceEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create user device event: %w", err)
	}

	log.Printf("Published create user device event for user id %s", userDeviceEvent.UserID)
	return nil
}
