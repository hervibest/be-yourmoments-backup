package adapter

import (
	"context"
	"log"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	discovery "github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/utils"
)

type PhotoAdapter interface {
	CreateCreator(ctx context.Context, userId string) (*entity.Creator, error)
	GetCreator(ctx context.Context, creatorId string) (*entity.Creator, error)
}

type photoAdapter struct {
	client photopb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry, logs logger.Log) (PhotoAdapter, error) {
	photoServiceName := utils.GetEnv("PHOTO_SVC_NAME")
	conn, err := discovery.NewGrpcClient(photoServiceName)
	if err != nil {
		logs.CustomError("failed to connect to the photo service due to an error : ", err)
		return nil, err
	}

	log.Printf("successfuly connected to %s", photoServiceName)
	client := photopb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
	}, nil
}

func (a *photoAdapter) CreateCreator(ctx context.Context, userId string) (*entity.Creator, error) {
	pbRequest := &photopb.CreateCreatorRequest{
		UserId: userId,
	}

	pbResponse, err := a.client.CreateCreator(context.Background(), pbRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	creator := &entity.Creator{
		Id: pbResponse.GetCreator().GetId(),
	}

	return creator, nil
}

func (a *photoAdapter) GetCreator(ctx context.Context, userId string) (*entity.Creator, error) {
	pbRequest := &photopb.GetCreatorRequest{
		UserId: userId,
	}

	pbResponse, err := a.client.GetCreator(context.Background(), pbRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	creator := &entity.Creator{
		Id: pbResponse.GetCreator().GetId(),
	}

	return creator, nil
}
