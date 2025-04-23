package adapter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/helper"
	discovery "be-yourmoments/user-svc/internal/helper/discovery"
	"context"

	"github.com/be-yourmoments/pb"
)

type PhotoAdapter interface {
	CreateCreator(ctx context.Context, userId string) (*entity.Creator, error)
	GetCreator(ctx context.Context, creatorId string) (*entity.Creator, error)
}

type photoAdapter struct {
	client pb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "photo-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
	client := pb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
	}, nil
}

func (a *photoAdapter) CreateCreator(ctx context.Context, userId string) (*entity.Creator, error) {
	pbRequest := &pb.CreateCreatorRequest{
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
	pbRequest := &pb.GetCreatorRequest{
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
