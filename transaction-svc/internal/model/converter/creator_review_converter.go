package converter

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

func ReviewToResponse(review *entity.CreatorReview) *model.CreatorReviewResponse {
	return &model.CreatorReviewResponse{
		Id:                  review.Id,
		TransactionDetailId: review.TransactionDetailId,
		CreatorId:           review.CreatorId,
		UserId:              review.UserId,
		Rating:              review.Rating,
		Comment:             nullable.SQLStringToPtr(review.Comment),
		CreatedAt:           review.CreatedAt,
		UpdatedAt:           review.UpdatedAt,
	}
}

func ReviewsToResponses(reviews *[]*entity.CreatorReview) *[]*model.CreatorReviewResponse {
	responses := make([]*model.CreatorReviewResponse, 0)
	for _, review := range *reviews {
		response := ReviewToResponse(review)
		responses = append(responses, response)
	}
	return &responses
}
