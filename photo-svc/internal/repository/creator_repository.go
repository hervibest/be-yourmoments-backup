package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type creatorPreparedStmt struct {
	findByUserId *sqlx.Stmt
}

func newCreatorPreparedStmt(db *sqlx.DB) (*creatorPreparedStmt, error) {
	findByUserIdStmt, err := db.Preparex("SELECT * FROM creators WHERE user_id = $1")
	if err != nil {
		return nil, err
	}

	return &creatorPreparedStmt{
		findByUserId: findByUserIdStmt,
	}, nil
}

type CreatorRepository interface {
	Create(ctx context.Context, tx Querier, creator *entity.Creator) (*entity.Creator, error)
	FindByUserId(ctx context.Context, userId string) (*entity.Creator, error)
}

type creatorRepository struct {
	creatorPreparedStmt *creatorPreparedStmt
}

func NewCreatorRepository(db *sqlx.DB) (CreatorRepository, error) {
	creatorPreparedStmt, err := newCreatorPreparedStmt(db)
	if err != nil {
		return nil, err
	}

	return &creatorRepository{
		creatorPreparedStmt: creatorPreparedStmt,
	}, nil
}

func (r *creatorRepository) Create(ctx context.Context, tx Querier, creator *entity.Creator) (*entity.Creator, error) {
	query := `INSERT INTO creators  (id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := tx.ExecContext(ctx, query, creator.Id, creator.UserId, creator.CreatedAt, creator.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert creator: %w", err)
	}

	return creator, nil
}

func (r *creatorRepository) FindByUserId(ctx context.Context, userId string) (*entity.Creator, error) {
	creator := new(entity.Creator)

	row := r.creatorPreparedStmt.findByUserId.QueryRowxContext(ctx, userId)
	if err := row.StructScan(creator); err != nil {
		return nil, err
	}

	return creator, nil
}

func (r *creatorRepository) UpdateCreatorRating(ctx context.Context, tx Querier, creator *entity.Creator) (*entity.Creator, error) {
	query := `INSERT INTO creators  (id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := tx.ExecContext(ctx, query, creator.Id, creator.UserId, creator.CreatedAt, creator.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert creator: %w", err)
	}

	return creator, nil
}
