package usecase

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/enum"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/be-yourmoments/pb"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type UserSimilarUsecase interface {
	CreateUserSimilar(ctx context.Context, request *pb.CreateUserSimilarPhotoRequest) error
	CreateUserFacecam(ctx context.Context, request *pb.CreateUserSimilarFacecamRequest) error
}

type userSimilarUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	facecamRepo     repository.FacecamRepository
	userSimilarRepo repository.UserSimilarRepository
	logs            *logger.Log
}

func NewUserSimilarUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, logs *logger.Log) UserSimilarUsecase {
	return &userSimilarUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
		logs:            logs,
	}
}

func (u *userSimilarUsecase) CreateUserSimilar(ctx context.Context, request *pb.CreateUserSimilarPhotoRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	photo := &entity.Photo{
		Id:             request.GetPhotoDetail().PhotoId,
		IsThisYouURL:   sql.NullString{String: "", Valid: true},
		YourMomentsUrl: sql.NullString{String: request.GetPhotoDetail().GetUrl(), Valid: true},
		UpdatedAt:      time.Now(),
	}

	err = u.photoRepo.UpdateProcessedUrl(tx, photo)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update processed photo url in database", err)
	}

	//TODO TYPE , CHECKSUM, HEIGHT, WIOTH

	//YourMoments Type (AI Result)
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

	_, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		log.Println(err)
		return err
	}

	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		userSimilarPhoto := &entity.UserSimilarPhoto{
			Id:         ulid.Make().String(),
			PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
			UserId:     userSimilarPhotoRequest.GetUserId(),
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity().String()),
			CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
			UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
		}

		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
		log.Println("id : " + userSimilarPhoto.Id)
		log.Println("photo id : " + userSimilarPhoto.PhotoId)
		log.Println("user id : " + userSimilarPhoto.UserId)
		log.Println("similarity : " + userSimilarPhoto.Similarity)
	}

	err = u.userSimilarRepo.InsertOrUpdateByPhotoId(tx, request.GetPhotoDetail().PhotoId, &userSimilarPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update photo in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil

}

func (u *userSimilarUsecase) CreateUserFacecam(ctx context.Context, request *pb.CreateUserSimilarFacecamRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	facecam := &entity.Facecam{
		UserId:      request.GetFacecam().GetUserId(),
		IsProcessed: request.GetFacecam().GetIsProcessed(),
		UpdatedAt:   request.GetFacecam().GetUpdatedAt().AsTime(),
	}

	err = u.facecamRepo.UpdatedProcessedFacecam(tx, facecam)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update processed facecam url in database", err)
	}

	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		u.logs.Log("UPDATE UserSimilarPhoto from facecams")
		userSimilarPhoto := &entity.UserSimilarPhoto{
			Id:         ulid.Make().String(),
			PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
			UserId:     userSimilarPhotoRequest.GetUserId(),
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity().String()),
			CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
			UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
		}
		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
	}

	err = u.userSimilarRepo.InserOrUpdateByUserId(tx, request.GetFacecam().UserId, &userSimilarPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update user facececam in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
