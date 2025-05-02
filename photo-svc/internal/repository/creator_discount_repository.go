package repository

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"

	"github.com/jmoiron/sqlx"
)

type discountPreparedStmt struct {
	findManyByCreatorIds *sqlx.Stmt
}

func newDiscountPreparedStmt(db *sqlx.DB) (*discountPreparedStmt, error) {
	findManyByCreatorIdsStmt, err := db.Preparex(`
	SELECT * FROM creator_discounts WHERE creator_id = ANY($1) AND is_active = true 
	ORDER BY creator_id, min_quantity DESC`)
	if err != nil {
		return nil, err
	}

	return &discountPreparedStmt{
		findManyByCreatorIds: findManyByCreatorIdsStmt,
	}, nil
}

type CreatorDiscountRepository interface {
	Create(ctx context.Context, tx Querier, discount *entity.CreatorDiscount) (*entity.CreatorDiscount, error)
	Activate(ctx context.Context, tx Querier, discountId string) error
	Deactivate(ctx context.Context, tx Querier, discountId string) error
	FindByIdAndCreatorId(ctx context.Context, tx Querier, discountId, creatorId string) (*entity.CreatorDiscount, error)
	GetDiscountRules(ctx context.Context, creatorIds []string) (*[]*entity.CreatorDiscount, error)
	FindAll(ctx context.Context, tx Querier, creatorId string) (*[]*entity.CreatorDiscount, error)
}

type creatorDiscountRepository struct {
	discountPreparedStmt *discountPreparedStmt
}

func NewCreatorDiscountRepository(db *sqlx.DB) (CreatorDiscountRepository, error) {
	discountPreparedStmt, err := newDiscountPreparedStmt(db)
	if err != nil {
		return nil, err
	}

	return &creatorDiscountRepository{
		discountPreparedStmt: discountPreparedStmt,
	}, nil
}

func (r *creatorDiscountRepository) Create(ctx context.Context, tx Querier, discount *entity.CreatorDiscount) (*entity.CreatorDiscount, error) {
	query := `
	INSERT INTO creator_discounts 
	(id, creator_id, name, min_quantity, discount_type, value, is_active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	_, err := tx.ExecContext(ctx, query, discount.Id, discount.CreatorId, discount.Name,
		discount.MinQuantity, discount.DiscountType, discount.Value, discount.IsActive,
		discount.CreatedAt, discount.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return discount, nil
}

func (r *creatorDiscountRepository) Activate(ctx context.Context, tx Querier, discountId string) error {
	query := `UPDATE creator_discounts SET is_active = true WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, discountId)
	if err != nil {
		return err
	}
	return nil
}

func (r *creatorDiscountRepository) Deactivate(ctx context.Context, tx Querier, discountId string) error {
	query := `UPDATE creator_discounts SET is_active = false WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, discountId)
	if err != nil {
		return err
	}
	return nil
}

func (r *creatorDiscountRepository) FindByIdAndCreatorId(ctx context.Context, tx Querier, discountId, creatorId string) (*entity.CreatorDiscount, error) {
	discount := new(entity.CreatorDiscount)
	query := `SELECT * from creator_discounts WHERE id = $1 AND creator_id =  $2`
	if err := tx.GetContext(ctx, discount, query, discountId, creatorId); err != nil {
		return nil, err
	}
	return discount, nil
}

func (r *creatorDiscountRepository) FindAll(ctx context.Context, tx Querier, creatorId string) (*[]*entity.CreatorDiscount, error) {
	discounts := make([]*entity.CreatorDiscount, 0)
	query := `SELECT * from creator_discounts WHERE creator_id =  $1`
	if err := tx.SelectContext(ctx, &discounts, query, creatorId); err != nil {
		return nil, err
	}
	return &discounts, nil
}

// func (r *creatorDiscountRepository) FindAll(ctx context.Context, tx Querier, discountId string) (*entity.CreatorDiscount, error) {
// 	discount := new(entity.CreatorDiscount)
// 	query := `SELECT * from creator_discounts WHERE id = $1`
// 	if err := tx.GetContext(ctx, discount, query, discountId); err != nil {
// 		log.Println("Error happen in FindById  discount with error ", err)
// 		return nil, err
// 	}
// 	return discount, nil
// }

func (r *creatorDiscountRepository) GetDiscountRules(ctx context.Context, creatorIds []string) (*[]*entity.CreatorDiscount, error) {
	photos := make([]*entity.CreatorDiscount, 0)
	err := r.discountPreparedStmt.findManyByCreatorIds.SelectContext(ctx, &photos, creatorIds)
	if err != nil {
		return nil, err
	}
	return &photos, nil
}
