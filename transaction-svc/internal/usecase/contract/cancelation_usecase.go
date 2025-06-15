package contract

import "context"

type CancelationUseCase interface {
	ExpirePendingTransaction(ctx context.Context, transactionId string) error
	CancelPendingTransaction(ctx context.Context, transactionId string) error
}
