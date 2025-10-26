package usecase

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type FacecamUseCaseWorker interface {
	CreateFacecam(ctx context.Context, request *event.CreateFacecamEvent) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type facecamUseCaseWorker struct {
	db              *sqlx.DB
	facecamRepo     repository.FacecamRepository
	userSimilarRepo repository.UserSimilarRepository
	// aiAdapter       adapter.AiAdapter
	storageAdapter adapter.StorageAdapter
	logs           *logger.Log
}

func NewFacecamUseCaseWorker(db *sqlx.DB, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, storageAdapter adapter.StorageAdapter,
	logs *logger.Log) FacecamUseCaseWorker {
	return &facecamUseCaseWorker{
		db:              db,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
		// aiAdapter:       aiAdapter,
		storageAdapter: storageAdapter,
		logs:           logs,
	}
}

func (u *facecamUseCaseWorker) CreateFacecam(ctx context.Context, request *event.CreateFacecamEvent) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	newPhoto := &entity.Facecam{
		Id:       ulid.Make().String(),
		UserId:   request.Facecam.UserId,
		FileName: request.Facecam.FileName,
		FileKey:  request.Facecam.FileKey,
		Title:    request.Facecam.Title,

		Size: request.Facecam.Size,
		Url:  request.Facecam.Url,

		OriginalAt: *request.Facecam.OriginalAt,
		CreatedAt:  *request.Facecam.CreatedAt,
		UpdatedAt:  *request.Facecam.UpdatedAt,
	}

	newPhoto, err = u.facecamRepo.Create(tx, newPhoto)
	if err != nil {
		return err
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil

}
