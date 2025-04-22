package adapter

import (
	"be-yourmoments/upload-svc/internal/helper"
	discovery "be-yourmoments/upload-svc/internal/helper/discovery"
	"be-yourmoments/upload-svc/internal/helper/logger"

	"github.com/be-yourmoments/pb"

	"context"
	"log"
)

type AiAdapter interface {
	ProcessPhoto(ctx context.Context, fileId, fileUrl string) error
	ProcessFacecam(ctx context.Context, userId, fileUrl string) error
}

type aiAdapter struct {
	client pb.AiServiceClient
	logs   *logger.Log
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (AiAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "ai-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
	logs.Log("successfuly connected to ai-svc-grpc")
	client := pb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
		logs:   logs,
	}, nil
}

func (a *aiAdapter) ProcessPhoto(ctx context.Context, userId, fileUrl string) error {
	processPhotoRequest := &pb.ProcessPhotoRequest{
		Id:  userId,
		Url: fileUrl,
	}

	_, err := a.client.ProcessPhoto(ctx, processPhotoRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *aiAdapter) ProcessFacecam(ctx context.Context, fileId, fileUrl string) error {
	log.Println("REQUESTED PROCESS FACECAM VIA GRPC TO AI SERVER")
	processPhotoRequest := &pb.ProcessFacecamRequest{
		Id:  fileId,
		Url: fileUrl,
	}

	_, err := a.client.ProcessFacecam(ctx, processPhotoRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}
