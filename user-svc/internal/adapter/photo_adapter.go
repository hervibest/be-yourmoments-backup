package adapter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/helper"
	discovery "be-yourmoments/user-svc/internal/helper/discovery"
	"context"

	photopb "github.com/be-yourmoments/pb/photo"
)

type PhotoAdapter interface {
	CreateCreator(ctx context.Context, userId string) (*entity.Creator, error)
	GetCreator(ctx context.Context, creatorId string) (*entity.Creator, error)
}

type photoAdapter struct {
	client photopb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "photo-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
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
