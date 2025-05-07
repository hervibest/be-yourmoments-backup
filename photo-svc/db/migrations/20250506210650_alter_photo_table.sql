-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'photo_status') THEN
        CREATE TYPE photo_status AS ENUM (
            'AVAILABLE',
            'IN_TRANSACTION',
            'SOLD'
        );
    END IF;
END$$;

ALTER TABLE photos ADD COLUMN status photo_status DEFAULT 'AVAILABLE' NOT NULL;

CREATE INDEX IF NOT EXISTS idx_photo_status_creator ON photos (status, creator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photo_status_creator;
ALTER TABLE photos DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
