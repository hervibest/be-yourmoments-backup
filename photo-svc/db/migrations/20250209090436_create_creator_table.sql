-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS creators (
    id CHAR(26) PRIMARY KEY NOT NULL,
    user_id CHAR(26) NOT NULL UNIQUE,
    rating DECIMAL(3, 1) DEFAULT 0.0,
    rating_count INT DEFAULT 0,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_creator_user_id ON creators (user_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS creators;
DROP INDEX IF EXISTS idx_creator_user_id;
-- +goose StatementEnd

