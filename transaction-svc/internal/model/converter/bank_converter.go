package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func BankToResponse(bank *entity.Bank) *model.BankResponse {
	return &model.BankResponse{
		Id:        bank.Id,
		BankCode:  bank.BankCode,
		Name:      bank.Name,
		Alias:     nullable.SQLStringToPtr(bank.Alias),
		SwiftCode: nullable.SQLStringToPtr(bank.SwiftCode),
		LogoUrl:   nullable.SQLStringToPtr(bank.LogoUrl),
		CreatedAt: bank.CreatedAt,
		UpdatedAt: bank.UpdatedAt,
	}
}

func BanksToResponses(banks *[]*entity.Bank) *[]*model.BankResponse {
	bankReponses := make([]*model.BankResponse, 0)
	for _, bank := range *banks {
		bankResponse := BankToResponse(bank)
		bankReponses = append(bankReponses, bankResponse)
	}
	return &bankReponses
}
