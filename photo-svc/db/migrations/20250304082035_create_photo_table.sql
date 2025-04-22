-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS photos (
    id CHAR(26) PRIMARY KEY NOT NULL,
    creator_id CHAR(26) NOT NULL,
    title VARCHAR(100) NOT NULL,
    owned_by_user_id CHAR(26),
    compressed_url TEXT,
    is_this_you_url TEXT,
    your_moments_url TEXT,
    collection_url TEXT,
    price INT NOT NULL,
    price_str VARCHAR(100) NOT NULL,
    original_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY(creator_id) REFERENCES creators(id)
);

CREATE INDEX IF NOT EXISTS idx_photos_owned_by_user_id ON photos (owned_by_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photos_owned_by_user_id;
DROP TABLE IF EXISTS photos;
-- +goose StatementEnd
