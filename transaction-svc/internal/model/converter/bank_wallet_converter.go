package converter

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/model"
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
		bankWalletResponse := &model.BankWalletResponse{
			Id:            bankWallet.Id,
			BankId:        bankWallet.BankId,
			WalletId:      bankWallet.WalletId,
			FullName:      bankWallet.FullName,
			AccountNumber: bankWallet.AccountNumber,
			CreatedAt:     bankWallet.CreatedAt,
			UpdatedAt:     bankWallet.UpdatedAt,
		}

		responses = append(responses, bankWalletResponse)
	}
	return &responses
}
