-- +goose Up
-- +goose StatementBegin
ALTER TABLE photos ADD COLUMN bulk_photo_id CHAR(26);
ALTER TABLE photos ADD CONSTRAINT fk_bulk_photo FOREIGN KEY (bulk_photo_id) REFERENCES bulk_photos(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos DROP CONSTRAINT IF EXISTS fk_bulk_photo;
ALTER TABLE photos DROP COLUMN IF EXISTS bulk_photo_id;
-- +goose StatementEnd
