-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_wallets (
    id CHAR(30) PRIMARY KEY,
    wallet_id CHAR(30) NOT NULL,
    transaction_detail_id CHAR(30) UNIQUE,
    amount INT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (transaction_detail_id) REFERENCES transaction_details(id)
) 
-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_wallets;
-- +goose StatementEnd