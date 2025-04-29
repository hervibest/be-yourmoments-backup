package converter

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
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
