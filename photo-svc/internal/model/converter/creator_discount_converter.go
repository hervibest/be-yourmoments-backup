package converter

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
)

func CreatorDiscountToResponse(discount *entity.CreatorDiscount) *model.CreatorDiscountResponse {
	return &model.CreatorDiscountResponse{
		Id:           discount.Id,
		CreatorId:    discount.CreatorId,
		Name:         discount.Name,
		MinQuantity:  discount.MinQuantity,
		DiscountType: discount.DiscountType,
		Value:        discount.Value,
		Active:       discount.Active,
		CreatedAt:    discount.CreatedAt,
		UpdatedAt:    discount.UpdatedAt,
	}
}

func CreatorDiscountsToResponses(discounts []*entity.CreatorDiscount) *[]*model.CreatorDiscountResponse {
	discountResponses := make([]*model.CreatorDiscountResponse, 0, len(discounts))
	for _, discount := range discounts {
		discountResponse := &model.CreatorDiscountResponse{
			Id:           discount.Id,
			CreatorId:    discount.CreatorId,
			Name:         discount.Name,
			MinQuantity:  discount.MinQuantity,
			DiscountType: discount.DiscountType,
			Value:        discount.Value,
			Active:       discount.Active,
			CreatedAt:    discount.CreatedAt,
			UpdatedAt:    discount.UpdatedAt,
		}
		discountResponses = append(discountResponses, discountResponse)
	}
	return &discountResponses
}
