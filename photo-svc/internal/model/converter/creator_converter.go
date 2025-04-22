package converter

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/model"
)

func CreatorToResponse(creator *entity.Creator) *model.CreatorResponse {
	return &model.CreatorResponse{
		Id:          creator.Id,
		UserId:      creator.UserId,
		Rating:      creator.Rating,
		RatingCount: creator.RatingCount,
		VerifiedAt:  creator.VerifiedAt,
		CreatedAt:   creator.CreatedAt,
		UpdatedAt:   creator.UpdatedAt,
	}
}
