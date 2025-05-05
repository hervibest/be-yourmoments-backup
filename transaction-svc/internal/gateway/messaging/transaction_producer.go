package producer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
)

type TransactionProducer interface {
	ScheduleTransactionTaskExpiration(ctx context.Context, transactionID string) error
}

type transactionProducer struct {
	transactionTTL time.Duration
	cacheAdapter   adapter.CacheAdapter
}

func NewTransactionProducer(cacheAdapter adapter.CacheAdapter) TransactionProducer {
	ttlStr := utils.GetEnv("TRANSACTION_EXPIRATION_TTL") // misal "60"
	ttlInt, err := strconv.Atoi(ttlStr)
	if err != nil || ttlInt <= 0 {
		log.Printf("Invalid TRANSACTION_EXPIRATION_TTL value: %q, defaulting to 60 seconds", ttlStr)
		ttlInt = 60
	}

	return &transactionProducer{
		transactionTTL: time.Duration(ttlInt) * time.Second,
		cacheAdapter:   cacheAdapter,
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
