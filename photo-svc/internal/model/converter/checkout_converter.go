package converter

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
)

func CheckoutItemToResponse(checkoutItems *[]*model.CheckoutItem, totalPrice, totalDiscount int32, createdAt *time.Time) *model.PreviewCheckoutResponse {
	chekoutItemsWeb := make([]*model.CheckoutItemWeb, 0, len(*checkoutItems))
	for _, checkoutItem := range *checkoutItems {
		var discount *model.DiscountItem
		if checkoutItem.Discount != 0 && checkoutItem.DiscountId != "" {
			discount = &model.DiscountItem{
				Discount:            checkoutItem.Discount,
				DiscountMinQuantity: checkoutItem.DiscountMinQuantity,
				DiscountValue:       checkoutItem.DiscountValue,
				DiscountId:          checkoutItem.DiscountId,
				DiscountType:        checkoutItem.DiscountType,
			}
		}

		checkoutItemWeb := &model.CheckoutItemWeb{
			PhotoId:        checkoutItem.PhotoId,
			CreatorId:      checkoutItem.CreatorId,
			Title:          checkoutItem.Title,
			YourMomentsUrl: checkoutItem.YourMomentsUrl,
			Price:          checkoutItem.Price,
			Discount:       discount,
			FinalPrice:     checkoutItem.FinalPrice,
		}
		chekoutItemsWeb = append(chekoutItemsWeb, checkoutItemWeb)
	}
	return &model.PreviewCheckoutResponse{
		Items:         &chekoutItemsWeb,
		TotalPrice:    totalPrice,
		TotalDiscount: totalDiscount,
		CreatedAt:     createdAt,
	}
}
