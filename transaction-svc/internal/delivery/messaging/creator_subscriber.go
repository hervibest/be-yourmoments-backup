package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

type CreatorSubscriber struct {
	js            nats.JetStreamContext
	walletUseCase usecase.WalletUseCase
	subject       string
	consumerName  string
	durableName   string
	logs          *logger.Log
}

func NewCreatorSubscriber(js nats.JetStreamContext, walletUseCase usecase.WalletUseCase, logs *logger.Log) *CreatorSubscriber {
	return &CreatorSubscriber{
		js:            js,
		walletUseCase: walletUseCase,
		subject:       "creator.created",
		consumerName:  "transaction_svc_consumer",
		durableName:   "transaction_svc_durable",
		logs:          logs,
	}
}

func (s *CreatorSubscriber) Start(ctx context.Context) error {
	sub, err := s.js.PullSubscribe(
		s.subject,
		s.durableName,
		nats.BindStream("CREATOR_STREAM"),
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
					event := new(event.CreatorEvent)
					if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
						s.logs.CustomError("failed to unmarshal event: %v", err)
						_ = msg.Nak()
						continue
					}

					s.logs.CustomLog("Processing event: %+v", event)

					request := &model.CreateWalletRequest{
						CreatorId: event.Id,
					}

					if _, err := s.walletUseCase.CreateWallet(ctx, request); err != nil {
						s.logs.CustomError("failed to create wallet: %v", err)
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
