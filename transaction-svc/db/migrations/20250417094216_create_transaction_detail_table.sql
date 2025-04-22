-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_details (
    id CHAR(30) PRIMARY KEY,
    transaction_id uuid NOT NULL,
    creator_id CHAR(30) NOT NULL,
    subtotal_price INT NOT NULL,
    creator_discount_id CHAR(30) NOT NULL,
    is_reviewed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
) 
-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_details;
-- +goose StatementEnd