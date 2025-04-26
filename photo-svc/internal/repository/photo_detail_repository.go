package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"context"
	"fmt"
	"log"
)

type PhotoDetailRepository interface {
	Create(tx Querier, photoDetail *entity.PhotoDetail) (*entity.PhotoDetail, error)
	BulkCreate(ctx context.Context, tx Querier, items []*entity.PhotoDetail) (*[]*entity.PhotoDetail, error)
}

type photoDetailRepository struct {
}

func NewPhotoDetailRepository() PhotoDetailRepository {
	return &photoDetailRepository{}
}

func (r *photoDetailRepository) Create(tx Querier, photoDetail *entity.PhotoDetail) (*entity.PhotoDetail, error) {
	log.Println("Created accessed")
	query := `INSERT INTO photo_details 
			  (id, photo_id, file_name, file_key, size, type, checksum, width, height, url, your_moments_type, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := tx.Exec(query, photoDetail.Id, photoDetail.PhotoId, photoDetail.FileName, photoDetail.FileKey, photoDetail.Size, photoDetail.Type,
		photoDetail.Checksum, photoDetail.Width, photoDetail.Height, photoDetail.Url, photoDetail.YourMomentsType, photoDetail.CreatedAt, photoDetail.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert photoDetail: %w", err)
	}

	return photoDetail, nil
}

func (r *photoDetailRepository) BulkCreate(ctx context.Context, tx Querier, items []*entity.PhotoDetail) (*[]*entity.PhotoDetail, error) {
	query := `INSERT INTO photo_details (id, photo_id, file_name, file_key, size, type, checksum, width, height, url, your_moments_type, created_at, updated_at)
	          VALUES (:id, :photo_id, :file_name, :file_key, :size, :type, :checksum, :width, :height, :url, :your_moments_type, :created_at, :updated_at)`

	_, err := tx.NamedExecContext(ctx, query, items)
	if err != nil {
		log.Printf("error inserting bulk photo details: %v", err)
		return nil, err
	}

	return &items, nil
}
