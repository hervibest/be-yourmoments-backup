package converter

import (
	"be-yourmoments/photo-svc/internal/model"

	photopb "github.com/be-yourmoments/pb/photo"
)

func GrpcToCreateRequest(req *photopb.UpdatePhotographerPhotoRequest) *model.RequestUpdateProcessedPhoto {

	userId := make([]string, len(req.UserId))
	for _, value := range req.UserId {
		userId = append(userId, value)
	}

	return &model.RequestUpdateProcessedPhoto{
		Id:     req.GetId(),
		UserId: userId,
	}

}
