-- +goose Up
-- +goose StatementBegin
ALTER TABLE photos ADD COLUMN total_user_similar INT DEFAULT 0 NOT NULL;
CREATE INDEX IF NOT EXISTS idx_photo_creator_id ON photos (creator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photo_creator_id;
ALTER TABLE photos DROP COLUMN IF EXISTS total_user_similar;
-- +goose StatementEnd

