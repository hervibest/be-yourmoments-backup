package converter

import (
	"be-yourmoments/upload-svc/internal/model"

	photopb "github.com/be-yourmoments/pb/photo"
)

func GrpcToCreateRequest(req *photopb.UpdatePhotographerPhotoRequest) *model.RequestUpdateProcessedPhoto {
	userId := make([]string, len(req.UserId))
	userId = append(userId, req.UserId...)

	return &model.RequestUpdateProcessedPhoto{
		Id:     req.GetId(),
		UserId: userId,
	}

}
