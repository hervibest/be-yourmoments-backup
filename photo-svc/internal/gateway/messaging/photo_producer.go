package producer

import (
	"context"
	"fmt"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
)

type PhotoProducer interface {
	ProduceBulkPhoto(ctx context.Context, creatorEvent *event.BulkPhotoEvent) error
	ProduceSinglePhoto(ctx context.Context, creatorEvent *event.SinglePhotoEvent) error
	ProduceSingleFacecam(ctx context.Context, creatorEvent *event.SingleFacecamEvent) error
	ProducePersistFacecam(ctx context.Context, persistEvent *event.PersistFacecamEvent) error
}

type photoProducer struct {
	messagingAdapter adapter.MessagingAdapter
	logs             *logger.Log
}

func NewPhotoProducer(messagingAdapter adapter.MessagingAdapter, logs *logger.Log) PhotoProducer {
	return &creatorProducer{
		messagingAdapter: messagingAdapter,
		logs:             logs,
	}
}

func (s *creatorProducer) ProduceBulkPhoto(ctx context.Context, creatorEvent *event.BulkPhotoEvent) error {
	subject := "photo.bulk"

	err := s.messagingAdapter.Publish(ctx, subject, creatorEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create creator event: %w", err)
	}

	s.logs.Log(fmt.Sprintf("Published bulk photo for event id %s", creatorEvent.EventID))
	return nil
}

func (s *creatorProducer) ProduceSinglePhoto(ctx context.Context, creatorEvent *event.SinglePhotoEvent) error {
	subject := "photo.single.photo"

	err := s.messagingAdapter.Publish(ctx, subject, creatorEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create creator event: %w", err)
	}

	s.logs.Log(fmt.Sprintf("Published single photo for event id %s", creatorEvent.EventID))
	return nil
}

func (s *creatorProducer) ProduceSingleFacecam(ctx context.Context, creatorEvent *event.SingleFacecamEvent) error {
	subject := "photo.single.facecam"

	err := s.messagingAdapter.Publish(ctx, subject, creatorEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create creator event: %w", err)
	}

	s.logs.Log(fmt.Sprintf("Published single facecam for event id %s", creatorEvent.EventID))
	return nil
}

func (s *creatorProducer) ProducePersistFacecam(ctx context.Context, creatorEvent *event.PersistFacecamEvent) error {
	subject := "photo.persist.facecam"

	err := s.messagingAdapter.Publish(ctx, subject, creatorEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create creator event: %w", err)
	}

	s.logs.Log(fmt.Sprintf("Published persist facecam for user id %s", creatorEvent.UserID))
	return nil
}
