package repository

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"

	"github.com/jmoiron/sqlx"
)

type userProfilePreparedStmt struct {
	findById     *sqlx.Stmt
	findByUserId *sqlx.Stmt
}

func newUserProfilePreraredStmt(db *sqlx.DB) (*userProfilePreparedStmt, error) {
	findByUserIdStmt, err := db.Preparex("SELECT * FROM user_profiles  WHERE user_id = $1")
	if err != nil {
		return nil, err
	}

	findByIdStmt, err := db.Preparex("SELECT * FROM user_profiles WHERE id = $1")
	if err != nil {
		return nil, err
	}

	return &userProfilePreparedStmt{
		findById:     findByIdStmt,
		findByUserId: findByUserIdStmt,
	}, nil
}

type UserProfileRepository interface {
	Close() error
	CreateWithProfileUrl(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	Create(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	Update(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	FindByUserId(ctx context.Context, userId string) (*entity.UserProfile, error)

	// UpdateUserProfileImage(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	// UpdateUserProfileCover(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
}

type userProfileRepository struct {
	userProfilePreparedStmt *userProfilePreparedStmt
}

func NewUserProfileRepository(db *sqlx.DB) (UserProfileRepository, error) {

	userProfilePreparedStmt, err := newUserProfilePreraredStmt(db)
	if err != nil {
		return nil, err
	}

	return &userProfileRepository{
		userProfilePreparedStmt: userProfilePreparedStmt,
	}, nil
}

func (r *userProfileRepository) Close() error {
	if err := r.userProfilePreparedStmt.findById.Close(); err != nil {
		return err
	}

	if err := r.userProfilePreparedStmt.findByUserId.Close(); err != nil {
		return err
	}

	return nil
}

func (r *userProfileRepository) CreateWithProfileUrl(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
	query := `INSERT INTO user_profiles 
	(id, user_id, nickname, profile_url, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, userProfile.Id, userProfile.UserId, userProfile.Nickname, userProfile.ProfileUrl, userProfile.CreatedAt, userProfile.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (r *userProfileRepository) Create(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
	query := `INSERT INTO user_profiles 
	(id, user_id, birth_date, nickname, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, userProfile.Id, userProfile.UserId, userProfile.BirthDate, userProfile.Nickname, userProfile.CreatedAt, userProfile.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
	query := `UPDATE user_profiles set birth_date = $1, nickname = $2, 
	biography = $3, updated_at = $4 WHERE user_id = $5 RETURNING *`

	row := tx.QueryRowxContext(ctx, query, userProfile.BirthDate, userProfile.Nickname, userProfile.Biography, userProfile.UpdatedAt, userProfile.UserId)
	if err := row.StructScan(userProfile); err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (r *userProfileRepository) FindByUserId(ctx context.Context, userId string) (*entity.UserProfile, error) {
	userProfile := new(entity.UserProfile)
	row := r.userProfilePreparedStmt.findByUserId.QueryRowxContext(ctx, userId)
	if err := row.StructScan(userProfile); err != nil {

		return nil, err
	}

	return userProfile, nil
}

// func (r *userProfileRepository) UpdateUserProfileImage(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
// 	query := `UPDATE user_profiles set profile_url = $1, updated_at = $2 WHERE user_id = $3`
// 	_, err := tx.ExecContext(ctx, query, userProfile.ProfileUrl, userProfile.UpdatedAt, userProfile.UserId)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update user profile image: %w", err)
// 	}

// 	return userProfile, nil
// }

// func (r *userProfileRepository) UpdateUserProfileCover(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
// 	query := `UPDATE user_profiles set profile_cover_url = $1, updated_at = $2 WHERE user_id = $3`
// 	_, err := tx.ExecContext(ctx, query, userProfile.ProfileCoverUrl, userProfile.UpdatedAt, userProfile.UserId)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update user profile cover: %w", err)
// 	}

// 	return userProfile, nil
// }
