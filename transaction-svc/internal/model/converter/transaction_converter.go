package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func TransactionToResponse(transaction *entity.Transaction, redirectUrl string) *model.CreateTransactionResponse {
	return &model.CreateTransactionResponse{
		TransactionId: transaction.Id,
		SnapToken:     transaction.SnapToken.String,
		RedirectURL:   redirectUrl,
	}
}
