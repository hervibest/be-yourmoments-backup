package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type photoPreparedStmt struct {
	findByPhotoId *sqlx.Stmt
}

func newPhotoPreparedStmt(db *sqlx.DB) (*photoPreparedStmt, error) {
	findByPhotoIdStmt, err := db.Preparex("SELECT * FROM photos WHERE id = $1")
	if err != nil {
		return nil, err
	}

	return &photoPreparedStmt{
		findByPhotoId: findByPhotoIdStmt,
	}, nil
}

type PhotoRepository interface {
	Create(tx Querier, photo *entity.Photo) (*entity.Photo, error)
	FindByPhotoId(ctx context.Context, tx Querier, photoId string) (*entity.Photo, error)
	FindBuyableByPhotoId(ctx context.Context, tx Querier, photoId string, forUpdate bool) (*entity.Photo, error)
	UpdateProcessedUrl(tx Querier, photo *entity.Photo) error
	UpdateCompressedUrl(tx Querier, photo *entity.Photo) error
	GetSimilarPhotosByIDs(ctx context.Context, tx Querier, userId, creatorId string, ids []string, forUpdate bool, converter func(string) string) (*[]*entity.Photo, error)
	GetManyInTransactionByIDsAndUserID(ctx context.Context, tx Querier, userId string, ids []string, forUpdate bool) (*[]*entity.Photo, error)
	UpdatePhotoOwnerAndStatusByIds(ctx context.Context, tx Querier, ownerID string, photoIDs []string) error
	BulkCreate(ctx context.Context, tx Querier, items []*entity.Photo) (*[]*entity.Photo, error) // UpdatePhotoStatus(ctx context.Context, db Querier, photo *entity.Photo) error
	UpdateProcessedUrlBulk(tx Querier, photos []*entity.Photo) error
	BulkIncrementTotal(ctx context.Context, tx Querier, photoIDs []string) error
	BulkAddPhotoTotals(ctx context.Context, tx Querier, photoCountMap map[string]int32) error
	AddPhotoTotal(ctx context.Context, tx Querier, photoID string, count int) error
	UserGetPhotoWithDetail(ctx context.Context, tx Querier, photoIDs []string, userID string) ([]*entity.PhotoWithDetail, error)
	UpdatePhotoStatusesByIDs(ctx context.Context, tx Querier, status enum.PhotoStatusEnum, ids []string) error
}

type photoRepository struct {
	photoPreparedStmt *photoPreparedStmt
}

func NewPhotoRepository(db *sqlx.DB) (PhotoRepository, error) {
	photoPreparedStmt, err := newPhotoPreparedStmt(db)
	if err != nil {
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

func (r *photoRepository) FindBuyableByPhotoId(ctx context.Context, tx Querier, photoId string, forUpdate bool) (*entity.Photo, error) {
	photo := new(entity.Photo)
	query := "SELECT * FROM photos WHERE id = $1 AND owned_by_user_id IS NULL AND status = 'AVAILABLE'"
	if forUpdate {
		query += " FOR UPDATE"
	}

	if err := tx.GetContext(ctx, photo, query, photoId); err != nil {
		return nil, err
	}

	return photo, nil
}

func (r *photoRepository) GetManyInTransactionByIDsAndUserID(ctx context.Context, tx Querier, userId string, ids []string, forUpdate bool) (*[]*entity.Photo, error) {
	photos := make([]*entity.Photo, 0)
	query := `
		SELECT 
			id,
			creator_id,
			title,
			is_this_you_url,
			your_moments_url,
			price
		FROM 
			photos AS p
		JOIN 
			user_similar_photos AS up
		ON 
			up.photo_id = p.id
		WHERE 
			id = ANY($1) 
		AND 
			up.user_id = $2 
		AND 
			p.owned_by_user_id IS NULL
		AND
			p.status = 'IN_TRANSACTION'
			`
	if forUpdate {
		query += ` FOR UPDATE`
	}

	if err := tx.SelectContext(ctx, &photos, query, ids, userId); err != nil {
		return nil, err
	}
	return &photos, nil
}

// TODO HANDLE BEST PRACTIE URL, CDN AND LOCKING FOR UPDATE
func (r *photoRepository) GetSimilarPhotosByIDs(ctx context.Context, tx Querier, userId, creatorId string, ids []string, forUpdate bool, cdn func(string) string) (*[]*entity.Photo, error) {
	photos := make([]*entity.Photo, 0)
	query := `
		SELECT 
			id,
			creator_id,
			title,
			is_this_you_url,
			your_moments_url,
			price,

			pd.file_name,
			pd.file_key,
			pd.your_moments_type
		FROM 
			photos AS p
		JOIN 
			user_similar_photos AS up
		ON 
			up.photo_id = p.id
		LEFT JOIN LATERAL (
			SELECT 
				pd.file_name,
				pd.file_key,
				pd.your_moments_type
			FROM photo_details pd
			WHERE pd.photo_id = p.id
			AND pd.your_moments_type = 
				CASE 
					WHEN p.owned_by_user_id IS NULL THEN 'YOU'::your_moments_type
					ELSE 'COLLECTION'::your_moments_type
				END
			LIMIT 1
		) pd ON TRUE
		WHERE 
			id = ANY($1) 
		AND 
			up.user_id = $2 
		AND 
			p.owned_by_user_id IS NULL
		AND 
			p.creator_id != $3
		AND
			p.status = 'AVAILABLE'
			`
	// if forUpdate {
	// 	query += ` FOR UPDATE`
	// }

	if err := tx.SelectContext(ctx, &photos, query, ids, userId, creatorId); err != nil {
		return nil, err
	}
	for i := range photos {
		log.Println("File Key:", photos[i].FileKey.String)
		log.Println("Photo Type:", photos[i].PhotoType.String)
		urlCdn := converter.ToCollectionOrIsYouURL(photos[i].PhotoType.String, photos[i].FileKey.String, cdn)
		photos[i].YourMomentsUrl = sql.NullString{String: urlCdn.IsThisYouURL, Valid: true}
		log.Println("File Key:", photos[i])

	}
	return &photos, nil
}

func (r *photoRepository) UpdatePhotoOwnerAndStatusByIds(ctx context.Context, tx Querier, ownerID string, photoIDs []string) error {
	query := "UPDATE photos SET owned_by_user_id = $1, status = $2 WHERE id = ANY($3)"
	_, err := tx.ExecContext(ctx, query, ownerID, enum.PhotoStatusSoldEnum, photoIDs)
	if err != nil {
		return err
	}

	return nil
}

// TODO NamedExecContext doesnt use bulk upload
func (r *photoRepository) BulkCreate(ctx context.Context, tx Querier, items []*entity.Photo) (*[]*entity.Photo, error) {
	query := `INSERT INTO photos (id, creator_id, bulk_photo_id, title, collection_url, price, price_str, latitude, longitude, description, original_at, created_at, updated_at)
	          VALUES (:id, :creator_id, :bulk_photo_id, :title, :collection_url, :price, :price_str, :latitude, :longitude, :description, :original_at, :created_at, :updated_at)`

	_, err := tx.NamedExecContext(ctx, query, items)
	if err != nil {
		log.Printf("error inserting bulk photo: %v", err)
		return nil, err
	}

	return &items, nil
}

func (r *photoRepository) UpdateProcessedUrlBulk(tx Querier, photos []*entity.Photo) error {
	if len(photos) == 0 {
		return nil
	}

	caseIsThisYouURL := "CASE id"
	caseYourMomentsURL := "CASE id"
	caseUpdatedAt := "CASE id"

	args := make([]interface{}, 0, len(photos)*4) // *4 karena 3 CASE + 1 WHERE IN
	argPos := 1                                   // PostgreSQL bind starts with $1
	idArgs := make([]string, 0, len(photos))

	for _, photo := range photos {
		caseIsThisYouURL += fmt.Sprintf(" WHEN $%d THEN $%d", argPos, argPos+1)
		args = append(args, photo.Id, photo.IsThisYouURL.String)
		argPos += 2

		caseYourMomentsURL += fmt.Sprintf(" WHEN $%d THEN $%d", argPos, argPos+1)
		args = append(args, photo.Id, photo.YourMomentsUrl.String)
		argPos += 2

		caseUpdatedAt += fmt.Sprintf(" WHEN $%d THEN $%d::timestamptz", argPos, argPos+1)
		args = append(args, photo.Id, photo.UpdatedAt)
		argPos += 2

		idArgs = append(idArgs, photo.Id)
	}

	caseIsThisYouURL += " END"
	caseYourMomentsURL += " END"
	caseUpdatedAt += " END"

	// Untuk WHERE IN kita pakai Array
	args = append(args, pq.Array(idArgs))

	query := fmt.Sprintf(`
		UPDATE photos
		SET
			is_this_you_url = %s,
			your_moments_url = %s,
			updated_at = %s
		WHERE id = ANY($%d)
	`, caseIsThisYouURL, caseYourMomentsURL, caseUpdatedAt, argPos)

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to bulk update photo: %w", err)
	}
	return nil
}

func (r *photoRepository) BulkIncrementTotal(ctx context.Context, tx Querier, photoIDs []string) error {
	if len(photoIDs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(photoIDs))
	args := make([]interface{}, 0, len(photoIDs))

	for i, id := range photoIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, 1)", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		UPDATE photos AS p
		SET total_user_similar = p.total_user_similar + v.increment
		FROM (
			VALUES %s
		) AS v(photo_id, increment)
		WHERE p.id = v.photo_id;
	`, strings.Join(valueStrings, ", "))

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *photoRepository) BulkAddPhotoTotals(ctx context.Context, tx Querier, photoCountMap map[string]int32) error {
	if len(photoCountMap) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(photoCountMap))
	args := make([]interface{}, 0, len(photoCountMap)*2)
	i := 1

	for photoID, count := range photoCountMap {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::text, $%d::int)", i, i+1))
		args = append(args, photoID, count)
		i += 2
	}

	query := fmt.Sprintf(`
		UPDATE photos AS p
		SET total_user_similar = p.total_user_similar + v.increment
		FROM (
			VALUES %s
		) AS v(photo_id, increment)
		WHERE p.id = v.photo_id;
	`, strings.Join(valueStrings, ", "))

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *photoRepository) AddPhotoTotal(ctx context.Context, tx Querier, photoID string, count int) error {
	query := `
		UPDATE photos
		SET total_user_similar = total_user_similar + $1
		WHERE id = $2;
	`
	_, err := tx.ExecContext(ctx, query, count, photoID)
	if err != nil {
		return err
	}
	return nil
}

// ISSUE  handle autorization for this query
func (r *photoRepository) UserGetPhotoWithDetail(ctx context.Context, tx Querier, photoIDs []string, userID string) ([]*entity.PhotoWithDetail, error) {
	photoWithDetails := make([]*entity.PhotoWithDetail, 0)
	query := ` 
				SELECT 
			p.id AS photo_id,
			p.creator_id,
			p.title,
			p.price,
			p.price_str,
			p.latitude,
			p.longitude,
			p.description,
			p.original_at,
			p.created_at,
			p.updated_at,

			pd.file_name,
			pd.file_key,
			pd.size,
			pd.type,
			pd.width,
			pd.height,
			pd.your_moments_type
		FROM photos p
		LEFT JOIN LATERAL (
			SELECT 
				pd.file_name,
				pd.file_key,
				pd.size,
				pd.type,
				pd.width,
				pd.height,
				pd.your_moments_type
			FROM photo_details pd
			WHERE pd.photo_id = p.id
			AND pd.your_moments_type = 
				CASE 
					WHEN p.owned_by_user_id IS NULL THEN 'ISYOU'::your_moments_type
					ELSE 'COLLECTION'::your_moments_type
				END
			LIMIT 1
		) pd ON TRUE
		WHERE p.id = ANY($1) 
		AND (
			p.owned_by_user_id IS NULL
			OR p.owned_by_user_id = $2
		)
		`
	if err := tx.SelectContext(ctx, &photoWithDetails, query, photoIDs, userID); err != nil {
		return nil, err
	}
	return photoWithDetails, nil
}

func (r *photoRepository) UpdatePhotoStatusesByIDs(ctx context.Context, tx Querier, status enum.PhotoStatusEnum, ids []string) error {
	query := `UPDATE photos SET status = ($1) WHERE id = ANY($2)`
	_, err := tx.ExecContext(ctx, query, status, ids)
	if err != nil {
		return err
	}
	return nil
}
