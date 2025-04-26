-- +goose Up
-- +goose StatementBegin
ALTER TABLE photos
ADD COLUMN bulk_photo_id CHAR(30),
FOREIGN KEY (bulk_photo_id) REFERENCES bulk_photos(bulk_photo_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos
DROP COLUMN IF EXISTS bulk_photo_id;
-- +goose StatementEnd
