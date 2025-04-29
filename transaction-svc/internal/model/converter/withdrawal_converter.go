package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func WithdrawalToResponse(withdrawal *entity.Withdrawal) *model.WithdrawalResponse {
	return &model.WithdrawalResponse{
		Id:           withdrawal.Id,
		WalletId:     withdrawal.WalletId,
		BankWalletId: withdrawal.BankWalletId,
		Amount:       withdrawal.Amount,
		Status:       withdrawal.Status,
		Description:  withdrawal.Description,
		CreatedAt:    withdrawal.CreatedAt,
		UpdatedAt:    withdrawal.UpdatedAt,
	}
}

func WithdrawalsToResponses(withdrawals *[]*entity.Withdrawal) *[]*model.WithdrawalResponse {
	responses := make([]*model.WithdrawalResponse, 0)
	for _, withdrawal := range *withdrawals {
		withdrawalReponse := &model.WithdrawalResponse{
			Id:           withdrawal.Id,
			WalletId:     withdrawal.WalletId,
			BankWalletId: withdrawal.BankWalletId,
			Amount:       withdrawal.Amount,
			Status:       withdrawal.Status,
			Description:  withdrawal.Description,
			CreatedAt:    withdrawal.CreatedAt,
			UpdatedAt:    withdrawal.UpdatedAt,
		}

		responses = append(responses, withdrawalReponse)
	}
	return &responses
}
