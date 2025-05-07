package contract

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error)
	CheckAndUpdateTransaction(ctx context.Context, request *model.CheckAndUpdateTransactionRequest) error
	UserGetWithDetail(ctx context.Context, request *model.GetTransactionWithDetail) (*model.TransactionWithDetail, error)
	GetAllUserTransaction(ctx context.Context, request *model.GetAllUsertTransaction) (*[]*model.UserTransaction, *model.PageMetadata, error)
	CheckPaymentSignature(signatureKey, transcationId, statusCode, grossAmount string) (bool, string)
}
