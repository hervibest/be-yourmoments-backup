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

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type UserSimilarWorkerUseCase interface {
	CreateBulkUserSimilarPhotos(ctx context.Context, request *event.BulkUserSimilarPhotoEvent) error
	CreateUserFacecam(ctx context.Context, request *event.UserSimiliarFacecamEvent) error
	CreateUserSimilar(ctx context.Context, request *event.UserSimilarEvent) error
}

type userSimilarWorkerUseCase struct {
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

func NewUserSimilarWorkerUseCase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, bulkPhotoRepo repository.BulkPhotoRepository,
	userAdapter adapter.UserAdapter, photoProducer producer.PhotoProducer, logs *logger.Log) UserSimilarWorkerUseCase {
	return &userSimilarWorkerUseCase{
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

func (u *userSimilarWorkerUseCase) CreateUserSimilar(ctx context.Context, request *event.UserSimilarEvent) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	photo := &entity.Photo{
		Id:             request.PhotoDetail.PhotoID,
		IsThisYouURL:   sql.NullString{String: "", Valid: true},
		YourMomentsUrl: sql.NullString{String: request.PhotoDetail.Url, Valid: true},
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

	_, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		log.Println(err)
		return err
	}

	userIds := make([]string, 0, len(request.UserSimilarPhoto))
	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.UserSimilarPhoto))
	for _, userSimilarPhotoRequest := range request.UserSimilarPhoto {
		userSimilarPhoto := &entity.UserSimilarPhoto{
			PhotoId:    userSimilarPhotoRequest.PhotoID,
			UserId:     userSimilarPhotoRequest.UserID,
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.Similarity),
			CreatedAt:  *userSimilarPhotoRequest.CreatedAt,
			UpdatedAt:  *userSimilarPhotoRequest.UpdatedAt,
		}

		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
		log.Println("photo id : " + userSimilarPhoto.PhotoId)
		log.Println("user id : " + userSimilarPhoto.UserId)
		log.Println("similarity : ", userSimilarPhoto.Similarity)

		userIds = append(userIds, userSimilarPhotoRequest.UserID)
	}

	err = u.userSimilarRepo.InsertOrUpdateByPhotoId(tx, request.PhotoDetail.PhotoID, &userSimilarPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update photo in database", err)
	}

	err = u.photoRepo.AddPhotoTotal(ctx, tx, request.PhotoDetail.PhotoID, len(request.UserSimilarPhoto))
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

func (u *userSimilarWorkerUseCase) CreateUserFacecam(ctx context.Context, request *event.UserSimiliarFacecamEvent) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	facecam := &entity.Facecam{
		UserId:      request.Facecam.UserId,
		IsProcessed: request.Facecam.IsProcessed,
		UpdatedAt:   *request.Facecam.UpdatedAt,
	}

	err = u.facecamRepo.UpdatedProcessedFacecam(tx, facecam)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update processed facecam url in database", err)
	}

	photoIDs := make([]string, 0, len(request.UserSimilarPhoto))
	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.UserSimilarPhoto))
	for _, userSimilarPhotoRequest := range request.UserSimilarPhoto {
		u.logs.Log("UPDATE UserSimilarPhoto from facecams")
		userSimilarPhoto := &entity.UserSimilarPhoto{
			PhotoId:    userSimilarPhotoRequest.PhotoID,
			UserId:     userSimilarPhotoRequest.UserID,
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.Similarity),
			CreatedAt:  *userSimilarPhotoRequest.CreatedAt,
			UpdatedAt:  *userSimilarPhotoRequest.UpdatedAt,
		}
		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
		photoIDs = append(photoIDs, userSimilarPhotoRequest.PhotoID)

		u.logs.CustomLog("Ini adalalah photo ID", userSimilarPhoto.PhotoId)
		u.logs.CustomLog("Ini adalalah user ID", userSimilarPhoto.UserId)
		u.logs.CustomLog("Ini adalalah similarity", userSimilarPhoto.Similarity)
	}

	err = u.userSimilarRepo.InserOrUpdateByUserId(tx, request.Facecam.UserId, &userSimilarPhotos)
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

	if len(request.UserSimilarPhoto) != 0 {
		go func() {
			singleFacecamEvent := &event.SingleFacecamEvent{
				EventID:     uuid.NewString(),
				UserID:      request.Facecam.UserId,
				CountPhotos: len(request.UserSimilarPhoto),
			}
			if err := u.photoProducer.ProduceSingleFacecam(ctx, singleFacecamEvent); err != nil {
				u.logs.Error(err)
			}
		}()
	}

	return nil
}

func (u *userSimilarWorkerUseCase) CreateBulkUserSimilarPhotos(ctx context.Context, request *event.BulkUserSimilarPhotoEvent) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}
	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	// Update BulkPhoto entity
	bulkPhoto := &entity.BulkPhoto{
		Id:              request.BulkPhoto.Id,
		CreatorId:       request.BulkPhoto.CreatorId,
		BulkPhotoStatus: enum.BulkPhotoStatus(request.BulkPhoto.BulkPhotoStatus),
		UpdatedAt:       time.Now(),
	}
	_, err = u.bulkPhotoRepo.Update(ctx, tx, bulkPhoto)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update bulk photo entity in database", err)
	}

	allPhotos := make([]*entity.Photo, 0)

	for _, bulkUserSimilar := range request.BulkUserSimilarPhoto {
		photo := &entity.Photo{
			Id:             bulkUserSimilar.PhotoDetail.PhotoID,
			IsThisYouURL:   sql.NullString{String: "", Valid: true},
			YourMomentsUrl: sql.NullString{String: bulkUserSimilar.PhotoDetail.Url, Valid: true},
			UpdatedAt:      time.Now(),
		}
		allPhotos = append(allPhotos, photo)
	}

	err = u.photoRepo.UpdateProcessedUrlBulk(tx, allPhotos)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to bulk update photo processed url", err)
	}

	photoUserSimilarMap := make(map[string][]*entity.UserSimilarPhoto)

	for _, bulkUserSimilar := range request.BulkUserSimilarPhoto {
		newPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         bulkUserSimilar.PhotoDetail.PhotoID,
			FileName:        bulkUserSimilar.PhotoDetail.FileName,
			FileKey:         bulkUserSimilar.PhotoDetail.FileKey,
			Size:            bulkUserSimilar.PhotoDetail.Size,
			Type:            "JPG",
			Checksum:        "1212",
			Height:          121,
			Width:           1212,
			Url:             bulkUserSimilar.PhotoDetail.Url,
			YourMomentsType: enum.YourMomentsType(bulkUserSimilar.PhotoDetail.YourMomentsType),
			CreatedAt:       *bulkUserSimilar.PhotoDetail.CreatedAt,
			UpdatedAt:       *bulkUserSimilar.PhotoDetail.UpdatedAt,
		}

		_, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
		if err != nil {
			log.Println(err)
			return err
		}

		userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(bulkUserSimilar.UserSimilarPhoto))
		for _, userSimilarPhotoRequest := range bulkUserSimilar.UserSimilarPhoto {
			userSimilarPhotos = append(userSimilarPhotos, &entity.UserSimilarPhoto{
				PhotoId:    userSimilarPhotoRequest.PhotoID,
				UserId:     userSimilarPhotoRequest.UserID,
				Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.Similarity),
				CreatedAt:  *userSimilarPhotoRequest.CreatedAt,
				UpdatedAt:  *userSimilarPhotoRequest.UpdatedAt,
			})
		}

		photoUserSimilarMap[bulkUserSimilar.PhotoDetail.PhotoID] = userSimilarPhotos
	}

	for photoID, userSimilars := range photoUserSimilarMap {
		fmt.Println("Photo:", photoID)
		seen := map[string]bool{}
		for _, u := range userSimilars {
			key := fmt.Sprintf("%s-%s", photoID, u.UserId)
			if seen[key] {
				fmt.Println("⚠️  Duplicate pair detected:", key)
			}
			seen[key] = true
			fmt.Println("  User:", u.UserId)
		}
	}

	err = u.userSimilarRepo.InsertOrUpdateBulk(ctx, tx, photoUserSimilarMap)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to insert or update bulk user similar photos in database", err)
	}

	photoCountMap := u.countPhotosParallel(request.BulkUserSimilarPhoto)
	err = u.photoRepo.BulkAddPhotoTotals(ctx, tx, photoCountMap)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update add photos total in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	countMap := u.countUsersParallel(request.BulkUserSimilarPhoto)
	u.logs.Log(fmt.Sprintf("[USER][BULK USER SIMILAR PHOTO] Count user in photo service: %v", countMap))
	for id, count := range countMap {
		u.logs.Log(fmt.Sprintf("[USER][BULK USER SIMILAR PHOTO] User ID %s has count %d", id, count))
	}

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

func (u *userSimilarWorkerUseCase) countUsersParallel(datas []event.BulkUserSimilarPhoto) map[string]int32 {
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
		go func(part []event.BulkUserSimilarPhoto) {
			defer wg.Done()
			localCount := make(map[string]int32)

			for _, photo := range part {
				for _, user := range photo.UserSimilarPhoto {
					u.logs.Log(fmt.Sprintf("[CountUsersParallel] User ID %s in photo %s", user.UserID, photo.PhotoDetail.PhotoID))
					localCount[user.UserID]++
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

func (u *userSimilarWorkerUseCase) countPhotosParallel(datas []event.BulkUserSimilarPhoto) map[string]int32 {
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
		go func(part []event.BulkUserSimilarPhoto) {
			defer wg.Done()
			localCount := make(map[string]int32)

			for _, photo := range part {
				localCount[photo.PhotoDetail.PhotoID] += int32(len(photo.UserSimilarPhoto))
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
