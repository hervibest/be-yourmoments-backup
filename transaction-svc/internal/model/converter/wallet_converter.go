package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
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
