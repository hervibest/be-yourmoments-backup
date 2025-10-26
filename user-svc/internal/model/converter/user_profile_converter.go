package converter

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
)

func UserProfileToResponse(userProfile *entity.UserProfile) *model.UserProfileResponse {
	return &model.UserProfileResponse{
		Id:              userProfile.Id,
		UserId:          userProfile.UserId,
		BirthDate:       nullable.TimeToString(userProfile.BirthDate, time.DateOnly),
		Nickname:        userProfile.Nickname,
		Biography:       nullable.SQLStringToPtr(userProfile.Biography),
		ProfileUrl:      nullable.SQLStringToPtr(userProfile.ProfileUrl),
		ProfileCoverUrl: nullable.SQLStringToPtr(userProfile.ProfileCoverUrl),
		Similarity:      userProfile.Similarity,
		CreatedAt:       userProfile.CreatedAt,
		UpdatedAt:       userProfile.UpdatedAt,
	}
}
