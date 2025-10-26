package producer

import (
	"context"
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model/event"
)

type UploadProducer interface {
	ProcessBulkPhoto(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo) error
	ProcessFacecam(ctx context.Context, request *model.ProcessFacecam) error
	ProcessPhoto(ctx context.Context, request *model.ProcessPhoto) error
	CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error
	CreatePhotos(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo, photoDetails *[]*entity.PhotoDetail) error
	UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error
	CreateFacecam(ctx context.Context, facecam *entity.Facecam) error
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

func (a *uploadProducer) CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error {
	photoDetailEvent := &event.PhotoDetail{
		Id:              facecam.Id,
		PhotoId:         facecam.PhotoId,
		FileName:        facecam.FileName,
		FileKey:         facecam.FileKey,
		Size:            facecam.Size,
		Type:            facecam.Type,
		Checksum:        facecam.Checksum,
		Width:           int64(facecam.Width),
		Height:          int64(facecam.Height),
		Url:             facecam.Url,
		YourMomentsType: string(facecam.YourMomentsType),
		CreatedAt:       &facecam.CreatedAt,
		UpdatedAt:       &facecam.UpdatedAt,
	}

	photoEvent := &event.Photo{
		ID:            photo.Id,
		UserID:        photo.UserId,
		CreatorID:     photo.CreatorId,
		Title:         photo.Title,
		CollectionURL: photo.CollectionUrl,
		Price:         int32(photo.Price),
		PriceStr:      photo.PriceStr,
		OriginalAt:    &photo.OriginalAt,
		CreatedAt:     &photo.CreatedAt,
		UpdatedAt:     &photo.UpdatedAt,
		PhotoDetail:   *photoDetailEvent,
		Latitude:      nullable.SQLFloat64ToPtr(photo.Latitude),
		Longitude:     nullable.SQLFloat64ToPtr(photo.Longitude),
		Description:   nullable.SQLStringToPtr(photo.Description),
	}

	createPhotoEvent := &event.CreatePhotoEvent{Photo: *photoEvent}
	if err := a.messagingAdapter.Publish(ctx, "upload.single.photo", createPhotoEvent); err != nil {
		a.logs.Log("Failed to publish ProcessBulkPhoto message")
		return err
	}

	a.logs.Log("Successfully published ProcessBulkPhoto message")
	return nil
}

func (a *uploadProducer) CreatePhotos(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo, photoDetails *[]*entity.PhotoDetail) error {
	bulkPhotoEvent := &event.BulkPhoto{
		Id:              bulkPhoto.Id,
		CreatorId:       bulkPhoto.CreatorId,
		BulkPhotoStatus: string(bulkPhoto.BulkPhotoStatus),
		CreatedAt:       &bulkPhoto.CreatedAt,
		UpdatedAt:       &bulkPhoto.UpdatedAt,
	}

	photoEvents := make([]event.Photo, len(*photos))
	for i, photo := range *photos {
		detail := (*photoDetails)[i]
		photoEvents[i] = event.Photo{
			ID:            photo.Id,
			UserID:        photo.UserId,
			CreatorID:     photo.CreatorId,
			Title:         photo.Title,
			CollectionURL: photo.CollectionUrl,
			Price:         int32(photo.Price),
			PriceStr:      photo.PriceStr,
			OriginalAt:    &photo.OriginalAt,
			CreatedAt:     &photo.CreatedAt,
			UpdatedAt:     &photo.UpdatedAt,
			PhotoDetail: event.PhotoDetail{
				Id:              detail.Id,
				PhotoId:         detail.PhotoId,
				FileName:        detail.FileName,
				FileKey:         detail.FileKey,
				Size:            detail.Size,
				Type:            detail.Type,
				Checksum:        detail.Checksum,
				Width:           int64(detail.Width),
				Height:          int64(detail.Height),
				Url:             detail.Url,
				YourMomentsType: string(detail.YourMomentsType),
				CreatedAt:       &detail.CreatedAt,
				UpdatedAt:       &detail.UpdatedAt,
			},
			Latitude:    nullable.SQLFloat64ToPtr(photo.Latitude),
			Longitude:   nullable.SQLFloat64ToPtr(photo.Longitude),
			Description: nullable.SQLStringToPtr(photo.Description),
		}
	}

	createBulkPhotoEvent := &event.CreateBulkPhotoEvent{
		BulkPhoto: *bulkPhotoEvent,
		Photos:    photoEvents,
	}

	if err := a.messagingAdapter.Publish(ctx, "upload.bulk.photo", createBulkPhotoEvent); err != nil {
		a.logs.Log("Failed to publish CreatePhotos message")
		return err
	}

	a.logs.Log("Successfully published CreatePhotos message")
	return nil
}

func (a *uploadProducer) UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error {
	photoDetailEvent := &event.PhotoDetail{
		Id:              facecam.Id,
		PhotoId:         facecam.PhotoId,
		FileName:        facecam.FileName,
		FileKey:         facecam.FileKey,
		Size:            facecam.Size,
		Type:            facecam.Type,
		Checksum:        facecam.Checksum,
		Width:           int64(facecam.Width),
		Height:          int64(facecam.Height),
		Url:             facecam.Url,
		YourMomentsType: string(facecam.YourMomentsType),
		CreatedAt:       &facecam.CreatedAt,
		UpdatedAt:       &facecam.UpdatedAt,
	}

	updatePhotoDetailEvent := &event.UpdatePhotoDetailEvent{PhotoDetail: *photoDetailEvent}
	if err := a.messagingAdapter.Publish(ctx, "upload.update.photo", updatePhotoDetailEvent); err != nil {
		a.logs.Log("Failed to publish UpdatePhotoDetail message")
		return err
	}

	a.logs.Log("Successfully published UpdatePhotoDetail message")
	return nil
}

func (a *uploadProducer) CreateFacecam(ctx context.Context, facecam *entity.Facecam) error {
	facecamEvent := &event.CreateFacecamEvent{
		Facecam: event.Facecam{
			Id:         facecam.Id,
			UserId:     facecam.UserId,
			FileName:   facecam.FileName,
			FileKey:    facecam.FileKey,
			Title:      facecam.Title,
			Size:       facecam.Size,
			Url:        facecam.Url,
			OriginalAt: &facecam.OriginalAt,
			CreatedAt:  &facecam.CreatedAt,
			UpdatedAt:  &facecam.UpdatedAt,
		},
	}

	if err := a.messagingAdapter.Publish(ctx, "upload.single.facecam", facecamEvent); err != nil {
		a.logs.Log("Failed to publish CreateFacecam message")
		return err
	}

	a.logs.Log("Successfully published CreateFacecam message")
	return nil
}
