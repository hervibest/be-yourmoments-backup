-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
ADD COLUMN photo_ids JSON
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
DROP COLUMN IF EXISTS photo_ids
-- +goose StatementEnd
