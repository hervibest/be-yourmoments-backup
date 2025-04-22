package usecase

import (
	"be-yourmoments/photo-svc/internal/adapter"
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/repository"
	"context"

	"github.com/be-yourmoments/pb"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type FacecamUseCase interface {
	CreateFacecam(ctx context.Context, request *pb.CreateFacecamRequest) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type facecamUseCase struct {
	db              *sqlx.DB
	facecamRepo     repository.FacecamRepository
	userSimilarRepo repository.UserSimilarRepository
	aiAdapter       adapter.AiAdapter
	uploadAdapter   adapter.UploadAdapter
	logs            *logger.Log
}

func NewFacecamUseCase(db *sqlx.DB, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, aiAdapter adapter.AiAdapter,
	uploadAdapter adapter.UploadAdapter, logs *logger.Log) FacecamUseCase {
	return &facecamUseCase{
		db:              db,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
		aiAdapter:       aiAdapter,
		uploadAdapter:   uploadAdapter,
		logs:            logs,
	}
}

func (u *facecamUseCase) CreateFacecam(ctx context.Context, request *pb.CreateFacecamRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	newPhoto := &entity.Facecam{
		Id:       ulid.Make().String(),
		UserId:   request.GetFacecam().GetUserId(),
		FileName: request.GetFacecam().GetFileName(),
		FileKey:  request.GetFacecam().GetFileKey(),
		Title:    request.GetFacecam().GetTitle(),

		Size: request.GetFacecam().GetSize(),
		Url:  request.GetFacecam().GetUrl(),

		OriginalAt: request.GetFacecam().GetOriginalAt().AsTime(),
		CreatedAt:  request.GetFacecam().GetCreatedAt().AsTime(),
		UpdatedAt:  request.GetFacecam().GetUpdatedAt().AsTime(),
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
