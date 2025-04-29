package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func WalletsToResponses(transactionWallets *[]*entity.TransactionWallet) *[]*model.TransactionWalletResponse {
	responses := make([]*model.TransactionWalletResponse, 0)
	for _, transactionWallet := range *transactionWallets {

		response := &model.TransactionWalletResponse{
			Id:                  transactionWallet.Id,
			WalletId:            transactionWallet.WalletId,
			TransactionDetailId: transactionWallet.TransactionDetailId,
			Amount:              transactionWallet.Amount,
			CreatedAt:           transactionWallet.CreatedAt,
			UpdatedAt:           transactionWallet.UpdatedAt,
		}
		responses = append(responses, response)
	}
	return &responses
}
