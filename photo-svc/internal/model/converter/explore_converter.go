package converter

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
)

func ToDiscountIfValid(explore *entity.Explore) *model.CreatorDiscountResponse {
	if !explore.Name.Valid || explore.Name.String == "" {
		return nil
	}
	return &model.CreatorDiscountResponse{
		Name:         explore.Name.String,
		MinQuantity:  int(explore.MinQuantity.Int32),
		DiscountType: enum.DiscountType(explore.DiscountType.String),
		IsActive:     explore.IsActive.Bool,
		Value:        explore.Value.Int32,
	}
}

func ToCollectionOrIsYouURL(photoDetailType string, filekey string, generateCDN func(string) string) *model.PhotoUrlResponse {
	if photoDetailType == string(enum.YourMomentTypeCollection) {
		return &model.PhotoUrlResponse{
			CollectionUrl: generateCDN(filekey),
		}
	} else {
		return &model.PhotoUrlResponse{
			IsThisYouURL: generateCDN(filekey),
		}
	}
}

func ExploresToResponses(explores *[]*entity.Explore, generateCDN func(string) string) *[]*model.ExploreUserSimilarResponse {
	responses := make([]*model.ExploreUserSimilarResponse, 0)
	for _, explore := range *explores {
		photoUrlResponse := ToCollectionOrIsYouURL(explore.PhotoDetailType, explore.FileKey, generateCDN)

		photoStageResponse := &model.PhotoStageResponse{
			IsWishlist: explore.IsWishlist,
			IsResend:   explore.IsResend,
			IsCart:     explore.IsCart,
			IsFavorite: explore.IsFavorite,
		}

		discount := ToDiscountIfValid(explore)

		response := &model.ExploreUserSimilarResponse{
			PhotoId:    explore.PhotoId,
			UserId:     explore.UserId,
			Similarity: explore.Similarity,
			PhotoStage: photoStageResponse,
			CreatorId:  explore.CreatorId,
			Title:      explore.Title,
			PhotoUrl:   photoUrlResponse,
			Price:      explore.Price,
			PriceStr:   explore.PriceStr,
			Discount:   discount,
			OriginalAt: explore.OriginalAt,
			CreatedAt:  explore.CreatedAt,
			UpdatedAt:  explore.UpdatedAt,
		}

		responses = append(responses, response)
	}

	return &responses
}
