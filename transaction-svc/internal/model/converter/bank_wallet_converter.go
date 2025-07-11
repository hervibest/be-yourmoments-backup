package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func BankWalletToResponse(bankWallet *entity.BankWallet) *model.BankWalletResponse {
	return &model.BankWalletResponse{
		Id:            bankWallet.Id,
		BankId:        bankWallet.BankId,
		WalletId:      bankWallet.WalletId,
		FullName:      bankWallet.FullName,
		AccountNumber: bankWallet.AccountNumber,
		CreatedAt:     bankWallet.CreatedAt,
		UpdatedAt:     bankWallet.UpdatedAt,
	}
}

func BankWalletsToResponses(bankWallets *[]*entity.BankWallet) *[]*model.BankWalletResponse {
	responses := make([]*model.BankWalletResponse, 0)
	for _, bankWallet := range *bankWallets {
		bankWalletResponse := BankWalletToResponse(bankWallet)
		responses = append(responses, bankWalletResponse)
	}
	return &responses
}
