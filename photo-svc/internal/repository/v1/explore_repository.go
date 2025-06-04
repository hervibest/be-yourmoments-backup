package v1

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"strconv"

// 	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
// 	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
// 	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
// 	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"

// 	"github.com/jmoiron/sqlx"
// )

// type ExploreRepository interface {
// 	FindAllExploreSimilar(ctx context.Context, tx Querier, page, size int, similarity uint32, userId, creatorId string, isWishlist, isFavorite,
// 		isCart bool) ([]*entity.Explore, *model.PageMetadata, error)
// 	// FindAllUserCart(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
// 	// 	FindAllUserFavorite(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
// 	// FindAllUserWishlist(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
// 	// UserAddCart(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
// 	// UserAddFavorite(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
// 	// UserAddWishlist(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
// 	// UserDeleteCart(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
// 	// UserDeleteFavorite(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
// 	// UserDeleteWishlist(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
// 	UserAddStage(ctx context.Context, tx Querier, photoId, userId string, stage enum.PhotoStageEnum) error
// 	UserDeleteStage(ctx context.Context, tx Querier, photoId, userId string, stage enum.PhotoStageEnum) error
// }

// type exploreRepository struct {
// }

// func NewExploreRepository(db *sqlx.DB) (ExploreRepository, error) {
// 	return &exploreRepository{}, nil
// }

// // ISSUE #5 : Creators cannot explore their own  photos
// // User who owned picture can see but the other no
// // You have to make sure in AI logic only past different user_id (not same)
// func (r *exploreRepository) FindAllExploreSimilar(ctx context.Context, tx Querier, page, size int,
// 	similarity uint32, userId, creatorId string, isWishlist, isFavorite,
// 	isCart bool) ([]*entity.Explore, *model.PageMetadata, error) {
// 	results := make([]*entity.Explore, 0)

// 	var totalItems int
// 	countQuery := `
// 	SELECT
// 		COUNT(*)
// 	FROM
// 		user_similar_photos
// 	AS usp
// 	JOIN
// 		photos
// 	AS p on p.id = usp.photo_id
// 	WHERE
// 		usp.user_id = $1
// 	AND
// 		usp.similarity >= $2
// 	AND
// 		p.creator_id != $3
// 	AND
// 		(p.status = 'AVAILABLE' OR p.status= 'SOLD')
// 	AND
// 		(p.owned_by_user_id IS NULL OR p.owned_by_user_id = $1)
// 	 `

// 	var countArgs []interface{}

// 	query := `
// 	SELECT
// 		usp.photo_id,
// 		usp.user_id,
// 		usp.similarity,
// 		usp.is_wishlist,
// 		usp.is_resend,
// 		usp.is_cart,
// 		usp.is_favorite,

// 		p.creator_id,
// 		p.title,
// 		p.is_this_you_url,
// 		p.price,
// 		p.price_str,
// 		p.original_at,
// 		p.created_at,
// 		p.updated_at,

// 		cd.name,
// 		cd.min_quantity,
// 		cd.discount_type,
// 		cd.value,
// 		cd.is_active,

// 		pd.file_name,
// 		pd.file_key,
// 		pd.your_moments_type AS photo_detail_type

// 	FROM
// 		user_similar_photos AS usp

// 	JOIN
// 		photos AS p ON p.id = usp.photo_id

//     LEFT JOIN LATERAL (
//       SELECT *
//       FROM creator_discounts cd
//       WHERE cd.creator_id = p.creator_id AND cd.is_active = TRUE
//       ORDER BY cd.min_quantity ASC, cd.created_at DESC
//       LIMIT 1
//     ) cd ON TRUE

// 	LEFT JOIN LATERAL (
// 			SELECT
// 				pd.file_name,
// 				pd.file_key,
// 				pd.your_moments_type
// 			FROM photo_details pd
// 			WHERE pd.photo_id = p.id
// 			AND pd.your_moments_type =
// 				CASE
// 					WHEN p.owned_by_user_id IS NULL THEN 'YOU'::your_moments_type
// 					ELSE 'COLLECTION'::your_moments_type
// 				END
// 			LIMIT 1
// 		) pd ON TRUE

// 	-- LEFT JOIN photo_details AS pd
// 	--	ON p.id = pd.photo_id AND pd.your_moments_type = 'YOU'

// 	WHERE usp.user_id = $1
// 	AND usp.similarity >= $2
// 	AND p.creator_id != $3
// 	AND (p.status = 'AVAILABLE' OR p.status= 'SOLD')
// 	AND (p.owned_by_user_id IS NULL OR p.owned_by_user_id = $1)
// 	`

// 	if isWishlist {
// 		countQuery += " AND usp.is_wishlist = TRUE "
// 		query += " AND usp.is_wishlist = TRUE "
// 	}

// 	if isFavorite {
// 		countQuery += " AND usp.is_favorite = TRUE "
// 		query += " AND usp.is_favorite = TRUE "
// 	}

// 	if isCart {
// 		countQuery += " AND usp.is_cart = TRUE "
// 		query += " AND usp.is_cart = TRUE "
// 	}

// 	var queryArgs []interface{}

// 	// var conditions []string
// 	// var args []interface{}
// 	argIndex := 4

// 	// if username != "" {
// 	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
// 	// 	args = append(args, "%"+username+"%")
// 	// 	argIndex++
// 	// }

// 	// if len(conditions) > 0 {
// 	// 	query += " WHERE " + strings.Join(conditions, " AND ")
// 	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
// 	// }

// 	countArgs = append(countArgs, userId, similarity, creatorId)

// 	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
// 		log.Printf("Get context error in explore %v", err)
// 		return nil, nil, err
// 	}

// 	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

// 	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
// 	queryArgs = append(queryArgs, userId, similarity, creatorId, pageMetadata.Size, pageMetadata.Offset)

// 	if err := tx.SelectContext(ctx, &results, query, queryArgs...); err != nil {
// 		return nil, nil, err
// 	}

// 	return results, pageMetadata, nil
// }

// // UNUSED
// // func (r *exploreRepository) FindAllUserWishlist(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {
// // 	results := make([]*entity.Explore, 0)

// // 	var totalItems int
// // 	countQuery := `SELECT COUNT(*) FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_wishlist = true`

// // 	var countArgs []interface{}

// // 	query := `
// // 	SELECT
// // 		usp.photo_id,
// // 		usp.user_id,
// // 		usp.similarity,
// // 		usp.is_wishlist,
// // 		usp.is_resend,
// // 		usp.is_cart,
// // 		usp.is_favorite,

// // 		p.creator_id,
// // 		p.title,
// // 		p.is_this_you_url,
// // 		p.your_moments_url,
// // 		p.price,
// // 		p.price_str,
// // 		p.original_at,
// // 		p.created_at,
// // 		p.updated_at
// // 	FROM
// // 		user_similar_photos AS usp
// // 	JOIN
// // 		photos AS p ON p.id = usp.photo_id
// // 	WHERE
// // 		usp.user_id = $1
// // 		AND usp.is_wishlist = true
// // 		-- TODO: Add logic for additional filtering (e.g., photo availability, ownership)
// // 		`

// // 	var queryArgs []interface{}

// // 	// var conditions []string
// // 	// var args []interface{}
// // 	argIndex := 2

// // 	// if username != "" {
// // 	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
// // 	// 	args = append(args, "%"+username+"%")
// // 	// 	argIndex++
// // 	// }

// // 	// if len(conditions) > 0 {Fin
// // 	// 	query += " WHERE " + strings.Join(conditions, " AND ")
// // 	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
// // 	// }

// // 	// log.Print(countQuery)
// // 	// log.Print(countArgs...)

// // 	countArgs = append(countArgs, userId)

// // 	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
// // 		return nil, nil, err
// // 	}

// // 	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

// // 	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
// // 	queryArgs = append(queryArgs, userId, pageMetadata.Size, pageMetadata.Offset)
// // 	// log.Println(query)
// // 	// log.Println(queryArgs...)

// // 	if err := tx.SelectContext(ctx, &results, query, queryArgs...); err != nil {
// // 		return nil, nil, err
// // 	}

// // 	return results, pageMetadata, nil
// // }

// func (r *exploreRepository) UserAddWishlist(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
// 	query := "UPDATE user_similar_photos SET is_wishlist = true WHERE photo_id = $1 AND user_id = $2"
// 	_, err := tx.ExecContext(ctx, query, photoId, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *exploreRepository) UserDeleteWishlist(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
// 	query := "UPDATE user_similar_photos SET is_wishlist = false WHERE photo_id = $1 AND user_id = $2"
// 	_, err := tx.ExecContext(ctx, query, photoId, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // func (r *exploreRepository) FindAllUserFavorite(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {

// // 	results := make([]*entity.Explore, 0)

// // 	var totalItems int
// // 	countQuery := `SELECT COUNT(*) FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_favorite = true`

// // 	var countArgs []interface{}

// // 	query := `SELECT usp.photo_id, usp.user_id, usp.similarity, usp.is_wishlist,
// // 	usp.is_resend, usp.is_cart, usp.is_favorite, p.creator_id, p.title, p.is_this_you_url,
// // 	p.your_moments_url, p.price, p.price_str, p. original_at, p.created_at, p.updated_at
// // 	FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_favorite = true`
// // 	//TO DO ADD logic for
// // 	var queryArgs []interface{}

// // 	// var conditions []string
// // 	// var args []interface{}
// // 	argIndex := 2

// // 	// if username != "" {
// // 	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
// // 	// 	args = append(args, "%"+username+"%")
// // 	// 	argIndex++
// // 	// }

// // 	// if len(conditions) > 0 {Fin
// // 	// 	query += " WHERE " + strings.Join(conditions, " AND ")
// // 	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
// // 	// }

// // 	// log.Print(countQuery)
// // 	// log.Print(countArgs...)

// // 	countArgs = append(countArgs, userId)

// // 	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
// // 		return nil, nil, err
// // 	}

// // 	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

// // 	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
// // 	queryArgs = append(queryArgs, userId, pageMetadata.Size, pageMetadata.Offset)
// // 	// log.Println(query)
// // 	// log.Println(queryArgs...)

// // 	if err := tx.SelectContext(ctx, &results, query, queryArgs...); err != nil {
// // 		return nil, nil, err
// // 	}

// // 	return results, pageMetadata, nil
// // }

// func (r *exploreRepository) UserAddFavorite(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
// 	query := "UPDATE user_similar_photos SET is_favorite = true WHERE photo_id = $1 AND user_id = $2"
// 	_, err := tx.ExecContext(ctx, query, photoId, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *exploreRepository) UserDeleteFavorite(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
// 	query := "UPDATE user_similar_photos SET is_favorite = false WHERE photo_id = $1 AND user_id = $2"
// 	_, err := tx.ExecContext(ctx, query, photoId, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // func (r *exploreRepository) FindAllUserCart(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {

// // 	results := make([]*entity.Explore, 0)

// // 	var totalItems int
// // 	countQuery := `SELECT COUNT(*) FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_cart = true`

// // 	var countArgs []interface{}

// // 	query := `SELECT usp.photo_id, usp.user_id, usp.similarity, usp.is_wishlist,
// // 	usp.is_resend, usp.is_cart, usp.is_favorite, p.creator_id, p.title, p.is_this_you_url,
// // 	p.your_moments_url, p.price, p.price_str, p. original_at, p.created_at, p.updated_at
// // 	FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_cart = true`
// // 	//TO DO ADD logic for
// // 	var queryArgs []interface{}

// // 	// var conditions []string
// // 	// var args []interface{}
// // 	argIndex := 2

// // 	// if username != "" {
// // 	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
// // 	// 	args = append(args, "%"+username+"%")
// // 	// 	argIndex++
// // 	// }

// // 	// if len(conditions) > 0 {Fin
// // 	// 	query += " WHERE " + strings.Join(conditions, " AND ")
// // 	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
// // 	// }

// // 	// log.Print(countQuery)
// // 	// log.Print(countArgs...)

// // 	countArgs = append(countArgs, userId)

// // 	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
// // 		return nil, nil, err
// // 	}

// // 	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

// // 	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
// // 	queryArgs = append(queryArgs, userId, pageMetadata.Size, pageMetadata.Offset)
// // 	// log.Println(query)
// // 	// log.Println(queryArgs...)

// // 	if err := tx.SelectContext(ctx, &results, query, queryArgs...); err != nil {
// // 		return nil, nil, err
// // 	}
// // 	return results, pageMetadata, nil
// // }

// // ISSUE #1 : USER CAN BYPASS THIS UPDATE U SHOULD HANDLE EVERY PHOTOS LOGIC
// func (r *exploreRepository) UserAddStage(ctx context.Context, tx Querier, photoId, userId string,
// 	stage enum.PhotoStageEnum) error {
// 	var queryStage string
// 	var stageCondition string
// 	switch stage {
// 	case enum.PhotoStageWishlist:
// 		queryStage = "is_wishlist = true"
// 		stageCondition = "is_wishlist = false"
// 	case enum.PhotoStageFavorite:
// 		queryStage = "is_favorite = true"
// 		stageCondition = "is_favorite = false"
// 	case enum.PhotoStageCart:
// 		queryStage = "is_cart = true"
// 		stageCondition = "is_cart = false"
// 	default:
// 		return errors.New("invalid photo stage")
// 	}

// 	log.Printf("ini adalah photo id : %s dan ini adalah user id : %s", photoId, userId)
// 	query := fmt.Sprintf("UPDATE user_similar_photos SET %s WHERE photo_id = $1 AND user_id = $2 AND %s", queryStage, stageCondition)
// 	log.Print(query)
// 	res, err := tx.ExecContext(ctx, query, photoId, userId)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("no rows updated, possibly already unset or not found")
// 	}

// 	return nil
// }

// func (r *exploreRepository) UserDeleteStage(ctx context.Context, tx Querier, photoId, userId string,
// 	stage enum.PhotoStageEnum) error {
// 	var queryStage string
// 	var stageCondition string
// 	switch stage {
// 	case enum.PhotoStageWishlist:
// 		queryStage = "is_wishlist = false"
// 		stageCondition = "is_wishlist = true"
// 	case enum.PhotoStageFavorite:
// 		queryStage = "is_favorite = false"
// 		stageCondition = "is_favorite = true"
// 	case enum.PhotoStageCart:
// 		queryStage = "is_cart = false"
// 		stageCondition = "is_cart = true"
// 	default:
// 		return errors.New("invalid photo stage")
// 	}

// 	query := fmt.Sprintf("UPDATE user_similar_photos SET %s WHERE photo_id = $1 AND user_id = $2 AND %s", queryStage, stageCondition)
// 	res, err := tx.ExecContext(ctx, query, photoId, userId)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("no rows updated, possibly already unset or not found")
// 	}

// 	return nil
// }
