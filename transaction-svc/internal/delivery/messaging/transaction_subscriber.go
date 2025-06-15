package messaging

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase/contract"
)

type TransactionSubscriber interface {
	SubscribeTransactionExpire(ctx context.Context)
}

type transactionSubscriber struct {
	cacheAdapter       adapter.CacheAdapter
	cancelationUseCase contract.CancelationUseCase
	logs               *logger.Log
}

func NewTransactionSubsciber(cacheAdapter adapter.CacheAdapter,
	cancelationUseCase contract.CancelationUseCase,
	logs *logger.Log) TransactionSubscriber {
	return &transactionSubscriber{
		cacheAdapter:       cacheAdapter,
		cancelationUseCase: cancelationUseCase,
		logs:               logs,
	}
}

func (s *transactionSubscriber) SubscribeTransactionExpire(ctx context.Context) {
	s.logs.Log("Successfully subscribed to redis for transaction expire")
	pubsub := s.cacheAdapter.PSubscribe(ctx, "__keyevent@0__:expired") // Gunakan DB 0 secara eksplisit
	for {
		select {
		case <-ctx.Done():
			s.logs.Log("Shutting down Redis subscriber gracefully...")
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				// Jika ctx sudah dibatalkan, keluar juga
				if ctx.Err() != nil {
					s.logs.CustomLog("Context canceled, exiting subscriber:", ctx.Err())
					return
				}
				s.logs.CustomError("Error when getting redis message:", err)
				continue
			}

			log.Println("Received redis message:", msg.Channel, msg.Payload)

			payload := msg.Payload

			if strings.HasPrefix(payload, "task:") && strings.HasSuffix(payload, ":expire") {
				orderID := strings.TrimSuffix(strings.TrimPrefix(payload, "task:"), ":expire")

				s.logs.CustomLog("Begin expiration process of task ID %s", orderID)
				if err := s.cancelationUseCase.ExpirePendingTransaction(ctx, orderID); err != nil {
					s.logs.CustomError("Failed to expire task:", err)
				} else {
					s.logs.Log(fmt.Sprintf("Task ID %s sucessfully expire", orderID))
				}
			} else {
				s.logs.CustomLog("Invalid payload format:", payload)
			}
		}
	}
}
