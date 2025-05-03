-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallets (
    id CHAR(26) PRIMARY KEY,
    creator_id CHAR(26) NOT NULL UNIQUE, 
    balance INT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
) ;
CREATE INDEX IF NOT EXISTS idx_wallet_creator_id ON wallets(creator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
DROP INDEX IF EXISTS idx_ wallet_creator_id;
-- +goose StatementEnd

