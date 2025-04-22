package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type photoPreparedStmt struct {
	findByPhotoId *sqlx.Stmt
	findManyByIds *sqlx.Stmt
}

func newPhotoPreparedStmt(db *sqlx.DB) (*photoPreparedStmt, error) {
	findByPhotoIdStmt, err := db.Preparex("SELECT * FROM photos WHERE id = $1")
	if err != nil {
		return nil, err
	}

	findManyByIdsStmt, err := db.Preparex(`
	SELECT 
			id,
			creator_id,
			title,
			is_this_you_url,
			your_moments_url,
			price
		FROM photos AS p
		JOIN user_similar_photos AS up
		ON up.photo_id = p.id
		WHERE id = ANY($1) AND up.user_id = $2 AND p.owned_by_user_id`)
	if err != nil {
		return nil, err
	}

	return &photoPreparedStmt{
		findByPhotoId: findByPhotoIdStmt,
		findManyByIds: findManyByIdsStmt,
	}, nil
}

type PhotoRepository interface {
	Create(tx Querier, photo *entity.Photo) (*entity.Photo, error)
	FindByPhotoId(ctx context.Context, tx Querier, photoId string) (*entity.Photo, error)
	UpdateProcessedUrl(tx Querier, photo *entity.Photo) error
	UpdateCompressedUrl(tx Querier, photo *entity.Photo) error
	GetPhotosByIDs(ctx context.Context, userId string, ids []string) (*[]*entity.Photo, error)
	UpdatePhotoOwnerByPhotoIds(ctx context.Context, tx Querier, ownerID string, photoIDs []string) error
	// UpdateClaimedPhoto(ctx context.Context, db Querier, photo *entity.Photo) error
	// UpdatePhotoStatus(ctx context.Context, db Querier, photo *entity.Photo) error
}

type photoRepository struct {
	photoPreparedStmt *photoPreparedStmt
}

func NewPhotoRepository(db *sqlx.DB) (PhotoRepository, error) {
	photoPreparedStmt, err := newPhotoPreparedStmt(db)
	if err != nil {
		log.Print("error initialize photo statement : ", err)
		return nil, err
	}

	return &photoRepository{
		photoPreparedStmt: photoPreparedStmt,
	}, nil
}

func (r *photoRepository) Create(tx Querier, photo *entity.Photo) (*entity.Photo, error) {
	query := `INSERT INTO photos 
			  (id, creator_id, title, collection_url, price, price_str, latitude, longitude, description, original_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := tx.Exec(query, photo.Id, photo.CreatorId, photo.Title, photo.CollectionUrl, photo.Price, photo.PriceStr, photo.Latitude, photo.Longitude, photo.Description,
		photo.OriginalAt, photo.CreatedAt, photo.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert photo: %w", err)
	}

	return photo, nil
}

func (r *photoRepository) UpdateCompressedUrl(tx Querier, photo *entity.Photo) error {
	log.Println("Updated accesed")
	query := `UPDATE photos 
			  SET compressed_url = $1, updated_at = $2
			  WHERE id = $3`

	_, err := tx.Exec(query, photo.CompressedUrl, photo.UpdatedAt, photo.Id)
	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}

func (r *photoRepository) UpdateProcessedUrl(tx Querier, photo *entity.Photo) error {
	log.Println("Updated accesed")
	query := `UPDATE photos 
			  SET is_this_you_url = $1, your_moments_url = $2, updated_at = $3 
			  WHERE id = $4`

	_, err := tx.Exec(query, photo.IsThisYouURL, photo.YourMomentsUrl, photo.UpdatedAt, photo.Id)
	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}

func (r *photoRepository) FindByPhotoId(ctx context.Context, tx Querier, photoId string) (*entity.Photo, error) {
	photo := new(entity.Photo)

	row := r.photoPreparedStmt.findByPhotoId.QueryRowxContext(ctx, photoId)
	if err := row.StructScan(photo); err != nil {
		return nil, err
	}

	return photo, nil
}

func (r *photoRepository) GetPhotosByIDs(ctx context.Context, userId string, ids []string) (*[]*entity.Photo, error) {
	photos := make([]*entity.Photo, 0)
	if err := r.photoPreparedStmt.findManyByIds.SelectContext(ctx, &photos, ids, userId); err != nil {
		return nil, err
	}
	return &photos, nil
}

func (r *photoRepository) UpdatePhotoOwnerByPhotoIds(ctx context.Context, tx Querier, ownerID string, photoIDs []string) error {
	query := "UPDATE photos SET owned_by_user_id = $1 WHERE id = ANY($2)"
	_, err := tx.ExecContext(ctx, query, ownerID, photoIDs)
	if err != nil {
		log.Println("Error happen in updatePhotoOwner", err)
	}

	return nil
}
