package transactionconsumer

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

// consumer.go
type TransactionConsumer struct {
	checkoutUC   usecase.CheckoutUseCase
	js           nats.JetStreamContext
	logs         *logger.Log
	subjects     []string
	durableNames map[string]string
}

func NewTransactionConsumer(
	checkoutUC usecase.CheckoutUseCase,
	js nats.JetStreamContext,
	logs *logger.Log,
) *TransactionConsumer {
	return &TransactionConsumer{
		checkoutUC: checkoutUC,
		js:         js,
		logs:       logs,
		subjects: []string{
			"transaction.settled",
			"transaction.canceled",
		},
		durableNames: map[string]string{
			"transaction.settled":  "transaction_settled_consumer",
			"transaction.canceled": "transaction_canceled_consumer",
		},
	}
}
