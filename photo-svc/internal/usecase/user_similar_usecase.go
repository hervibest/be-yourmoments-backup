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
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type UserSimilarUsecase interface {
	CreateUserSimilar(ctx context.Context, request *photopb.CreateUserSimilarPhotoRequest) error
	CreateUserFacecam(ctx context.Context, request *photopb.CreateUserSimilarFacecamRequest) error
	CreateBulkUserSimilarPhotos(ctx context.Context, request *photopb.CreateBulkUserSimilarPhotoRequest) error
}

type userSimilarUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	facecamRepo     repository.FacecamRepository
	userSimilarRepo repository.UserSimilarRepository
	bulkPhotoRepo   repository.BulkPhotoRepository
	userAdapter     adapter.UserAdapter
	logs            *logger.Log
}

func NewUserSimilarUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, bulkPhotoRepo repository.BulkPhotoRepository,
	userAdapter adapter.UserAdapter, logs *logger.Log) UserSimilarUsecase {
	return &userSimilarUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
		bulkPhotoRepo:   bulkPhotoRepo,
		userAdapter:     userAdapter,
		logs:            logs,
	}
}

func (u *userSimilarUsecase) CreateUserSimilar(ctx context.Context, request *photopb.CreateUserSimilarPhotoRequest) error {
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
			PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
			UserId:     userSimilarPhotoRequest.GetUserId(),
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity().String()),
			CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
			UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
		}

		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
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

func (u *userSimilarUsecase) CreateUserFacecam(ctx context.Context, request *photopb.CreateUserSimilarFacecamRequest) error {
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

func (u *userSimilarUsecase) CreateBulkUserSimilarPhotos(ctx context.Context, request *photopb.CreateBulkUserSimilarPhotoRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}
	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	// Update BulkPhoto entity
	bulkPhoto := &entity.BulkPhoto{
		Id:              request.GetBulkPhoto().GetId(),
		CreatorId:       request.GetBulkPhoto().GetCreatorId(),
		BulkPhotoStatus: enum.BulkPhotoStatus(request.GetBulkPhoto().GetBulkPhotoStatus()),
		UpdatedAt:       time.Now(),
	}
	_, err = u.bulkPhotoRepo.Update(ctx, tx, bulkPhoto)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update bulk photo entity in database", err)
	}

	allPhotos := make([]*entity.Photo, 0)

	for _, bulkUserSimilar := range request.GetBulkUserSimilarPhoto() {
		photo := &entity.Photo{
			Id:             bulkUserSimilar.GetPhotoDetail().GetPhotoId(),
			IsThisYouURL:   sql.NullString{String: "", Valid: true},
			YourMomentsUrl: sql.NullString{String: bulkUserSimilar.GetPhotoDetail().GetUrl(), Valid: true},
			UpdatedAt:      time.Now(),
		}
		allPhotos = append(allPhotos, photo)
	}

	err = u.photoRepo.UpdateProcessedUrlBulk(tx, allPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to bulk update photo processed url", err)
	}

	photoUserSimilarMap := make(map[string][]*entity.UserSimilarPhoto)

	for _, bulkUserSimilar := range request.GetBulkUserSimilarPhoto() {
		newPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         bulkUserSimilar.GetPhotoDetail().GetPhotoId(),
			FileName:        bulkUserSimilar.GetPhotoDetail().GetFileName(),
			FileKey:         bulkUserSimilar.GetPhotoDetail().GetFileKey(),
			Size:            bulkUserSimilar.GetPhotoDetail().GetSize(),
			Type:            "JPG",
			Checksum:        "1212",
			Height:          121,
			Width:           1212,
			Url:             bulkUserSimilar.GetPhotoDetail().GetUrl(),
			YourMomentsType: enum.YourMomentsType(bulkUserSimilar.GetPhotoDetail().GetYourMomentsType()),
			CreatedAt:       bulkUserSimilar.GetPhotoDetail().GetCreatedAt().AsTime(),
			UpdatedAt:       bulkUserSimilar.GetPhotoDetail().GetUpdatedAt().AsTime(),
		}

		_, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
		if err != nil {
			log.Println(err)
			return err
		}

		userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(bulkUserSimilar.GetUserSimilarPhoto()))
		for _, userSimilarPhotoRequest := range bulkUserSimilar.GetUserSimilarPhoto() {
			userSimilarPhotos = append(userSimilarPhotos, &entity.UserSimilarPhoto{
				PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
				UserId:     userSimilarPhotoRequest.GetUserId(),
				Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity().String()),
				CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
				UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
			})
		}

		photoUserSimilarMap[bulkUserSimilar.GetPhotoDetail().GetPhotoId()] = userSimilarPhotos
	}

	err = u.userSimilarRepo.InsertOrUpdateBulk(ctx, tx, photoUserSimilarMap)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update bulk user similar photos in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	go func() {
		if _, err := u.userAdapter.SendBulkPhotoNotification(ctx, request.GetBulkUserSimilarPhoto()); err != nil {
			u.logs.Error(err)
		}
	}()
	return nil
}
