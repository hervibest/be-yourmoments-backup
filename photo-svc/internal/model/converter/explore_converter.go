package converter

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/model"
)

func ExploresToResponses(explores *[]*entity.Explore) *[]*model.ExploreUserSimilarResponse {
	responses := make([]*model.ExploreUserSimilarResponse, 0)
	for _, explore := range *explores {

		photoUrlResponse := &model.PhotoUrlResponse{
			IsThisYouURL:   explore.IsThisYouURL.String,
			YourMomentsUrl: explore.YourMomentsUrl.String,
		}

		photoStageResponse := &model.PhotoStageResponse{
			IsWishlist: explore.IsWishlist,
			IsResend:   explore.IsResend,
			IsCart:     explore.IsCart,
			IsFavorite: explore.IsFavorite,
		}

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
			OriginalAt: explore.OriginalAt,
			CreatedAt:  explore.CreatedAt,
			UpdatedAt:  explore.UpdatedAt,
		}

		responses = append(responses, response)
	}

	return &responses
}
