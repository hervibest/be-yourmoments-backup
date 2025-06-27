package converter

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
)

func UserToResponse(u *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Id:       u.Id,
		Username: u.Username,
		Email:    nullable.SQLStringToPtr(u.Email),
		// EmailVerifiedAt:       u.EmailVerifiedAt,
		PhoneNumber: nullable.SQLStringToPtr(u.PhoneNumber),
		// PhoneNumberVerifiedAt: u.PhoneNumberVerifiedAt,
		GoogleId:  nullable.SQLStringToPtr(u.GoogleId),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
