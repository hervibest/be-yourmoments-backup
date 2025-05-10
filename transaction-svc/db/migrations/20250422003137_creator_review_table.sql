-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS creator_reviews (
    id CHAR(26) PRIMARY KEY,
    transaction_detail_id CHAR(26) UNIQUE NOT NULL,
    creator_id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (transaction_detail_id) REFERENCES transaction_details(id)
);
CREATE INDEX  IF NOT EXISTS idx_creator_review_creator_id ON creator_reviews(creator_id);
CREATE INDEX  IF NOT EXISTS idx_creator_review_user_id ON creator_reviews(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS creator_reviews;
DROP INDEX IF EXISTS idx_creator_review_creator_id;
DROP INDEX IF EXISTS idx_creator_review_user_id;
-- +goose StatementEnd
