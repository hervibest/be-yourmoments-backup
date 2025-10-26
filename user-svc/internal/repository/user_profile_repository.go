package repository

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"

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
	CreateWithProfileUrl(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	Create(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	Update(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	FindByUserId(ctx context.Context, userId string) (*entity.UserProfile, error)
	UpdateSimilarity(ctx context.Context, tx Querier, similarity enum.SimilarityLevelEnum, userID string) error
	UpdateImageURL(ctx context.Context, tx Querier, url, userProfId string, imageType enum.ImageTypeEnum) error
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
	query := `
	INSERT INTO user_profiles 
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
	query := `
	UPDATE 
		user_profiles
	SET
		birth_date = $1, 
		nickname = $2, 
		biography = $3, 
		updated_at = $4 
	WHERE 
		user_id = $5 
	RETURNING *`

	if err := tx.GetContext(ctx, userProfile, query, userProfile.BirthDate, userProfile.Nickname, userProfile.Biography, userProfile.UpdatedAt, userProfile.UserId); err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (r *userProfileRepository) UpdateImageURL(ctx context.Context, tx Querier, url, userProfId string, imageType enum.ImageTypeEnum) error {
	query := "UPDATE user_profiles SET "
	if imageType == enum.ImageTypeProfile {
		query += "profile_url = $1"
	} else {
		query += "profile_cover_url = $1"
	}
	query += " WHERE id = $2"
	_, err := tx.ExecContext(ctx, query, url, userProfId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userProfileRepository) UpdateSimilarity(ctx context.Context, tx Querier, similarity enum.SimilarityLevelEnum, userID string) error {
	query := `UPDATE user_profiles set similarity = $1 WHERE user_id = $2 `
	_, err := tx.ExecContext(ctx, query, similarity, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userProfileRepository) FindByUserId(ctx context.Context, userId string) (*entity.UserProfile, error) {
	userProfile := new(entity.UserProfile)
	row := r.userProfilePreparedStmt.findByUserId.QueryRowxContext(ctx, userId)
	if err := row.StructScan(userProfile); err != nil {

		return nil, err
	}

	return userProfile, nil
}
