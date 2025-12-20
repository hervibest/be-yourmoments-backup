package subscriber

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

type UserSubscriber struct {
	js           nats.JetStreamContext
	useCase      usecase.UserDeviceUseCase
	subject      string
	consumerName string
	durableName  string
	logs         logger.Log
}

func NewUserSubscriber(js nats.JetStreamContext, useCase usecase.UserDeviceUseCase, logs logger.Log) *UserSubscriber {
	return &UserSubscriber{
		js:           js,
		useCase:      useCase,
		subject:      "user.device.created",
		consumerName: "notification_svc_consumer",
		durableName:  "notification_svc_durable",
		logs:         logs,
	}
}

func (s *UserSubscriber) Start(ctx context.Context) error {
	sub, err := s.js.PullSubscribe(
		s.subject,
		s.durableName,
		nats.BindStream("USER_DEVICE_STREAM"),
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
					event := new(event.UserDeviceEvent)
					if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
						s.logs.CustomError("failed to unmarshal event: %v", err)
						_ = msg.Nak()
						continue
					}

					s.logs.CustomLog("Processing event: %+v", event)

					request := &model.CreateDeviceRequest{
						UserID:      event.UserID,
						DeviceToken: event.DeviceToken,
						Platform:    event.Platform,
					}

					if err := s.useCase.CreateDevice(ctx, request); err != nil {
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
