package model

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"
)

type RequestGetUserProfile struct {
	UserId string `json:"user_id" validate:"required"`
}

type RequestUpdateUserProfile struct {
	UserId       string     `validate:"required"`
	BirthDate    *time.Time `json:"-"`
	BirthDateStr string     `json:"birth_date" validate:"required"`
	Nickname     string     `json:"nickname" validate:"required"`
	Biography    string     `json:"biography" validate:"required"`
}

type UserProfileResponse struct {
	Id              string     `json:"id"`
	UserId          string     `json:"user_id"`
	BirthDate       string     `json:"birth_date,omitempty"`
	Nickname        string     `json:"nickname"`
	Biography       *string    `json:"biography"`
	ProfileUrl      *string    `json:"profile_url"`
	ProfileCoverUrl *string    `json:"profile_cover_url"`
	Similarity      uint       `json:"similarity"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

type RequestGetAllPublicUser struct {
	Username string `json:"username" validate:"required"`
	Page     int    `json:"page" validate:"required"`
	Size     int    `json:"size" validate:"required"`
}

type GetAllPublicUserResponse struct {
	UserId     string `json:"user_id"`
	Username   string `json:"username"`
	ProfileUrl string `json:"profile_url,omitempty"`
}

type RequestUpdateSimilarity struct {
	Similarity enum.SimilarityLevelEnum `json:"similarity" validate:"required,gte=1,lte=9"`
	UserID     string                   `json:"user_id" validate:"required"`
}

type UpdateSeimilarityResponse struct {
	Similarity enum.SimilarityLevelEnum `json:"similarity"`
}
