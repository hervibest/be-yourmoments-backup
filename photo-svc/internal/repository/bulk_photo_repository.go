package repository

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
)

type BulkPhotoRepository interface {
	Create(ctx context.Context, tx Querier, bulkPhoto *entity.BulkPhoto) (*entity.BulkPhoto, error)
	FindDetailById(ctx context.Context, tx Querier, bulkPhotoID, creatorID string) (*[]*entity.BulkPhotoDetail, error)
	Update(ctx context.Context, tx Querier, bulkPhoto *entity.BulkPhoto) (*entity.BulkPhoto, error)
}
type bulkPhotoRepository struct{}

func NewBulkPhotoRepository() BulkPhotoRepository {
	return &bulkPhotoRepository{}
}

func (r *bulkPhotoRepository) Create(ctx context.Context, tx Querier, bulkPhoto *entity.BulkPhoto) (*entity.BulkPhoto, error) {
	query := "INSERT INTO bulk_photos (id, creator_id, bulk_photo_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.ExecContext(ctx, query, bulkPhoto.Id, bulkPhoto.CreatorId, bulkPhoto.BulkPhotoStatus, bulkPhoto.CreatedAt, bulkPhoto.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return bulkPhoto, nil
}

// TODO mencoba kalau creator_id engga include di query apakah aman untuk non nill field
func (r *bulkPhotoRepository) FindDetailById(ctx context.Context, tx Querier, bulkPhotoID, creatorID string) (*[]*entity.BulkPhotoDetail, error) {
	bulkPhotoDetails := make([]*entity.BulkPhotoDetail, 0)
	query := `
	SELECT 
	bp.id AS bulk_photo_id,
	bp.creator_id AS bulk_photo_creator_id,
	bp.bulk_photo_status,
	bp.created_at AS bulk_photo_created_at,
	bp.updated_at AS bulk_photo_updated_at,

	p.id AS photo_id,
	p.creator_id AS photo_creator_id,
	p.title AS photo_title,
	p.owned_by_user_id AS photo_owned_by_user_id,
	p.compressed_url AS photo_compressed_url,
	p.is_this_you_url AS photo_is_this_you_url,
	p.your_moments_url AS photo_your_moments_url,
	p.collection_url AS photo_collection_url,
	p.price AS photo_price,
	p.price_str AS photo_price_str,
	p.latitude AS photo_latitude,
	p.longitude AS photo_longitude,
	p.description AS photo_description,
	p.original_at AS photo_original_at,
	p.created_at AS photo_created_at,
	p.updated_at AS photo_updated_at

	FROM bulk_photos AS bp
	LEFT JOIN
	photos AS p
	ON
	bp.id = p.bulk_photo_id
	WHERE bp.id = $1 
	AND
	bp.creator_id = $2
	`
	if err := tx.SelectContext(ctx, &bulkPhotoDetails, query, bulkPhotoID, creatorID); err != nil {
		return nil, err
	}
	return &bulkPhotoDetails, nil
}

func (r *bulkPhotoRepository) Update(ctx context.Context, tx Querier, bulkPhoto *entity.BulkPhoto) (*entity.BulkPhoto, error) {
	query := "UPDATE bulk_photos SET bulk_photo_status = $1, updated_at = $2 WHERE id = $3 AND creator_id = $4"
	_, err := tx.ExecContext(ctx, query, bulkPhoto.BulkPhotoStatus, bulkPhoto.UpdatedAt, bulkPhoto.Id, bulkPhoto.CreatorId)
	if err != nil {
		return nil, err
	}

	return bulkPhoto, err
}
