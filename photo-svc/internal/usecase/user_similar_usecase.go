package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	producer "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/gateway/messaging"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
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
	photoProducer   producer.PhotoProducer
	logs            *logger.Log
}

func NewUserSimilarUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, bulkPhotoRepo repository.BulkPhotoRepository,
	userAdapter adapter.UserAdapter, photoProducer producer.PhotoProducer, logs *logger.Log) UserSimilarUsecase {
	return &userSimilarUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
		bulkPhotoRepo:   bulkPhotoRepo,
		userAdapter:     userAdapter,
		photoProducer:   photoProducer,
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

	//TODO ISSUE #3 TYPE , CHECKSUM, HEIGHT, WIOTH

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

	userIds := make([]string, 0, len(request.GetUserSimilarPhoto()))
	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		userSimilarPhoto := &entity.UserSimilarPhoto{
			PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
			UserId:     userSimilarPhotoRequest.GetUserId(),
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity()),
			CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
			UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
		}

		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
		log.Println("photo id : " + userSimilarPhoto.PhotoId)
		log.Println("user id : " + userSimilarPhoto.UserId)
		log.Println("similarity : ", userSimilarPhoto.Similarity)

		userIds = append(userIds, userSimilarPhotoRequest.GetUserId())
	}

	err = u.userSimilarRepo.InsertOrUpdateByPhotoId(tx, request.GetPhotoDetail().PhotoId, &userSimilarPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update photo in database", err)
	}

	err = u.photoRepo.AddPhotoTotal(ctx, tx, request.GetPhotoDetail().PhotoId, len(request.GetUserSimilarPhoto()))
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update photo in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	if len(userIds) != 0 {
		go func() {
			singlePhotoEvent := &event.SinglePhotoEvent{
				EventID: uuid.NewString(),
				UserIDs: userIds,
			}
			if err := u.photoProducer.ProduceSinglePhoto(ctx, singlePhotoEvent); err != nil {
				u.logs.Error(err)
			}
		}()

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

	photoIDs := make([]string, 0, len(request.GetUserSimilarPhoto()))
	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		u.logs.Log("UPDATE UserSimilarPhoto from facecams")
		userSimilarPhoto := &entity.UserSimilarPhoto{
			PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
			UserId:     userSimilarPhotoRequest.GetUserId(),
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity()),
			CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
			UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
		}
		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
		photoIDs = append(photoIDs, userSimilarPhotoRequest.PhotoId)
	}

	err = u.userSimilarRepo.InserOrUpdateByUserId(tx, request.GetFacecam().UserId, &userSimilarPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update user facececam in database", err)
	}

	err = u.photoRepo.BulkIncrementTotal(ctx, tx, photoIDs)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to bulk increment total photos in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	if len(request.GetUserSimilarPhoto()) != 0 {
		go func() {
			singleFacecamEvent := &event.SingleFacecamEvent{
				EventID:     uuid.NewString(),
				UserID:      request.GetFacecam().GetUserId(),
				CountPhotos: len(request.GetUserSimilarPhoto()),
			}
			if err := u.photoProducer.ProduceSingleFacecam(ctx, singleFacecamEvent); err != nil {
				u.logs.Error(err)
			}
		}()
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
				Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity()),
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

	photoCountMap := u.countPhotosParallel(request.GetBulkUserSimilarPhoto())
	err = u.photoRepo.BulkAddPhotoTotals(ctx, tx, photoCountMap)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update add photos total in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	countMap := u.countUsersParallel(request.GetBulkUserSimilarPhoto())

	if countMap != nil {
		go func() {
			bulkPhotoEvent := &event.BulkPhotoEvent{
				EventID:      uuid.NewString(),
				UserCountMap: countMap,
			}
			if err := u.photoProducer.ProduceBulkPhoto(ctx, bulkPhotoEvent); err != nil {
				u.logs.Error(err)
			}
		}()
	}

	return nil
}

func (u *userSimilarUsecase) countUsersParallel(datas []*photopb.BulkUserSimilarPhoto) map[string]int32 {
	u.logs.Log("[CountUsersParallel] Count user in photo service")
	countMap := make(map[string]int32)
	var mu sync.Mutex

	numWorkers := runtime.NumCPU()
	chunkSize := (len(datas) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(datas))
		if start >= len(datas) {
			continue
		}

		wg.Add(1)
		go func(part []*photopb.BulkUserSimilarPhoto) {
			defer wg.Done()
			localCount := make(map[string]int32)

			for _, photo := range part {
				for _, user := range photo.GetUserSimilarPhoto() {
					localCount[user.UserId]++
				}
			}

			// Merge localCount ke global countMap
			mu.Lock()
			for id, cnt := range localCount {
				countMap[id] += cnt
			}
			mu.Unlock()
		}(datas[start:end])
	}

	wg.Wait()

	return countMap
}

func (u *userSimilarUsecase) countPhotosParallel(datas []*photopb.BulkUserSimilarPhoto) map[string]int32 {
	u.logs.Log("[CountUsersParallel] Count user in photo service")
	countMap := make(map[string]int32)
	var mu sync.Mutex

	numWorkers := runtime.NumCPU()
	chunkSize := (len(datas) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(datas))
		if start >= len(datas) {
			continue
		}

		wg.Add(1)
		go func(part []*photopb.BulkUserSimilarPhoto) {
			defer wg.Done()
			localCount := make(map[string]int32)

			for _, photo := range part {
				localCount[photo.PhotoDetail.PhotoId] += int32(len(photo.GetUserSimilarPhoto()))
			}

			// Merge localCount ke global countMap
			mu.Lock()
			for id, cnt := range localCount {
				countMap[id] += cnt
				u.logs.Log(fmt.Sprintf("COUNT MAP ID %s memiliki CNT %d", id, cnt))
			}
			mu.Unlock()
		}(datas[start:end])
	}

	wg.Wait()

	return countMap
}
