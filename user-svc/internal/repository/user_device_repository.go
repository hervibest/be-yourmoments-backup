package repository

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"

	"github.com/lib/pq"
)

type UserDeviceRepository interface {
	Create(ctx context.Context, tx Querier, userDevice *entity.UserDevice) (*entity.UserDevice, error)
	FetchFCMTokensFromPostgre(ctx context.Context, tx Querier, userIDs []string) (*[]*entity.UserDevice, error)
	DeleteByUserID(ctx context.Context, tx Querier, userID string) error
}
type userDeviceRepository struct{}

func NewUserDeviceRepository() UserDeviceRepository {
	return &userDeviceRepository{}
}

func (r *userDeviceRepository) Create(ctx context.Context, tx Querier, userDevice *entity.UserDevice) (*entity.UserDevice, error) {
	query := `
	INSERT INTO user_devices
	(id, user_id, token, platform, created_at) 
	VALUES ($1, $2, $3, $4, $5) 
	ON CONFLICT (user_id, token) DO UPDATE
	SET platform = EXCLUDED.platform,
    created_at = current_timestamp;
	`
	_, err := tx.ExecContext(ctx, query, userDevice.Id, userDevice.UserId, userDevice.Token, userDevice.Platform, userDevice.CreatedAt)
	if err != nil {
		return nil, err
	}
	return userDevice, err
}

func (r *userDeviceRepository) FetchFCMTokensFromPostgre(ctx context.Context, tx Querier, userIDs []string) (*[]*entity.UserDevice, error) {
	userDevices := make([]*entity.UserDevice, 0)
	query := `
	SELECT * from user_devices 
	WHERE user_id = ANY($1)
	`

	if err := tx.SelectContext(ctx, &userDevices, query, pq.Array(userIDs)); err != nil {
		return nil, err
	}

	return &userDevices, nil
}

func (r *userDeviceRepository) DeleteByUserID(ctx context.Context, tx Querier, userID string) error {
	query := `
	DELETE 
	FROM user_devices
	WHERE user_id = $1
	`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
