package converter

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/model"
)

func TransactionToResponse(transaction *entity.Transaction, redirectUrl string) *model.CreateTransactionResponse {
	return &model.CreateTransactionResponse{
		TransactionId: transaction.Id,
		SnapToken:     transaction.SnapToken.String,
		RedirectURL:   redirectUrl,
	}
}
