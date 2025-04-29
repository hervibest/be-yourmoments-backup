package converter

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
)

func UserToResponse(u *entity.User) *model.UserResponse {
	var (
		emailPtr  *string
		phonePtr  *string
		googlePtr *string
	)

	if u.Email.Valid {
		emailPtr = &u.Email.String
	}
	if u.PhoneNumber.Valid {
		phonePtr = &u.PhoneNumber.String
	}
	if u.GoogleId.Valid {
		googlePtr = &u.GoogleId.String
	}

	return &model.UserResponse{
		Id:       u.Id,
		Username: u.Username,
		Email:    emailPtr,
		// EmailVerifiedAt:       u.EmailVerifiedAt,
		PhoneNumber: phonePtr,
		// PhoneNumberVerifiedAt: u.PhoneNumberVerifiedAt,
		GoogleId:  googlePtr,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
