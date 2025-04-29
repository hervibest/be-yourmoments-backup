package converter

import (
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
)

func GrpcToCreateRequest(req *photopb.UpdatePhotographerPhotoRequest) *model.RequestUpdateProcessedPhoto {
	userId := make([]string, len(req.UserId))
	userId = append(userId, req.UserId...)

	return &model.RequestUpdateProcessedPhoto{
		Id:     req.GetId(),
		UserId: userId,
	}

}
