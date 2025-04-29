package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	errorcode "github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum/error"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/repository"

	"github.com/oklog/ulid/v2"
)

type UserUseCase interface {
	GetUserProfile(ctx context.Context, userId string) (*model.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, request *model.RequestUpdateUserProfile) (*model.UserProfileResponse, error)
	UpdateUserProfileImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error)
	UpdateUserCoverImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error)
	GetPublicUserChat(ctx context.Context, request *model.RequestGetAllPublicUser) (*[]*model.GetAllPublicUserResponse, *model.PageMetadata, error)
}

type userUseCase struct {
	db                    repository.BeginTx
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
	userImageRepository   repository.UserImageRepository
	uploadAdapter         adapter.UploadAdapter
	logs                  *logger.Log
}

func NewUserUseCase(db repository.BeginTx, userRepository repository.UserRepository, userProfileRepository repository.UserProfileRepository,
	userImageRepository repository.UserImageRepository, uploadAdapter adapter.UploadAdapter, logs *logger.Log) UserUseCase {
	return &userUseCase{
		db:                    db,
		userRepository:        userRepository,
		userProfileRepository: userProfileRepository,
		userImageRepository:   userImageRepository,
		uploadAdapter:         uploadAdapter,
		logs:                  logs,
	}
}

func (u *userUseCase) GetUserProfile(ctx context.Context, userId string) (*model.UserProfileResponse, error) {
	userProfile, err := u.userProfileRepository.FindByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "invalid user id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find user profile by user id", err)
	}

	userImages, err := u.userImageRepository.FindByUserProfId(ctx, userProfile.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "invalid user profile id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find user image by user profile id", err)
	}

	profileUrl, coverUrl, err := u.getUserImageUrl(ctx, userImages)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to get user image url", err)
	}

	return converter.UserProfileToResponse(userProfile, profileUrl, coverUrl), nil
}

func (u *userUseCase) getUserImageUrl(ctx context.Context, userImages *[]*entity.UserImage) (string, string, error) {
	var (
		profileUrl string
		coverUrl   string
		err        error
	)

	for _, userImage := range *userImages {
		if userImage.ImageType == enum.ImageTypeProfile {
			profileUrl, err = u.uploadAdapter.GetPresignedUrl(ctx, userImage.FileKey)
			if err != nil {
				return "", "", fmt.Errorf("get presigned url image type profile : %+v", err)
			}
		} else if userImage.ImageType == enum.ImageTypeCover {
			coverUrl, err = u.uploadAdapter.GetPresignedUrl(ctx, userImage.FileKey)
			if err != nil {
				return "", "", fmt.Errorf("get presigned url image type cover : %+v", err)
			}
		}
	}

	return profileUrl, coverUrl, nil
}

func (u *userUseCase) UpdateUserProfile(ctx context.Context, request *model.RequestUpdateUserProfile) (*model.UserProfileResponse, error) {
	now := time.Now()
	userProfile := &entity.UserProfile{
		UserId:    request.UserId,
		BirthDate: request.BirthDate,
		Nickname:  request.Nickname,
		Biography: sql.NullString{
			Valid:  true,
			String: request.Biography,
		},
		UpdatedAt: &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	userProfile, err = u.userProfileRepository.Update(ctx, tx, userProfile)
	if err != nil {
		return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid user profile id")
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	userImages, err := u.userImageRepository.FindByUserProfId(ctx, userProfile.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid user profile id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find user image repository by user profile  id", err)
	}

	profileUrl, coverUrl, err := u.getUserImageUrl(ctx, userImages)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to get user image url", err)
	}

	return converter.UserProfileToResponse(userProfile, profileUrl, coverUrl), nil
}

func (u *userUseCase) UpdateUserProfileImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error) {
	ok, err := u.updateUserImage(ctx, file, userProfId, enum.ImageTypeProfile)
	if err != nil {
		return false, err
	}

	return ok, err
}

func (u *userUseCase) UpdateUserCoverImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error) {
	ok, err := u.updateUserImage(ctx, file, userProfId, enum.ImageTypeCover)
	if err != nil {
		return false, err
	}

	return ok, err
}

func (u *userUseCase) updateUserImage(ctx context.Context, file *multipart.FileHeader, userProfId string, imageType enum.ImageTypeEnum) (bool, error) {
	prevUserProfileImage, errRepo := u.userImageRepository.FindByUserProfIdAndType(ctx, userProfId, string(imageType))
	if errRepo != nil && !errors.Is(errRepo, sql.ErrNoRows) {
		return false, helper.WrapInternalServerError(u.logs, "failed to find user image by user profile id", errRepo)
	}

	uploadFile, err := file.Open()
	if err != nil {
		return false, helper.WrapInternalServerError(u.logs, "failed to open file user", errRepo)
	}
	defer uploadFile.Close()

	upload, err := u.uploadAdapter.UploadFile(ctx, file, uploadFile, fmt.Sprintf("user/profile/%s", userProfId))
	if err != nil {
		return false, helper.WrapInternalServerError(u.logs, "failed to upload file", errRepo)
	}

	now := time.Now()

	if errors.Is(errRepo, sql.ErrNoRows) {
		newUserProfileImage := &entity.UserImage{
			Id:            ulid.Make().String(),
			UserProfileId: userProfId,
			FileName:      upload.Filename,
			FileKey:       upload.FileKey,
			ImageType:     imageType,
			Size:          upload.Size,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}

		tx, err := repository.BeginTxx(u.db, ctx, u.logs)
		if err != nil {
			return false, err
		}

		defer func() {
			repository.Rollback(err, tx, ctx, u.logs)
		}()

		_, err = u.userImageRepository.Create(ctx, tx, newUserProfileImage)
		if err != nil {
			return false, helper.WrapInternalServerError(u.logs, "failed to create user profile image", err)
		}

		if err := repository.Commit(tx, u.logs); err != nil {
			return false, err
		}

	} else {
		updatedUserProfileImage := &entity.UserImage{
			UserProfileId: userProfId,
			FileName:      upload.Filename,
			FileKey:       upload.FileKey,
			ImageType:     imageType,
			Size:          upload.Size,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}

		tx, err := repository.BeginTxx(u.db, ctx, u.logs)
		if err != nil {
			return false, err
		}

		defer func() {
			repository.Rollback(err, tx, ctx, u.logs)
		}()

		_, err = u.userImageRepository.Update(ctx, tx, updatedUserProfileImage)
		if err != nil {
			return false, helper.WrapInternalServerError(u.logs, "failed to update user profile image", err)
		}

		_, err = u.uploadAdapter.DeleteFile(ctx, prevUserProfileImage.FileName)
		if err != nil {
			return false, helper.WrapInternalServerError(u.logs, "failed to delete user profile image from minio", err)
		}

		if err := repository.Commit(tx, u.logs); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (u *userUseCase) GetPublicUserChat(ctx context.Context, request *model.RequestGetAllPublicUser) (*[]*model.GetAllPublicUserResponse, *model.PageMetadata, error) {
	userPublicChat, pageMetadata, err := u.userRepository.FindAllPublicChat(ctx, u.db, request.Page, request.Size, request.Username)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all public chat (user list)", err)
	}

	responses := make([]*model.GetAllPublicUserResponse, 0)

	//TODO perlu dipertimbangkan buat kode eksekusinya berapa lama karena perlu ubah ke url
	for _, entity := range userPublicChat {
		response := &model.GetAllPublicUserResponse{
			UserId:   entity.UserId,
			Username: entity.Username,
		}

		if entity.FileKey.Valid {
			profileUrl, err := u.uploadAdapter.GetPresignedUrl(ctx, entity.FileKey.String)
			if err != nil {
				return nil, nil, helper.WrapInternalServerError(u.logs, "failed to get presigned url", err)
			}
			response.ProfileUrl = profileUrl
		}

		responses = append(responses, response)
	}

	return &responses, pageMetadata, nil
}

// func (u *userUseCase) UpdateUserProfileCover(ctx context.Context, userId string) (*model.UserProfileResponse, error) {
// 	now := time.Now()
// 	userProfile := &entity.UserProfile{
// 		UserId: userId,
// 		ProfileCoverUrl: sql.NullString{
// 			String: request.ProfileCoverUrl,
// 		},

// 		UpdatedAt: &now,
// 	}

// 	tx, err := u.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	userProfile, err = u.userProfileRepository.UpdateUserProfileCover(ctx, tx, userProfile)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return converter.UserProfileToResponse(userProfile), nil
// }
