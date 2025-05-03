package repository

import (
	"context"
	"log"
	"strconv"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"

	"github.com/jmoiron/sqlx"
)

type explorePreparedStmt struct {
	findByUserId *sqlx.Stmt
}

func newExplorePreparedStmt(db *sqlx.DB) (*explorePreparedStmt, error) {
	findByUserIdStmt, err := db.Preparex(`SELECT usp.photo_id, usp.user_id, usp.similarity, usp.is_wishlist,
	usp.is_resend, usp.is_cart, usp.is_favorite, p.creator_id, p.title, p.is_this_you_url,
	p.your_moments_url, p.price, p.price_str, p. original_at, p.created_at, p.updated_at
	FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1`)
	if err != nil {
		return nil, err
	}

	return &explorePreparedStmt{
		findByUserId: findByUserIdStmt,
	}, nil
}

type ExploreRepository interface {
	FindAllExploreSimilar(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
	FindAllUserCart(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
	FindAllUserFavorite(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
	FindAllUserWishlist(ctx context.Context, tx Querier, page int, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error)
	FindByUserId(ctx context.Context, userId string) (*[]*entity.Explore, error)
	UserAddCart(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
	UserAddFavorite(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
	UserAddWishlist(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
	UserDeleteCart(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
	UserDeleteFavorite(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
	UserDeleteWishlist(ctx context.Context, tx Querier, similarity uint32, photoId string, userId string) error
}

type exploreRepository struct {
	explorePreparedStmt *explorePreparedStmt
}

func NewExploreRepository(db *sqlx.DB) (ExploreRepository, error) {
	explorePreparedStmt, err := newExplorePreparedStmt(db)
	if err != nil {
		return nil, err
	}

	return &exploreRepository{
		explorePreparedStmt: explorePreparedStmt,
	}, nil
}

func (r *exploreRepository) FindByUserId(ctx context.Context, userId string) (*[]*entity.Explore, error) {
	explores := make([]*entity.Explore, 0)

	rows, err := r.explorePreparedStmt.findByUserId.QueryxContext(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		explore := new(entity.Explore)
		err := rows.StructScan(explore)
		if err != nil {
			return nil, err
		}

		explores = append(explores, explore)
	}

	return &explores, nil
}

// ISSUE #5 : Creators cannot explore their own  photos
// You have to make sure in AI logic only past different user_id (not same)
func (r *exploreRepository) FindAllExploreSimilar(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {
	results := make([]*entity.Explore, 0)

	var totalItems int
	countQuery :=
		`SELECT 
	COUNT(*) 
	FROM user_similar_photos 
	AS usp 
	JOIN photos 
	AS p on p.id = usp.photo_id 
	WHERE usp.user_id = $1 
	AND usp.similarity = $2
	`

	var countArgs []interface{}

	query := `
		SELECT 
		usp.photo_id,
		usp.user_id,
		usp.similarity,
		usp.is_wishlist,
		usp.is_resend,
		usp.is_cart,
		usp.is_favorite,

		p.creator_id,
		p.title,
		p.is_this_you_url,
		p.price,
		p.price_str,
		p.original_at,
		p.created_at,
		p.updated_at,

		cd.name,
		cd.min_quantity,
		cd.discount_type,
		cd.value,
		cd.is_active,

		pd.file_name,
		pd.file_key,
		pd.your_moments_type AS photo_detail_type

	FROM user_similar_photos AS usp

	JOIN photos AS p ON p.id = usp.photo_id

	LEFT JOIN creator_discounts AS cd 
		ON p.creator_id = cd.creator_id AND cd.is_active = TRUE

	LEFT JOIN photo_details AS pd 
		ON p.id = pd.photo_id AND pd.your_moments_type = 'YOU'

	WHERE usp.user_id = $1
	AND usp.similarity >= $2
	`

	var queryArgs []interface{}

	// var conditions []string
	// var args []interface{}
	argIndex := 3

	// if username != "" {
	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
	// 	args = append(args, "%"+username+"%")
	// 	argIndex++
	// }

	// if len(conditions) > 0 {
	// 	query += " WHERE " + strings.Join(conditions, " AND ")
	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
	// }

	// log.Print(countQuery)
	// log.Print(countArgs...)

	countArgs = append(countArgs, userId, similarity)

	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
		log.Printf("Get context error in explore %v", err)
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	queryArgs = append(queryArgs, userId, similarity, pageMetadata.Size, pageMetadata.Offset)
	// log.Println(query)
	// log.Println(queryArgs...)

	rows, err := tx.QueryxContext(ctx, query, queryArgs...)
	if err != nil {
		log.Printf("Queryx context error in explore %v", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		result := new(entity.Explore)
		if err := rows.StructScan(result); err != nil {
			return nil, nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}

func (r *exploreRepository) FindAllUserWishlist(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {

	results := make([]*entity.Explore, 0)

	var totalItems int
	countQuery := `SELECT COUNT(*) FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_wishlist = true`

	var countArgs []interface{}

	query := `SELECT usp.photo_id, usp.user_id, usp.similarity, usp.is_wishlist,
	usp.is_resend, usp.is_cart, usp.is_favorite, p.creator_id, p.title, p.is_this_you_url,
	p.your_moments_url, p.price, p.price_str, p. original_at, p.created_at, p.updated_at
	FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_wishlist = true`
	//TO DO ADD logic for
	var queryArgs []interface{}

	// var conditions []string
	// var args []interface{}
	argIndex := 2

	// if username != "" {
	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
	// 	args = append(args, "%"+username+"%")
	// 	argIndex++
	// }

	// if len(conditions) > 0 {Fin
	// 	query += " WHERE " + strings.Join(conditions, " AND ")
	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
	// }

	// log.Print(countQuery)
	// log.Print(countArgs...)

	countArgs = append(countArgs, userId)

	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	queryArgs = append(queryArgs, userId, pageMetadata.Size, pageMetadata.Offset)
	// log.Println(query)
	// log.Println(queryArgs...)

	rows, err := tx.QueryxContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		result := new(entity.Explore)
		if err := rows.StructScan(result); err != nil {
			return nil, nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}

func (r *exploreRepository) UserAddWishlist(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
	query := "UPDATE user_similar_photos SET is_wishlist = true WHERE photo_id = $1 AND user_id = $2"
	_, err := tx.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *exploreRepository) UserDeleteWishlist(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
	query := "UPDATE user_similar_photos SET is_wishlist = false WHERE photo_id = $1 AND user_id = $2"
	_, err := tx.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *exploreRepository) FindAllUserFavorite(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {

	results := make([]*entity.Explore, 0)

	var totalItems int
	countQuery := `SELECT COUNT(*) FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_favorite = true`

	var countArgs []interface{}

	query := `SELECT usp.photo_id, usp.user_id, usp.similarity, usp.is_wishlist,
	usp.is_resend, usp.is_cart, usp.is_favorite, p.creator_id, p.title, p.is_this_you_url,
	p.your_moments_url, p.price, p.price_str, p. original_at, p.created_at, p.updated_at
	FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_favorite = true`
	//TO DO ADD logic for
	var queryArgs []interface{}

	// var conditions []string
	// var args []interface{}
	argIndex := 2

	// if username != "" {
	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
	// 	args = append(args, "%"+username+"%")
	// 	argIndex++
	// }

	// if len(conditions) > 0 {Fin
	// 	query += " WHERE " + strings.Join(conditions, " AND ")
	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
	// }

	// log.Print(countQuery)
	// log.Print(countArgs...)

	countArgs = append(countArgs, userId)

	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	queryArgs = append(queryArgs, userId, pageMetadata.Size, pageMetadata.Offset)
	// log.Println(query)
	// log.Println(queryArgs...)

	rows, err := tx.QueryxContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		result := new(entity.Explore)
		if err := rows.StructScan(result); err != nil {
			return nil, nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}

func (r *exploreRepository) UserAddFavorite(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
	query := "UPDATE user_similar_photos SET is_favorite = true WHERE photo_id = $1 AND user_id = $2"
	_, err := tx.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *exploreRepository) UserDeleteFavorite(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
	query := "UPDATE user_similar_photos SET is_favorite = false WHERE photo_id = $1 AND user_id = $2"
	_, err := tx.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *exploreRepository) FindAllUserCart(ctx context.Context, tx Querier, page, size int, similarity uint32, userId string) ([]*entity.Explore, *model.PageMetadata, error) {

	results := make([]*entity.Explore, 0)

	var totalItems int
	countQuery := `SELECT COUNT(*) FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_cart = true`

	var countArgs []interface{}

	query := `SELECT usp.photo_id, usp.user_id, usp.similarity, usp.is_wishlist,
	usp.is_resend, usp.is_cart, usp.is_favorite, p.creator_id, p.title, p.is_this_you_url,
	p.your_moments_url, p.price, p.price_str, p. original_at, p.created_at, p.updated_at
	FROM user_similar_photos AS usp JOIN photos AS p on p.id = usp.photo_id WHERE usp.user_id = $1 AND usp.is_cart = true`
	//TO DO ADD logic for
	var queryArgs []interface{}

	// var conditions []string
	// var args []interface{}
	argIndex := 2

	// if username != "" {
	// 	conditions = append(conditions, "username LIKE $"+strconv.Itoa(argIndex))
	// 	args = append(args, "%"+username+"%")
	// 	argIndex++
	// }

	// if len(conditions) > 0 {Fin
	// 	query += " WHERE " + strings.Join(conditions, " AND ")
	// 	countQuery += " WHERE " + strings.Join(conditions, " AND ")
	// }

	// log.Print(countQuery)
	// log.Print(countArgs...)

	countArgs = append(countArgs, userId)

	if err := tx.GetContext(ctx, &totalItems, countQuery, countArgs...); err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	queryArgs = append(queryArgs, userId, pageMetadata.Size, pageMetadata.Offset)
	// log.Println(query)
	// log.Println(queryArgs...)

	rows, err := tx.QueryxContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		result := new(entity.Explore)
		if err := rows.StructScan(result); err != nil {
			return nil, nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}

func (r *exploreRepository) UserAddCart(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
	query := "UPDATE user_similar_photos SET is_cart = true WHERE photo_id = $1 AND user_id = $2"
	_, err := tx.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *exploreRepository) UserDeleteCart(ctx context.Context, tx Querier, similarity uint32, photoId, userId string) error {
	query := "UPDATE user_similar_photos SET is_cart = false WHERE photo_id = $1 AND user_id = $2"
	_, err := tx.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return err
	}
	return nil
}
