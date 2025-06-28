package usecase

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type PhotoWorkerUseCase interface {
	CreatePhoto(ctx context.Context, request *event.CreatePhotoEvent) error
	UpdatePhotoDetail(ctx context.Context, request *event.UpdatePhotoDetailEvent) error
	CreateBulkPhoto(ctx context.Context, request *event.CreateBulkPhotoEvent) error
}

type photoWorkerUseCase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	userSimilarRepo repository.UserSimilarRepository
	creatorRepo     repository.CreatorRepository
	bulkPhotoRepo   repository.BulkPhotoRepository
	// aiAdapter       adapter.AiAdapter
	storageAdapter adapter.StorageAdapter
	CDNAdapter     adapter.CDNAdapter
	logs           *logger.Log
}

func NewPhotoWorkerUseCase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository,
	userSimilarRepo repository.UserSimilarRepository,
	creatorRepo repository.CreatorRepository,
	bulkPhotoRepo repository.BulkPhotoRepository,
	// aiAdapter adapter.AiAdapter,
	storageAdapter adapter.StorageAdapter,
	CDNAdapter adapter.CDNAdapter, logs *logger.Log) PhotoWorkerUseCase {
	return &photoWorkerUseCase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		userSimilarRepo: userSimilarRepo,
		creatorRepo:     creatorRepo,
		bulkPhotoRepo:   bulkPhotoRepo,
		// aiAdapter:       aiAdapter,
		CDNAdapter:     CDNAdapter,
		storageAdapter: storageAdapter,
		logs:           logs,
	}
}

func (u *photoWorkerUseCase) CreatePhoto(ctx context.Context, request *event.CreatePhotoEvent) error {
	log.Print(request.Photo.UserID)
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	newPhoto := &entity.Photo{
		Id:          request.Photo.ID,
		CreatorId:   request.Photo.CreatorID,
		BulkPhotoId: nullable.ToSQLString(request.Photo.BulkPhotoID),
		Title:       request.Photo.Title,
		CollectionUrl: sql.NullString{
			Valid:  true,
			String: request.Photo.CollectionURL,
		},
		Price:    request.Photo.Price,
		PriceStr: request.Photo.PriceStr,

		OriginalAt: *request.Photo.OriginalAt,
		CreatedAt:  *request.Photo.CreatedAt,
		UpdatedAt:  *request.Photo.UpdatedAt,

		Latitude:    nullable.ToSQLFloat64(request.Photo.Latitude),
		Longitude:   nullable.ToSQLFloat64(request.Photo.Longitude),
		Description: nullable.ToSQLString(request.Photo.Description),
	}

	newPhoto, err = u.photoRepo.Create(tx, newPhoto)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new photo to database", err)
	}

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         newPhoto.Id,
		FileName:        request.Photo.PhotoDetail.FileName,
		FileKey:         request.Photo.PhotoDetail.FileKey,
		Size:            request.Photo.PhotoDetail.Size,
		Type:            request.Photo.PhotoDetail.Type,
		Checksum:        request.Photo.PhotoDetail.Checksum,
		Width:           int32(request.Photo.PhotoDetail.Width),
		Height:          int32(request.Photo.PhotoDetail.Height),
		Url:             request.Photo.PhotoDetail.Url,
		YourMomentsType: enum.YourMomentsType(request.Photo.PhotoDetail.YourMomentsType),
		CreatedAt:       *request.Photo.PhotoDetail.CreatedAt,
		UpdatedAt:       *request.Photo.PhotoDetail.UpdatedAt,
	}

	newPhotoDetail, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new photo detail to database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil

}

func (u *photoWorkerUseCase) UpdatePhotoDetail(ctx context.Context, request *event.UpdatePhotoDetailEvent) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	photo := &entity.Photo{
		Id: request.PhotoDetail.PhotoID,
		CompressedUrl: sql.NullString{
			Valid:  true,
			String: request.PhotoDetail.Url,
		}, UpdatedAt: time.Now(),
	}

	err = u.photoRepo.UpdateCompressedUrl(tx, photo)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update compressed url to photo in database", err)
	}

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         request.PhotoDetail.PhotoID,
		FileName:        request.PhotoDetail.FileName,
		FileKey:         request.PhotoDetail.FileKey,
		Size:            request.PhotoDetail.Size,
		Type:            "JPG",
		Checksum:        "1212",
		Height:          121,
		Width:           1212,
		Url:             request.PhotoDetail.Url,
		YourMomentsType: enum.YourMomentsType(request.PhotoDetail.YourMomentsType),
		CreatedAt:       *request.PhotoDetail.CreatedAt,
		UpdatedAt:       *request.PhotoDetail.UpdatedAt,
	}

	newPhotoDetail, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update insert new photo detail in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil

}

func (u *photoWorkerUseCase) CreateBulkPhoto(ctx context.Context, request *event.CreateBulkPhotoEvent) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	bulkPhotoRepo := &entity.BulkPhoto{
		Id:              request.BulkPhoto.Id,
		CreatorId:       request.BulkPhoto.CreatorId,
		BulkPhotoStatus: enum.BulkPhotoStatus(request.BulkPhoto.BulkPhotoStatus),
		CreatedAt:       *request.BulkPhoto.CreatedAt,
		UpdatedAt:       *request.BulkPhoto.UpdatedAt,
	}

	_, err = u.bulkPhotoRepo.Create(ctx, tx, bulkPhotoRepo)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new bulk photo to database", err)
	}

	photos := make([]*entity.Photo, 0)
	photoDetails := make([]*entity.PhotoDetail, 0)

	for _, photoEvent := range request.Photos {

		photo := &entity.Photo{
			Id:          photoEvent.ID,
			CreatorId:   photoEvent.CreatorID,
			Title:       photoEvent.Title,
			BulkPhotoId: nullable.ToSQLString(photoEvent.BulkPhotoID),
			CollectionUrl: sql.NullString{
				Valid:  true,
				String: photoEvent.CollectionURL,
			},
			Price:    photoEvent.Price,
			PriceStr: photoEvent.PriceStr,

			OriginalAt: *photoEvent.OriginalAt,
			CreatedAt:  *photoEvent.CreatedAt,
			UpdatedAt:  *photoEvent.UpdatedAt,

			Latitude:    nullable.ToSQLFloat64(photoEvent.Latitude),
			Longitude:   nullable.ToSQLFloat64(photoEvent.Longitude),
			Description: nullable.ToSQLString(photoEvent.Description),
		}

		photoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         photoEvent.ID,
			FileName:        photoEvent.PhotoDetail.FileName,
			FileKey:         photoEvent.PhotoDetail.FileKey,
			Size:            photoEvent.PhotoDetail.Size,
			Type:            photoEvent.PhotoDetail.Type,
			Checksum:        photoEvent.PhotoDetail.Checksum,
			Width:           int32(photoEvent.PhotoDetail.Width),
			Height:          int32(photoEvent.PhotoDetail.Height),
			Url:             photoEvent.PhotoDetail.Url,
			YourMomentsType: enum.YourMomentsType(photoEvent.PhotoDetail.YourMomentsType),
			CreatedAt:       *photoEvent.PhotoDetail.CreatedAt,
			UpdatedAt:       *photoEvent.PhotoDetail.UpdatedAt,
		}

		photos = append(photos, photo)
		photoDetails = append(photoDetails, photoDetail)
	}

	_, err = u.photoRepo.BulkCreate(ctx, tx, photos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new photos to database", err)
	}

	_, err = u.photoDetailRepo.BulkCreate(ctx, tx, photoDetails)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new photo details to database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
