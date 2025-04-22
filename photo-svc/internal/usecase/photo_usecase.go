package usecase

import (
	"be-yourmoments/photo-svc/internal/adapter"
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/enum"
	errorcode "be-yourmoments/photo-svc/internal/enum/error"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/helper/nullable"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/be-yourmoments/pb"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type PhotoUsecase interface {
	CreatePhoto(ctx context.Context, request *pb.CreatePhotoRequest) error
	UpdatePhotoDetail(ctx context.Context, request *pb.UpdatePhotoDetailRequest) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type photoUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	userSimilarRepo repository.UserSimilarRepository
	creatorRepo     repository.CreatorRepository
	aiAdapter       adapter.AiAdapter
	uploadAdapter   adapter.UploadAdapter
	logs            *logger.Log
}

func NewPhotoUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository,
	userSimilarRepo repository.UserSimilarRepository,
	creatorRepo repository.CreatorRepository,
	aiAdapter adapter.AiAdapter, uploadAdapter adapter.UploadAdapter,
	logs *logger.Log) PhotoUsecase {
	return &photoUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		userSimilarRepo: userSimilarRepo,
		creatorRepo:     creatorRepo,
		aiAdapter:       aiAdapter,
		uploadAdapter:   uploadAdapter,
		logs:            logs,
	}
}

func (u *photoUsecase) CreatePhoto(ctx context.Context, request *pb.CreatePhotoRequest) error {
	log.Print(request.Photo.GetUserId())
	creator, err := u.creatorRepo.FindByUserId(ctx, request.Photo.GetUserId())
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid user ids")
		}
		return helper.WrapInternalServerError(u.logs, "error find creator by user id google", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	newPhoto := &entity.Photo{
		Id:        request.GetPhoto().GetId(),
		CreatorId: creator.Id,
		Title:     request.GetPhoto().GetTitle(),
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

func (u *photoUsecase) UpdatePhotoDetail(ctx context.Context, request *pb.UpdatePhotoDetailRequest) error {
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

// func (u *photoUsecase) ClaimPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

// 	tx, err := u.db.Begin(ctx)
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:            req.Id,
// 		OwnedByUserId: req.UserId,
// 		Status:        "Claimed",
// 		UpdatedAt:     time.Now(),
// 	}

// 	err = u.photoRepo.UpdateClaimedPhoto(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	// err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
// 	// if err != nil {
// 	// 	return err, err
// 	// }

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	// Process photo service will be executed asyncronously by goroutine

// 	return nil, nil

// }

// func (u *photoUsecase) CancelClaimPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

// 	tx, err := u.db.Begin(ctx)
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:            req.Id,
// 		OwnedByUserId: "",
// 		Status:        "Unclaimed",
// 		UpdatedAt:     time.Now(),
// 	}

// 	err = u.photoRepo.UpdateClaimedPhoto(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	// err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
// 	// if err != nil {
// 	// 	return err, err
// 	// }

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	// Process photo service will be executed asyncronously by goroutine

// 	return nil, nil

// }

// func (u *photoUsecase) UpdateBuyyedPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

// 	tx, err := u.db.Begin(ctx)
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:        req.Id,
// 		Status:    "Owned",
// 		UpdatedAt: time.Now(),
// 	}

// 	err = u.photoRepo.UpdatePhotoStatus(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	err = u.userSimilarRepo.DeleteSimilarUsers(ctx, tx, req.Id)
// 	if err != nil {
// 		return err, err
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	// Process photo service will be executed asyncronously by goroutine

// 	return nil, nil

// }
