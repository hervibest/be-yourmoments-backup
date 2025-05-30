package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func ReviewToResponse(review *entity.CreatorReview) *model.CreatorReviewResponse {
	var (
		commentPtr *string
	)

	if review.Comment.Valid == true {
		commentPtr = &review.Comment.String
	}

	return &model.CreatorReviewResponse{
		Id:                  review.Id,
		TransactionDetailId: review.TransactionDetailId,
		CreatorId:           review.CreatorId,
		UserId:              review.UserId,
		Rating:              review.Rating,
		Comment:             commentPtr,
		CreatedAt:           review.CreatedAt,
		UpdatedAt:           review.UpdatedAt,
	}
}

func ReviewsToResponses(reviews *[]*entity.CreatorReview) *[]*model.CreatorReviewResponse {
	responses := make([]*model.CreatorReviewResponse, 0)
	for _, review := range *reviews {
		var (
			commentPtr *string
		)

		if review.Comment.Valid == true {
			commentPtr = &review.Comment.String
		}

		response := &model.CreatorReviewResponse{
			Id:                  review.Id,
			TransactionDetailId: review.TransactionDetailId,
			CreatorId:           review.CreatorId,
			Rating:              review.Rating,
			UserId:              review.UserId,
			Comment:             commentPtr,
			CreatedAt:           review.CreatedAt,
			UpdatedAt:           review.UpdatedAt,
		}
		responses = append(responses, response)
	}
	return &responses
}
