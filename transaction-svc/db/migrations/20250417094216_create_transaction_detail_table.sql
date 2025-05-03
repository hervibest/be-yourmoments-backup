-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_details (
    id CHAR(26) PRIMARY KEY,
    transaction_id uuid NOT NULL,
    creator_id CHAR(26) NOT NULL,
    subtotal_price INT NOT NULL,
    creator_discount_id CHAR(26) NOT NULL,
    is_reviewed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);
CREATE INDEX IF NOT EXISTS idx_transaction_detail_creator_id ON transaction_details(creator_id);
CREATE INDEX IF NOT EXISTS idx_transaction_detail_creator_discount_id ON transaction_details(creator_discount_id);
-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_details;
DROP INDEX IF EXISTS idx_transaction_detail_creator_id;
DROP INDEX IF EXISTS idx_transaction_detail_creator_discount_id;
-- +goose StatementEnd