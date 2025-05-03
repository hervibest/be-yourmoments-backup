-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bank_wallets (
    id CHAR(26) PRIMARY KEY,
    wallet_id CHAR(26) NOT NULL,
    bank_id CHAR(26) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (bank_id, wallet_id),
    FOREIGN KEY (wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (bank_id) REFERENCES banks(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bank_wallets;
-- +goose StatementEnd