package converter

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
)

func UserProfileToResponse(userProfile *entity.UserProfile, profileUrl, coverUrl string) *model.UserProfileResponse {
	var (
		biographyPtr       *string
		profileUrlPtr      *string
		profileCoverUrlPtr *string
	)

	if userProfile.Biography.Valid {
		biographyPtr = &userProfile.Biography.String
	}
	if userProfile.ProfileUrl.Valid {
		profileUrlPtr = &userProfile.ProfileUrl.String
	}
	if userProfile.ProfileCoverUrl.Valid {
		profileCoverUrlPtr = &userProfile.ProfileCoverUrl.String
	}

	return &model.UserProfileResponse{
		Id:              userProfile.Id,
		UserId:          userProfile.UserId,
		BirthDate:       userProfile.BirthDate,
		Nickname:        userProfile.Nickname,
		Biography:       biographyPtr,
		ProfileUrl:      profileUrlPtr,
		ProfileCoverUrl: profileCoverUrlPtr,
		Similarity:      userProfile.Similarity,
		CreatedAt:       userProfile.CreatedAt,
		UpdatedAt:       userProfile.UpdatedAt,
	}
}

func UserProfileToResponseV2(userProfile *entity.UserProfile) *model.UserProfileResponse {
	return &model.UserProfileResponse{
		Id:              userProfile.Id,
		UserId:          userProfile.UserId,
		BirthDate:       userProfile.BirthDate,
		Nickname:        userProfile.Nickname,
		Biography:       nullable.SQLStringToPtr(userProfile.Biography),
		ProfileUrl:      nullable.SQLStringToPtr(userProfile.ProfileUrl),
		ProfileCoverUrl: nullable.SQLStringToPtr(userProfile.ProfileCoverUrl),
		Similarity:      userProfile.Similarity,
		CreatedAt:       userProfile.CreatedAt,
		UpdatedAt:       userProfile.UpdatedAt,
	}
}
