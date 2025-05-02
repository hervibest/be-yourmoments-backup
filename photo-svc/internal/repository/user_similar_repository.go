package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"

	"github.com/lib/pq"
)

type UserSimilarRepository interface {
	InsertOrUpdateByPhotoId(tx Querier, photoId string, userSimilarPhotos *[]*entity.UserSimilarPhoto) error
	InserOrUpdateByUserId(tx Querier, userId string, userSimilarPhotos *[]*entity.UserSimilarPhoto) error
	InsertOrUpdateBulk(ctx context.Context, tx Querier, photoUserSimilarMap map[string][]*entity.UserSimilarPhoto) error // UpdateUsersForPhoto(ctx context.Context, db Querier, photoId string, userIds []string) error
	// GetSimilarPhotosByUser(ctx context.Context, db Querier, userId string) (*UserSimilarPhotosResponse, error)
	// DeleteSimilarUsers(ctx context.Context, db Querier, photoId string) error
}

type userSimilarRepository struct {
}

func NewUserSimilarRepository() UserSimilarRepository {
	return &userSimilarRepository{}
}

func (r *userSimilarRepository) InsertOrUpdateByPhotoId(tx Querier, photoId string, userSimilarPhotos *[]*entity.UserSimilarPhoto) error {
	now := time.Now()
	if len(*userSimilarPhotos) == 0 {
		log.Println("ALL DELETED BECAUSE OF ZERO")
		if _, err := tx.Exec("DELETE FROM user_similar_photos WHERE photo_id = $1", photoId); err != nil {
			return err
		}
		return nil
	}

	placeholders := make([]string, len(*userSimilarPhotos))
	deleteArgs := make([]interface{}, 0, len(*userSimilarPhotos)+1)
	deleteArgs = append(deleteArgs, photoId)
	for i, userSimilarPhoto := range *userSimilarPhotos {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		deleteArgs = append(deleteArgs, userSimilarPhoto.UserId)
	}

	deleteQuery := "DELETE FROM user_similar_photos WHERE photo_id = $1 AND user_id NOT IN (" + strings.Join(placeholders, ", ") + ")"
	if _, err := tx.Exec(deleteQuery, deleteArgs...); err != nil {
		log.Println("Error at delete query:", err)
		return err
	}

	insertValues := make([]string, 0, len(*userSimilarPhotos))
	insertArgs := make([]interface{}, 0, len(*userSimilarPhotos)*4)
	placeholderCounter := 1
	for _, userSimilarPhoto := range *userSimilarPhotos {
		// Misalnya, baris pertama: ($1, $2, $3, $4), baris kedua: ($5, $6, $7, $8), dst.
		insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3, placeholderCounter+4))
		insertArgs = append(insertArgs, photoId, userSimilarPhoto.UserId, userSimilarPhoto.Similarity, now, now)
		placeholderCounter += 5
	}

	insertQuery := "INSERT INTO user_similar_photos (photo_id, user_id, similarity, created_at, updated_at) VALUES " +
		strings.Join(insertValues, ", ") +
		" ON CONFLICT (user_id, photo_id) DO UPDATE SET updated_at = EXCLUDED.updated_at"

	if _, err := tx.Exec(insertQuery, insertArgs...); err != nil {
		log.Println("Error at insert query:", err)
		return err
	}

	return nil
}

func (r *userSimilarRepository) InserOrUpdateByUserId(tx Querier, userId string, userSimilarPhotos *[]*entity.UserSimilarPhoto) error {
	now := time.Now()
	if len(*userSimilarPhotos) == 0 {
		if _, err := tx.Exec("DELETE FROM user_similar_photos WHERE user_id = $1", userId); err != nil {
			return err
		}
		return nil
	}

	placeholders := make([]string, len(*userSimilarPhotos))
	deleteArgs := make([]interface{}, 0, len(*userSimilarPhotos)+1)
	deleteArgs = append(deleteArgs, userId)
	for i, userSimilarPhoto := range *userSimilarPhotos {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		deleteArgs = append(deleteArgs, userSimilarPhoto.UserId)
	}

	deleteQuery := "DELETE FROM user_similar_photos WHERE user_id = $1 AND photo_id NOT IN (" + strings.Join(placeholders, ", ") + ")"
	if _, err := tx.Exec(deleteQuery, deleteArgs...); err != nil {
		log.Println("Error at delete query:", err)
		return err
	}

	insertValues := make([]string, 0, len(*userSimilarPhotos))
	insertArgs := make([]interface{}, 0, len(*userSimilarPhotos)*4)
	placeholderCounter := 1
	for _, userSimilarPhoto := range *userSimilarPhotos {
		// Misalnya, baris pertama: ($1, $2, $3, $4), baris kedua: ($5, $6, $7, $8), dst.
		insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3, placeholderCounter+4))
		insertArgs = append(insertArgs, userId, userSimilarPhoto.PhotoId, userSimilarPhoto.Similarity, now, now)
		placeholderCounter += 5
	}

	insertQuery := "INSERT INTO user_similar_photos (user_id, photo_id, similarity, created_at, updated_at) VALUES " +
		strings.Join(insertValues, ", ") +
		" ON CONFLICT (photo_id, user_id) DO UPDATE SET updated_at = EXCLUDED.updated_at"

	if _, err := tx.Exec(insertQuery, insertArgs...); err != nil {
		log.Println("Error at insert query:", err)
		return err
	}

	return nil
}

func (r *userSimilarRepository) InsertOrUpdateBulk(ctx context.Context, tx Querier, photoUserSimilarMap map[string][]*entity.UserSimilarPhoto) error {
	now := time.Now()

	// Kumpulkan semua photo_id dan (photo_id, user_id) pair
	photoIDs := make([]string, 0, len(photoUserSimilarMap))
	userPairs := make([]string, 0)
	deleteArgs := make([]interface{}, 0)
	argCounter := 1

	for photoID, userSimilars := range photoUserSimilarMap {
		photoIDs = append(photoIDs, photoID)
		for _, userSimilar := range userSimilars {
			userPairs = append(userPairs, fmt.Sprintf("($%d, $%d)", argCounter, argCounter+1))
			deleteArgs = append(deleteArgs, photoID, userSimilar.UserId)
			argCounter += 2
		}
	}

	if len(photoIDs) == 0 {
		return nil
	}

	if len(userPairs) == 0 {
		return nil // Tidak ada data untuk dihapus atau diperbarui
	}

	// Hapus user_similar yang tidak ada dalam batch baru
	deleteQuery := `
		DELETE FROM user_similar_photos
		WHERE (photo_id, user_id) NOT IN (` + strings.Join(userPairs, ", ") + `)
		AND photo_id = ANY($` + fmt.Sprint(argCounter) + `)
	`

	deleteArgs = append(deleteArgs, pq.Array(photoIDs)) // postgres array

	if _, err := tx.ExecContext(ctx, deleteQuery, deleteArgs...); err != nil {
		log.Println("Error during bulk delete:", err)
		return err
	}

	// Masukkan semua user_similar baru
	insertValues := make([]string, 0)
	insertArgs := make([]interface{}, 0)
	counter := 1

	for photoID, userSimilars := range photoUserSimilarMap {
		for _, userSimilar := range userSimilars {
			insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", counter, counter+1, counter+2, counter+3, counter+4))
			insertArgs = append(insertArgs,
				photoID,
				userSimilar.UserId,
				userSimilar.Similarity,
				now,
				now,
			)
			counter += 5 // ✅ benar
		}
	}

	if len(insertValues) > 0 {
		insertQuery := `
			INSERT INTO user_similar_photos 
			(photo_id, user_id, similarity, created_at, updated_at)
			VALUES ` + strings.Join(insertValues, ", ") + `
			ON CONFLICT (user_id, photo_id) 
			DO UPDATE SET
				similarity = EXCLUDED.similarity,
				updated_at = NOW()
		`

		if _, err := tx.ExecContext(ctx, insertQuery, insertArgs...); err != nil {
			log.Println("Error during bulk insert:", err)
			return err
		}
	}

	return nil
}
