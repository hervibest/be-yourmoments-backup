package adapter

import (
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	discovery "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/utils"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"

	aipb "github.com/hervibest/be-yourmoments-backup/pb/ai"

	"context"
	"log"
)

type AiAdapter interface {
	ProcessPhoto(ctx context.Context, request *model.ProcessPhoto) error
	ProcessFacecam(ctx context.Context, request *model.ProcessFacecam) error
	ProcessBulkPhoto(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo) error
}

type aiAdapter struct {
	client aipb.AiServiceClient
	logs   logger.Log
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry, logs logger.Log) (AiAdapter, error) {
	aiServiceName := utils.GetEnv("AI_SVC_NAME")
	logs.Log("trying to connect to")
	conn, err := discovery.ServiceConnection(ctx, aiServiceName, registry, logs)
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

func (a *aiAdapter) ProcessPhoto(ctx context.Context, request *model.ProcessPhoto) error {
	processPhotoRequest := &aipb.ProcessPhotoRequest{
		Id:               request.PhotoId,
		CreatorId:        request.CreatorId,
		Url:              request.FileURL,
		OriginalFilename: request.OriginalFilename,
	}
	a.logs.Log("PROCESSED PHOTO")

	_, err := a.client.ProcessPhoto(ctx, processPhotoRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *aiAdapter) ProcessFacecam(ctx context.Context, request *model.ProcessFacecam) error {
	log.Println("REQUESTED PROCESS FACECAM VIA GRPC TO AI SERVER")
	processPhotoRequest := &aipb.ProcessFacecamRequest{
		Id:        request.UserId,
		CreatorId: request.CreatorId,
		Url:       request.FileURL,
	}

	a.logs.CustomLog(" ini adakah process face cam request", request.CreatorId)

	_, err := a.client.ProcessFacecam(ctx, processPhotoRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *aiAdapter) ProcessBulkPhoto(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo) error {
	log.Println("REQUESTED PROCESS  BULK PHOTO VIA GRPC TO AI SERVER")
	pbAIPhotos := make([]*aipb.AIPhoto, len(*photos))
	for i, photo := range *photos {
		pbAIPhotos[i] = &aipb.AIPhoto{
			Id:               photo.Id,
			CollectionUrl:    photo.CollectionUrl,
			OriginalFilename: photo.Title,
		}
		log.Println(pbAIPhotos[i])
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
