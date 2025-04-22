package converter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/model"
)

func UserProfileToResponse(userProfile *entity.UserProfile, profileUrl, coverUrl string) *model.UserProfileResponse {
	var (
		biographyPtr       *string
		profileUrlPtr      *string
		profileCoverUrlPtr *string
		similarityPtr      *string
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
	if userProfile.Similarity.Valid {
		profileCoverUrlPtr = &userProfile.Similarity.String
	}

	return &model.UserProfileResponse{
		Id:              userProfile.Id,
		UserId:          userProfile.UserId,
		BirthDate:       userProfile.BirthDate,
		Nickname:        userProfile.Nickname,
		Biography:       biographyPtr,
		ProfileUrl:      profileUrlPtr,
		ProfileCoverUrl: profileCoverUrlPtr,
		Similarity:      similarityPtr,
		CreatedAt:       userProfile.CreatedAt,
		UpdatedAt:       userProfile.UpdatedAt,
	}
}
