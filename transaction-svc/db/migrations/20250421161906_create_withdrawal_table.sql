-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'withdrawal_status') THEN
        CREATE TYPE withdrawal_status AS ENUM (
            'PENDING',
            'SUCCESS',
            'FAILED',
            'CANCELED',
            'EXPIRED'
        );
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS withdrawals (
    id CHAR(26) PRIMARY KEY,
    wallet_id CHAR(26) NOT NULL,
    bank_wallet_id CHAR(26) NOT NULL,
    amount INT NOT NULL,
    status withdrawal_status NOT NULL ,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (bank_wallet_id) REFERENCES bank_wallets(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;

DROP TYPE IF EXISTS withdrawal_status;
-- +goose StatementEnd