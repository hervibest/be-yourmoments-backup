package subscriber

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

type CreatorReviewSubscriber struct {
	js           nats.JetStreamContext
	useCase      usecase.CreatorUseCase
	subject      string
	consumerName string
	durableName  string
	logs         *logger.Log
}

func NewCreatorReviewSubscriber(js nats.JetStreamContext, useCase usecase.CreatorUseCase, logs *logger.Log) *CreatorReviewSubscriber {
	return &CreatorReviewSubscriber{
		js:           js,
		useCase:      useCase,
		subject:      "creator.review.updated",
		consumerName: "photo_svc_consumer",
		durableName:  "photo_svc_durable",
		logs:         logs,
	}
}

func (s *CreatorReviewSubscriber) Start(ctx context.Context) error {
	sub, err := s.js.PullSubscribe(
		s.subject,
		s.durableName,
		nats.BindStream("CREATOR_REVIEW_STREAM"), // ganti dengan stream name
	)
	if err != nil {
		return fmt.Errorf("failed to create pull subscription: %w", err)
	}

	s.logs.CustomLog("Started synchronous subscriber for", s.subject)

	func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping subscriber...")
				return
			default:
				msgs, err := sub.Fetch(10, nats.MaxWait(2*time.Second))
				if err != nil {
					if err == nats.ErrTimeout {
						time.Sleep(200 * time.Millisecond) // 🔥 WAJIB
						continue
					}
					s.logs.Error(fmt.Sprintf("failed to fetch messages with error %v", err))
					time.Sleep(time.Second)
					continue
				}
				for _, msg := range msgs {
					event := new(event.CreatorReviewCountEvent)
					if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
						s.logs.CustomError("failed to unmarshal event: %v", err)
						_ = msg.Nak()
						continue
					}

					s.logs.CustomLog("Processing event: %+v", event)

					request := &model.UpdateCreatorTotalRatingRequest{
						Id:          event.Id,
						Rating:      event.Rating,
						RatingCount: event.RatingCount,
					}

					if _, err := s.useCase.UpdateCreatorTotalReview(ctx, request); err != nil {
						s.logs.CustomError("failed to update creator review: %v", err)
						_ = msg.Nak()
						continue
					}

					if err := msg.Ack(); err != nil {
						s.logs.CustomError("failed to ack message: %v", err)
					}
				}
			}
		}
	}()

	return nil
}
