-- +goose Up
-- +goose StatementBegin
CREATE TYPE bulk_photo_status AS ENUM ('PROCESSED', 'FAILED', 'CANCELED', 'SUCCESS');

CREATE TABLE IF NOT EXISTS bulk_photos (
    id CHAR(26) PRIMARY KEY NOT NULL,
    creator_id CHAR(26) NOT NULL,
    bulk_photo_status bulk_photo_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_creator_id ON bulk_photos (creator_id);
CREATE INDEX IF NOT EXISTS idx_bulk_photo_status ON bulk_photos (bulk_photo_status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_creator_id;
DROP INDEX IF EXISTS idx_bulk_photo_status;
DROP TABLE IF EXISTS bulk_photos;
-- +goose StatementEnd
