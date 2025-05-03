-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banks (
    id CHAR(26) PRIMARY KEY,
    bank_code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    alias VARCHAR(50),
    swift_code VARCHAR(20),
    logo_url TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banks;
-- +goose StatementEnd
