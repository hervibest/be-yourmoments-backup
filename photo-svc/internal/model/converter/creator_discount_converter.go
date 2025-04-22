package converter

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/model"
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
