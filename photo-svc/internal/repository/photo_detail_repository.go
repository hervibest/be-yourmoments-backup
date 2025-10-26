package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
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

	log.Default().Printf("Ini adalah photo detail, id: %s, photo_id: %s, file_name: %s, file_key: %s, size: %d, type: %s, checksum: %s, width: %d, height: %d, url: %s, your_moments_type: %s, created_at: %s, updated_at: %s",
		photoDetail.Id, photoDetail.PhotoId, photoDetail.FileName, photoDetail.FileKey, photoDetail.Size, photoDetail.Type)
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

// TODO ISSUE -- bulk create
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

func (r *photoDetailRepository) CreateBulk(tx Querier, photoDetails []*entity.PhotoDetail) error {
	if len(photoDetails) == 0 {
		return nil
	}

	insertValues := make([]string, 0, len(photoDetails))
	insertArgs := make([]interface{}, 0, len(photoDetails)*10) // 10 kolom per photo detail
	counter := 1

	for _, detail := range photoDetails {
		insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			counter, counter+1, counter+2, counter+3, counter+4, counter+5, counter+6, counter+7, counter+8, counter+9))
		insertArgs = append(insertArgs,
			detail.Id,
			detail.PhotoId,
			detail.FileName,
			detail.FileKey,
			detail.Size,
			detail.Type,
			detail.Checksum,
			detail.Height,
			detail.Width,
			detail.Url,
		)
		counter += 10
	}

	query := `
		INSERT INTO photo_details 
		(id, photo_id, file_name, file_key, size, type, checksum, height, width, url) 
		VALUES ` + strings.Join(insertValues, ", ")

	_, err := tx.Exec(query, insertArgs...)
	if err != nil {
		log.Println("Error at CreateBulkPhotoDetail:", err)
		return err
	}

	return nil
}
