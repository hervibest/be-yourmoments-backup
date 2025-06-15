package producer

import (
	"context"
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model/event"
)

type UploadProducer interface {
	ProcessBulkPhoto(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo) error
	ProcessFacecam(ctx context.Context, request *model.ProcessFacecam) error
	ProcessPhoto(ctx context.Context, request *model.ProcessPhoto) error
}

type uploadProducer struct {
	messagingAdapter adapter.MessagingAdapter
	logs             logger.Log
}

func NewUploadProducer(messagingAdapter adapter.MessagingAdapter, logs logger.Log) UploadProducer {
	return &uploadProducer{
		messagingAdapter: messagingAdapter,
		logs:             logs,
	}
}

func (a *uploadProducer) ProcessPhoto(ctx context.Context, request *model.ProcessPhoto) error {
	msg := &event.ProcessPhotoMessage{
		PhotoID:          request.PhotoId,
		CreatorID:        request.CreatorId,
		URL:              request.FileURL,
		OriginalFilename: request.OriginalFilename,
		Timestamp:        time.Now().Unix(),
	}

	if err := a.messagingAdapter.Publish(ctx, "AI.PHOTO.PROCESS", msg); err != nil {
		a.logs.Log("Failed to publish ProcessPhoto message")
		return err
	}

	a.logs.Log("Successfully published ProcessPhoto message")
	return nil
}

func (a *uploadProducer) ProcessFacecam(ctx context.Context, request *model.ProcessFacecam) error {
	msg := &event.ProcessFacecamMessage{
		UserID:    request.UserId,
		CreatorID: request.CreatorId,
		URL:       request.FileURL,
		Timestamp: time.Now().Unix(),
	}

	if err := a.messagingAdapter.Publish(ctx, "AI.FACECAM.PROCESS", msg); err != nil {
		a.logs.Log("Failed to publish ProcessFacecam message")
		return err
	}

	a.logs.Log("Successfully published ProcessFacecam message")
	return nil
}

func (a *uploadProducer) ProcessBulkPhoto(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo) error {
	bulkItems := make([]*event.BulkPhotoItem, len(*photos))
	for i, photo := range *photos {
		bulkItems[i] = &event.BulkPhotoItem{
			ID:               photo.Id,
			CollectionURL:    photo.CollectionUrl,
			OriginalFilename: photo.Title,
		}
	}

	msg := &event.ProcessBulkPhotoMessage{
		BulkPhotoID: bulkPhoto.Id,
		CreatorID:   bulkPhoto.CreatorId,
		Photos:      bulkItems,
		Timestamp:   time.Now().Unix(),
	}

	if err := a.messagingAdapter.Publish(ctx, "AI.BULK_PHOTO.PROCESS", msg); err != nil {
		a.logs.Log("Failed to publish ProcessBulkPhoto message")
		return err
	}

	a.logs.Log("Successfully published ProcessBulkPhoto message")
	return nil
}
