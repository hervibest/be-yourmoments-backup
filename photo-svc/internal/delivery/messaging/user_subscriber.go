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

type UserSubscriber struct {
	js           nats.JetStreamContext
	useCase      usecase.CreatorUseCase
	subject      string
	consumerName string
	durableName  string
	logs         *logger.Log
}

func NewUserSubscriber(js nats.JetStreamContext, useCase usecase.CreatorUseCase, logs *logger.Log) *UserSubscriber {
	return &UserSubscriber{
		js:           js,
		useCase:      useCase,
		subject:      "user.created",
		consumerName: "photo_svc_consumer",
		durableName:  "photo_svc_durable",
		logs:         logs,
	}
}

func (s *UserSubscriber) Start(ctx context.Context) error {
	sub, err := s.js.PullSubscribe(
		s.subject,
		s.durableName,
		nats.BindStream("USER_STREAM"),
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
				if err != nil && err != nats.ErrTimeout {
					s.logs.CustomLog("Fetch error: %v", err)
					continue
				}

				for _, msg := range msgs {
					event := new(event.UserEvent)
					if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
						s.logs.CustomError("failed to unmarshal event: %v", err)
						_ = msg.Nak()
						continue
					}

					s.logs.CustomLog("Processing event: %+v", event)

					request := &model.CreateCreatorRequest{
						UserId: event.Id,
					}

					if _, err := s.useCase.CreateCreator(ctx, request); err != nil {
						s.logs.CustomError("failed to create creator: %v", err)
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
