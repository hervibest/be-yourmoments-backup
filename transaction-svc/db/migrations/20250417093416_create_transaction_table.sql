-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_status') THEN
        CREATE TYPE transaction_status AS ENUM (
            'PENDING',
            'EXPIRED',
            'SUCCESS',
            'CANCELED',
            'REFUNDED',
            'FAILED'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_internal_status') THEN
        CREATE TYPE transaction_internal_status AS ENUM (
            'PENDING',
            'TOKEN_READY',
            'EXPIRED',
            'EXPIRED_CHECKED_INVALID',
            'EXPIRED_CHECKED_VALID',
            'LATE_SETTLEMENT',
            'SETTLED',
            'CANCELED_BY_SYSTEM',
            'CANCELED_BY_USER',
            'REFUNDED',
            'FAILED'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'midtrans_payment_status') THEN
        CREATE TYPE midtrans_payment_status AS ENUM (
            'capture',
            'settlement',
            'pending',
            'deny',
            'cancel',
            'expire',
            'failure'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS transactions (
    id uuid PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    status transaction_status NOT NULL,
    internal_status transaction_internal_status NOT NULL,
    transaction_method_id CHAR(26),
    transaction_type_id CHAR(26),
    payment_type_id CHAR(26),
    payment_at TIMESTAMPTZ,
    checkout_at TIMESTAMPTZ NOT NULL,
    snap_token TEXT,
    external_settlement_at TIMESTAMPTZ,
    external_status midtrans_payment_status,
    external_callback_response JSON,
    amount INT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX  IF NOT EXISTS idx_transaction_user_id ON transactions(user_id);
CREATE INDEX  IF NOT EXISTS idx_transaction_internal_status ON transactions(internal_status);
CREATE INDEX  IF NOT EXISTS idx_transaction_status ON transactions(status);
CREATE INDEX  IF NOT EXISTS idx_transaction_midtrans_payment_status ON transactions(external_status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;

DROP TYPE IF EXISTS transaction_internal_status;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS midtrans_payment_status;
DROP INDEX IF EXISTS idx_transaction_user_id;
-- +goose StatementEnd
