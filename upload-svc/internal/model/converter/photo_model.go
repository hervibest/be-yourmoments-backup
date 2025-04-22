package converter

import (
	"be-yourmoments/upload-svc/internal/model"

	"github.com/be-yourmoments/pb"
)

func GrpcToCreateRequest(req *pb.UpdatePhotographerPhotoRequest) *model.RequestUpdateProcessedPhoto {

	userId := make([]string, len(req.UserId))
	for _, value := range req.UserId {
		userId = append(userId, value)
	}

	return &model.RequestUpdateProcessedPhoto{
		Id:     req.GetId(),
		UserId: userId,
	}

}
