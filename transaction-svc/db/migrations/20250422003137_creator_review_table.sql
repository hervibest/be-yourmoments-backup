-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS creator_reviews (
    id CHAR(30) PRIMARY KEY, -- ULID
    transaction_detail_id CHAR(30) UNIQUE NOT NULL,
    creator_id CHAR(30) NOT NULL,
    user_id CHAR(30) NOT NULL,
    star INT NOT NULL CHECK (star BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    
    FOREIGN KEY (transaction_detail_id) REFERENCES transaction_details(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS creator_reviews;
-- +goose StatementEnd
