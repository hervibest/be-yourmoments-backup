package converter

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/helper/nullable"
	"be-yourmoments/photo-svc/internal/model"
)

func BulkPhotoDetailToResponse(items *[]*entity.BulkPhotoDetail) *model.GetBulkPhotoDetailResponse {
	photoResponses := make([]*model.PhotoResponse, 0)
	for _, item := range *items {
		photoResponse := &model.PhotoResponse{
			Id:             item.PhotoId,
			CreatorId:      item.PhotoCreatorId,
			Title:          item.PhotoTitle,
			OwnedByUserId:  nullable.SQLStringToPtr(item.PhotoOwnedByUserId),
			CompressedUrl:  nullable.SQLStringToPtr(item.PhotoOwnedByUserId),
			IsThisYouURL:   nullable.SQLStringToPtr(item.PhotoOwnedByUserId),
			YourMomentsUrl: nullable.SQLStringToPtr(item.PhotoOwnedByUserId),
			CollectionUrl:  nullable.SQLStringToPtr(item.PhotoOwnedByUserId),
			Price:          item.PhotoPrice,
			PriceStr:       item.PhotoPriceStr,
			Latitude:       nullable.SQLFloat64ToPtr(item.PhotoLatitude),
			Longitude:      nullable.SQLFloat64ToPtr(item.PhotoLongitude),
			Description:    nullable.SQLStringToPtr(item.PhotoDescription),
			OriginalAt:     item.PhotoOriginalAt,
			CreatedAt:      item.PhotoCreatedAt,
			UpdatedAt:      item.PhotoUpdatedAt,
		}
		photoResponses = append(photoResponses, photoResponse)
	}

	bulkPhotoDetailResponse := &model.GetBulkPhotoDetailResponse{
		Id:        (*items)[0].BulkPhotoId,
		CreatorId: (*items)[0].BulkPhotoId,
		Status:    (*items)[0].BulkPhotoStatus,
		Photo:     &photoResponses,
		CreatedAt: (*items)[0].BulkPhotoCreatedAt,
		UpdatedAt: (*items)[0].BulkPhotoUpdatedAt,
	}

	return bulkPhotoDetailResponse
}
