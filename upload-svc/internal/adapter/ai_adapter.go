package adapter

import (
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	discovery "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"

	aipb "github.com/hervibest/be-yourmoments-backup/pb/ai"

	"context"
	"log"
)

type AiAdapter interface {
	ProcessPhoto(ctx context.Context, fileId, fileUrl string) error
	ProcessFacecam(ctx context.Context, userId, fileUrl string) error
}

type aiAdapter struct {
	client aipb.AiServiceClient
	logs   *logger.Log
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (AiAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "ai-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
	logs.Log("successfuly connected to ai-svc-grpc")
	client := aipb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
		logs:   logs,
	}, nil
}

func (a *aiAdapter) ProcessPhoto(ctx context.Context, userId, fileUrl string) error {
	processPhotoRequest := &aipb.ProcessPhotoRequest{
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
	processPhotoRequest := &aipb.ProcessFacecamRequest{
		Id:  fileId,
		Url: fileUrl,
	}

	_, err := a.client.ProcessFacecam(ctx, processPhotoRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *aiAdapter) CreatePhotos(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo) error {
	pbAIPhotos := make([]*aipb.AIPhoto, len(*photos))
	for _, photo := range *photos {
		pbAIPhoto := &aipb.AIPhoto{
			Id:            photo.Id,
			CollectionUrl: photo.CollectionUrl,
		}
		pbAIPhotos = append(pbAIPhotos, pbAIPhoto)
	}

	pbAIBulkPhoto := &aipb.AIBulkPhoto{
		Id:        bulkPhoto.Id,
		CreatorId: bulkPhoto.CreatorId,
	}

	pbRequest := &aipb.ProcessBulkPhotoRequest{
		ProcessBulkAi: pbAIBulkPhoto,
		ProcessAi:     pbAIPhotos,
	}

	_, err := a.client.ProcessBulkPhoto(ctx, pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}
