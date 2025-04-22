package converter

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/model"
)

func WalletToResponse(wallet *entity.Wallet) *model.WalletResponse {
	return &model.WalletResponse{
		Id:        wallet.Id,
		CreatorId: wallet.CreatorId,
		Balance:   int(wallet.Balance),
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
