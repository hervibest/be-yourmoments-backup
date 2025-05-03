package usecase

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type PhotoUseCase interface {
	CreatePhoto(ctx context.Context, request *photopb.CreatePhotoRequest) error
	UpdatePhotoDetail(ctx context.Context, request *photopb.UpdatePhotoDetailRequest) error
	CreateBulkPhoto(ctx context.Context, request *photopb.CreateBulkPhotoRequest) error
	GetBulkPhotoDetail(ctx context.Context, request *model.GetBulkPhotoDetailRequest) (*model.GetBulkPhotoDetailResponse, error)
	GetPhotoFile(ctx context.Context, filename string) (io.ReadCloser, error)
}

type photoUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	userSimilarRepo repository.UserSimilarRepository
	creatorRepo     repository.CreatorRepository
	bulkPhotoRepo   repository.BulkPhotoRepository
	aiAdapter       adapter.AiAdapter
	storageAdapter  adapter.StorageAdapter
	logs            *logger.Log
}

func NewPhotoUseCase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository,
	userSimilarRepo repository.UserSimilarRepository,
	creatorRepo repository.CreatorRepository,
	bulkPhotoRepo repository.BulkPhotoRepository,
	aiAdapter adapter.AiAdapter, storageAdapter adapter.StorageAdapter,
	logs *logger.Log) PhotoUseCase {
	return &photoUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		userSimilarRepo: userSimilarRepo,
		creatorRepo:     creatorRepo,
		bulkPhotoRepo:   bulkPhotoRepo,
		aiAdapter:       aiAdapter,
		storageAdapter:  storageAdapter,
		logs:            logs,
	}
}

// ISSUE #1 : creator_id should be called form pb contract and come from AuthMiddleware (centralized auth)
func (u *photoUsecase) CreatePhoto(ctx context.Context, request *photopb.CreatePhotoRequest) error {
	log.Print(request.Photo.GetUserId())
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	newPhoto := &entity.Photo{
		Id:          request.GetPhoto().GetId(),
		CreatorId:   request.GetPhoto().GetCreatorId(),
		BulkPhotoId: nullable.GRPCtoSQLString(request.GetPhoto().GetBulkPhotoId()),
		Title:       request.GetPhoto().GetTitle(),
		CollectionUrl: sql.NullString{
			Valid:  true,
			String: request.GetPhoto().GetCollectionUrl(),
		},
		Price:    request.GetPhoto().GetPrice(),
		PriceStr: request.GetPhoto().GetPriceStr(),

		OriginalAt: request.GetPhoto().GetOriginalAt().AsTime(),
		CreatedAt:  request.GetPhoto().GetCreatedAt().AsTime(),
		UpdatedAt:  request.GetPhoto().GetUpdatedAt().AsTime(),

		Latitude:    nullable.GRPCtoSQLDouble(request.GetPhoto().GetLatitude()),
		Longitude:   nullable.GRPCtoSQLDouble(request.GetPhoto().GetLongitude()),
		Description: nullable.GRPCtoSQLString(request.GetPhoto().GetDescription()),
	}

	newPhoto, err = u.photoRepo.Create(tx, newPhoto)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new photo to database", err)
	}

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         newPhoto.Id,
		FileName:        request.GetPhoto().GetDetail().GetFileName(),
		FileKey:         request.GetPhoto().GetDetail().GetFileKey(),
		Size:            request.GetPhoto().GetDetail().GetSize(),
		Type:            request.GetPhoto().GetDetail().GetType(),
		Checksum:        request.GetPhoto().GetDetail().GetChecksum(),
		Width:           request.GetPhoto().GetDetail().GetWidth(),
		Height:          request.GetPhoto().GetDetail().GetHeight(),
		Url:             request.GetPhoto().GetDetail().GetUrl(),
		YourMomentsType: enum.YourMomentsType(request.GetPhoto().GetDetail().GetYourMomentsType()),
		CreatedAt:       request.GetPhoto().GetDetail().GetCreatedAt().AsTime(),
		UpdatedAt:       request.GetPhoto().GetDetail().GetUpdatedAt().AsTime(),
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

func (u *photoUsecase) UpdatePhotoDetail(ctx context.Context, request *photopb.UpdatePhotoDetailRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	photo := &entity.Photo{
		Id: request.GetPhotoDetail().GetPhotoId(),
		CompressedUrl: sql.NullString{
			Valid:  true,
			String: request.GetPhotoDetail().GetUrl(),
		}, UpdatedAt: time.Now(),
	}

	err = u.photoRepo.UpdateCompressedUrl(tx, photo)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update compressed url to photo in database", err)
	}

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         request.GetPhotoDetail().GetPhotoId(),
		FileName:        request.GetPhotoDetail().GetFileName(),
		FileKey:         request.GetPhotoDetail().GetFileKey(),
		Size:            request.GetPhotoDetail().GetSize(),
		Type:            "JPG",
		Checksum:        "1212",
		Height:          121,
		Width:           1212,
		Url:             request.GetPhotoDetail().GetUrl(),
		YourMomentsType: enum.YourMomentsType(request.GetPhotoDetail().GetYourMomentsType()),
		CreatedAt:       request.GetPhotoDetail().GetCreatedAt().AsTime(),
		UpdatedAt:       request.GetPhotoDetail().GetUpdatedAt().AsTime(),
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

// ISSUE #1 : creator_id should be called form pb contract and come from AuthMiddleware (centralized auth)
func (u *photoUsecase) CreateBulkPhoto(ctx context.Context, request *photopb.CreateBulkPhotoRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	bulkPhotoRepo := &entity.BulkPhoto{
		Id:              request.GetBulkPhoto().GetId(),
		CreatorId:       request.GetBulkPhoto().GetCreatorId(),
		BulkPhotoStatus: enum.BulkPhotoStatus(request.GetBulkPhoto().BulkPhotoStatus),
		CreatedAt:       request.GetBulkPhoto().GetCreatedAt().AsTime(),
		UpdatedAt:       request.GetBulkPhoto().GetUpdatedAt().AsTime(),
	}

	_, err = u.bulkPhotoRepo.Create(ctx, tx, bulkPhotoRepo)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert new bulk photo to database", err)
	}

	photos := make([]*entity.Photo, 0)
	photoDetails := make([]*entity.PhotoDetail, 0)

	for _, pbPhoto := range request.GetPhotos() {

		photo := &entity.Photo{
			Id:          pbPhoto.GetId(),
			CreatorId:   pbPhoto.GetCreatorId(),
			Title:       pbPhoto.GetTitle(),
			BulkPhotoId: nullable.GRPCtoSQLString(pbPhoto.GetBulkPhotoId()),
			CollectionUrl: sql.NullString{
				Valid:  true,
				String: pbPhoto.GetCollectionUrl(),
			},
			Price:    pbPhoto.GetPrice(),
			PriceStr: pbPhoto.GetPriceStr(),

			OriginalAt: pbPhoto.GetOriginalAt().AsTime(),
			CreatedAt:  pbPhoto.GetCreatedAt().AsTime(),
			UpdatedAt:  pbPhoto.GetUpdatedAt().AsTime(),

			Latitude:    nullable.GRPCtoSQLDouble(pbPhoto.GetLatitude()),
			Longitude:   nullable.GRPCtoSQLDouble(pbPhoto.GetLongitude()),
			Description: nullable.GRPCtoSQLString(pbPhoto.GetDescription()),
		}

		photoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         pbPhoto.Id,
			FileName:        pbPhoto.GetDetail().GetFileName(),
			FileKey:         pbPhoto.GetDetail().GetFileKey(),
			Size:            pbPhoto.GetDetail().GetSize(),
			Type:            pbPhoto.GetDetail().GetType(),
			Checksum:        pbPhoto.GetDetail().GetChecksum(),
			Width:           pbPhoto.GetDetail().GetWidth(),
			Height:          pbPhoto.GetDetail().GetHeight(),
			Url:             pbPhoto.GetDetail().GetUrl(),
			YourMomentsType: enum.YourMomentsType(pbPhoto.GetDetail().GetYourMomentsType()),
			CreatedAt:       pbPhoto.GetDetail().GetCreatedAt().AsTime(),
			UpdatedAt:       pbPhoto.GetDetail().GetUpdatedAt().AsTime(),
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

func (u *photoUsecase) GetBulkPhotoDetail(ctx context.Context, request *model.GetBulkPhotoDetailRequest) (*model.GetBulkPhotoDetailResponse, error) {
	bulkPhotoDetails, err := u.bulkPhotoRepo.FindDetailById(ctx, u.db, request.BulkPhotoId, request.CreatorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrUserNotFound, "Invalid bulk photo id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	if len(*bulkPhotoDetails) == 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrUserNotFound, "Invalid bulk photo id")
	}

	return converter.BulkPhotoDetailToResponse(bulkPhotoDetails), nil
}

func (u *photoUsecase) GetPhotoFile(ctx context.Context, filename string) (io.ReadCloser, error) {
	object, err := u.storageAdapter.GetFile(ctx, filename)
	if err != nil {
		if strings.Contains(err.Error(), "file not found") {
			return nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, "File not found")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to get photo file from minio storage", err)
	}

	return object, nil
}
