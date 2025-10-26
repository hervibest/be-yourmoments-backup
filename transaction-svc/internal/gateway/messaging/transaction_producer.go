package producer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/event"
)

type TransactionProducer interface {
	ScheduleTransactionTaskExpiration(ctx context.Context, transactionID string) error
	ProduceCreateReviewEvent(ctx context.Context, creatorReviewCountEvent *event.CreatorReviewCountEvent) error
	ProduceTransactionSettledEvent(ctx context.Context, transactionSettledEvent *event.OwnerOwnPhotosEvent) error
	ProduceTransactionCanceledEvent(ctx context.Context, transactionCanceledEvent *event.CancelPhotosEvent) error
}

type transactionProducer struct {
	transactionTTL   time.Duration
	cacheAdapter     adapter.CacheAdapter
	messagingAdapter adapter.MessagingAdapter
	logs             *logger.Log
}

func NewTransactionProducer(cacheAdapter adapter.CacheAdapter, messagingAdapter adapter.MessagingAdapter, logs *logger.Log) TransactionProducer {
	ttlStr := utils.GetEnv("TRANSACTION_EXPIRATION_TTL") // misal "60"
	ttlInt, err := strconv.Atoi(ttlStr)
	if err != nil || ttlInt <= 0 {
		log.Printf("Invalid TRANSACTION_EXPIRATION_TTL value: %q, defaulting to 60 seconds", ttlStr)
		ttlInt = 60
	}

	return &transactionProducer{
		transactionTTL:   time.Duration(ttlInt) * time.Second,
		cacheAdapter:     cacheAdapter,
		messagingAdapter: messagingAdapter,
		logs:             logs,
	}
}

func (s *transactionProducer) ScheduleTransactionTaskExpiration(ctx context.Context, transactionID string) error {
	key := fmt.Sprintf("task:%s:expire", transactionID)
	err := s.cacheAdapter.SetEx(ctx, key, transactionID, s.transactionTTL)
	if err != nil {
		return fmt.Errorf("failed to schedule expiration for transaction %s: %w", transactionID, err)
	}

	log.Printf("Transaction ID %s scheduled to expire in %.2f seconds", transactionID, s.transactionTTL.Seconds())
	return nil
}

func (s *transactionProducer) ProduceCreateReviewEvent(ctx context.Context, creatorReviewCountEvent *event.CreatorReviewCountEvent) error {
	subject := "creator.review.updated"

	err := s.messagingAdapter.Publish(ctx, subject, creatorReviewCountEvent)
	if err != nil {
		return fmt.Errorf("failed to publish review update event: %w", err)
	}

	log.Printf("Published review update event for creator %s", creatorReviewCountEvent.Id)
	return nil
}

func (s *transactionProducer) ProduceTransactionSettledEvent(ctx context.Context, transactionSettledEvent *event.OwnerOwnPhotosEvent) error {
	subject := "transaction.settled"

	err := s.messagingAdapter.Publish(ctx, subject, transactionSettledEvent)
	if err != nil {
		return fmt.Errorf("failed to publish transaction settled event: %w", err)
	}

	// log.Printf("Published transaction settled event for transaction %s", transactionSettledEvent.TransactionId)
	return nil
}

func (s *transactionProducer) ProduceTransactionCanceledEvent(ctx context.Context, transactionCanceledEvent *event.CancelPhotosEvent) error {
	subject := "transaction.canceled"

	err := s.messagingAdapter.Publish(ctx, subject, transactionCanceledEvent)
	if err != nil {
		return fmt.Errorf("failed to publish transaction canceled event: %w", err)
	}

	// log.Printf("Published transaction canceled event for transaction %s", transactionCanceledEvent.TransactionId)
	return nil
}
