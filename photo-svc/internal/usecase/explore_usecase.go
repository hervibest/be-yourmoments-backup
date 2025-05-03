package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	oteltrace "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ExploreUseCase interface {
	GetUserCart(ctx context.Context, request *model.GetAllCartRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error)
	GetUserExploreSimilar(ctx context.Context, request *model.GetAllExploreSimilarRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error)
	GetUserFavorite(ctx context.Context, request *model.GetAllFavoriteRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error)
	GetUserWishlist(ctx context.Context, request *model.GetAllWishlistRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error)
	UserAddCart(ctx context.Context, request *model.UserAddCartRequest) error
	UserAddFavorite(ctx context.Context, request *model.UserAddFavoriteRequest) error
	UserAddWishlist(ctx context.Context, request *model.UserAddWishlistRequest) error
	UserDeleteCart(ctx context.Context, request *model.UserDeleteCartReqeust) error
	UserDeleteFavorite(ctx context.Context, request *model.UserDeleteFavoriteReqeust) error
	UserDeleteWishlist(ctx context.Context, request *model.UserDeleteWishlistReqeust) error
}

type exploreUseCase struct {
	db                *sqlx.DB
	exploreRepository repository.ExploreRepository
	photoRepository   repository.PhotoRepository
	CDNAdapter        adapter.CDNAdapter
	tracer            trace.Tracer
	logs              *logger.Log
}

func NewExploreUseCase(db *sqlx.DB, exploreRepository repository.ExploreRepository, photoRepository repository.PhotoRepository,
	CDNAdapter adapter.CDNAdapter, tracer trace.Tracer, logs *logger.Log) ExploreUseCase {
	return &exploreUseCase{
		db:                db,
		exploreRepository: exploreRepository,
		photoRepository:   photoRepository,
		CDNAdapter:        CDNAdapter,
		tracer:            tracer,
		logs:              logs,
	}
}

func (u *exploreUseCase) GetUserExploreSimilar(ctx context.Context, request *model.GetAllExploreSimilarRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error) {
	// Memulai span baru jika memang ingin menambahkan tracing khusus di use case
	_, span := u.tracer.Start(ctx, "exploreUseCase.GetUserExploreSimilar", oteltrace.WithAttributes(attribute.String("user.id", request.UserId)))
	defer span.End()

	explores, pageMetadata, err := u.exploreRepository.FindAllExploreSimilar(ctx, u.db, request.Page, request.Size, request.Similarity, request.UserId)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all explore similar in database", err)
	}

	return converter.ExploresToResponses(&explores, u.CDNAdapter.GenerateCDN), pageMetadata, nil
}

func (u *exploreUseCase) GetUserWishlist(ctx context.Context, request *model.GetAllWishlistRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error) {
	_, span := u.tracer.Start(ctx, "exploreUseCase.GetUserWishlist", oteltrace.WithAttributes(attribute.String("user.id", request.UserId)))
	defer span.End()

	explores, pageMetadata, err := u.exploreRepository.FindAllUserWishlist(ctx, u.db, request.Page, request.Size, request.Similarity, request.UserId)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all user wishlist photo in database", err)
	}

	return converter.ExploresToResponses(&explores, u.CDNAdapter.GenerateCDN), pageMetadata, nil
}

func (u *exploreUseCase) UserAddWishlist(ctx context.Context, request *model.UserAddWishlistRequest) error {
	_, err := u.photoRepository.FindByPhotoId(ctx, u.db, request.PhotoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.exploreRepository.UserAddWishlist(ctx, tx, request.Similarity, request.PhotoId, request.UserId); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to add photo wishlist in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *exploreUseCase) UserDeleteWishlist(ctx context.Context, request *model.UserDeleteWishlistReqeust) error {
	_, err := u.photoRepository.FindByPhotoId(ctx, u.db, request.PhotoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.exploreRepository.UserAddWishlist(ctx, tx, request.Similarity, request.PhotoId, request.UserId); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to delete wishlist in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *exploreUseCase) GetUserFavorite(ctx context.Context, request *model.GetAllFavoriteRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error) {
	_, span := u.tracer.Start(ctx, "exploreUseCase.GetUserFavorite", oteltrace.WithAttributes(attribute.String("user.id", request.UserId)))
	defer span.End()

	explores, pageMetadata, err := u.exploreRepository.FindAllUserFavorite(ctx, u.db, request.Page, request.Size, request.Similarity, request.UserId)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all user favorite photo in database", err)
	}

	return converter.ExploresToResponses(&explores, u.CDNAdapter.GenerateCDN), pageMetadata, nil
}

func (u *exploreUseCase) UserAddFavorite(ctx context.Context, request *model.UserAddFavoriteRequest) error {
	_, err := u.photoRepository.FindByPhotoId(ctx, u.db, request.PhotoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.exploreRepository.UserAddFavorite(ctx, tx, request.Similarity, request.PhotoId, request.UserId); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to add user favorite photo in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *exploreUseCase) UserDeleteFavorite(ctx context.Context, request *model.UserDeleteFavoriteReqeust) error {
	_, err := u.photoRepository.FindByPhotoId(ctx, u.db, request.PhotoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.exploreRepository.UserDeleteFavorite(ctx, tx, request.Similarity, request.PhotoId, request.UserId); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to delete user favorite photo in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *exploreUseCase) GetUserCart(ctx context.Context, request *model.GetAllCartRequest) (*[]*model.ExploreUserSimilarResponse, *model.PageMetadata, error) {

	_, span := u.tracer.Start(ctx, "exploreUseCase.GetUserCart", oteltrace.WithAttributes(attribute.String("user.id", request.UserId)))
	defer span.End()

	explores, pageMetadata, err := u.exploreRepository.FindAllUserCart(ctx, u.db, request.Page, request.Size, request.Similarity, request.UserId)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all user cart photo in database", err)
	}

	return converter.ExploresToResponses(&explores, u.CDNAdapter.GenerateCDN), pageMetadata, nil
}

func (u *exploreUseCase) UserAddCart(ctx context.Context, request *model.UserAddCartRequest) error {
	_, err := u.photoRepository.FindByPhotoId(ctx, u.db, request.PhotoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.exploreRepository.UserAddCart(ctx, tx, request.Similarity, request.PhotoId, request.UserId); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to add user cart photo in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *exploreUseCase) UserDeleteCart(ctx context.Context, request *model.UserDeleteCartReqeust) error {
	_, err := u.photoRepository.FindByPhotoId(ctx, u.db, request.PhotoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find photo by photo id in database", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.exploreRepository.UserDeleteCart(ctx, tx, request.Similarity, request.PhotoId, request.UserId); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to delete user cart photo in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
