package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
)

type UserDeviceRepository interface {
	Create(ctx context.Context, tx Querier, userDevice *entity.UserDevice) (*entity.UserDevice, error)
}
type userDeviceRepository struct{}

func NewUserDeviceRepository() UserDeviceRepository {
	return &userDeviceRepository{}
}

func (r *userDeviceRepository) Create(ctx context.Context, tx Querier, userDevice *entity.UserDevice) (*entity.UserDevice, error) {
	query := `
	INSERT INTO user_devices
	(id, user_id, token, platform, created_at) 
	VALUES $1, $2, $3, $4, $5
	`
	_, err := tx.ExecContext(ctx, query, userDevice.Id, userDevice.UserId, userDevice.Token, userDevice.Platform, userDevice.CreatedAt)
	if err != nil {
		return nil, err
	}
	return userDevice, err
}
