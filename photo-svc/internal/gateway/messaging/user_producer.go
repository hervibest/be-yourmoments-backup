package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
)

type CreatorProducer interface {
	ProduceCreatorCreated(ctx context.Context, creatorEvent *event.CreatorEvent) error
}

type creatorProducer struct {
	messagingAdapter adapter.MessagingAdapter
	logs             *logger.Log
}

func NewCreatorProducer(messagingAdapter adapter.MessagingAdapter, logs *logger.Log) CreatorProducer {
	return &creatorProducer{
		messagingAdapter: messagingAdapter,
		logs:             logs,
	}
}

func (s *creatorProducer) ProduceCreatorCreated(ctx context.Context, creatorEvent *event.CreatorEvent) error {
	subject := "creator.created"

	err := s.messagingAdapter.Publish(ctx, subject, creatorEvent)
	if err != nil {
		return fmt.Errorf("failed to publish create creator event: %w", err)
	}

	log.Printf("Published create creator event for creator id %s", creatorEvent.Id)
	return nil
}
