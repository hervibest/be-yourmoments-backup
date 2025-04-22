package converter

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/model"
)

func BankToResponse(bank *entity.Bank) *model.BankResponse {
	var (
		aliasPtr     *string
		swiftCodePtr *string
		LogoUrlPtr   *string
	)

	if bank.Alias.Valid {
		aliasPtr = &bank.Alias.String
	}
	if bank.SwiftCode.Valid {
		swiftCodePtr = &bank.SwiftCode.String
	}
	if bank.LogoUrl.Valid {
		LogoUrlPtr = &bank.LogoUrl.String
	}

	return &model.BankResponse{
		Id:        bank.Id,
		BankCode:  bank.BankCode,
		Name:      bank.Name,
		Alias:     aliasPtr,
		SwiftCode: swiftCodePtr,
		LogoUrl:   LogoUrlPtr,
		CreatedAt: bank.CreatedAt,
		UpdatedAt: bank.UpdatedAt,
	}
}

func BanksToResponses(banks *[]*entity.Bank) *[]*model.BankResponse {
	bankReponses := make([]*model.BankResponse, 0)
	for _, bank := range *banks {
		var (
			aliasPtr     *string
			swiftCodePtr *string
			LogoUrlPtr   *string
		)

		if bank.Alias.Valid {
			aliasPtr = &bank.Alias.String
		}
		if bank.SwiftCode.Valid {
			swiftCodePtr = &bank.SwiftCode.String
		}
		if bank.LogoUrl.Valid {
			LogoUrlPtr = &bank.LogoUrl.String
		}

		bankResponse := &model.BankResponse{
			Id:        bank.Id,
			BankCode:  bank.BankCode,
			Name:      bank.Name,
			Alias:     aliasPtr,
			SwiftCode: swiftCodePtr,
			LogoUrl:   LogoUrlPtr,
			CreatedAt: bank.CreatedAt,
			UpdatedAt: bank.UpdatedAt,
		}

		bankReponses = append(bankReponses, bankResponse)
	}
	return &bankReponses
}
