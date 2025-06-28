-- +goose Up
-- +goose StatementBegin
CREATE TYPE platform AS ENUM (
    'ANDROID', 'IOS', 'WEB'
);

CREATE TABLE IF NOT EXISTS user_devices (
    id CHAR(26) PRIMARY KEY NOT NULL,
    user_id CHAR(26) NOT NULL,
    token TEXT NOT NULL,
    platform platform NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_user_devices_token_hash ON user_devices(token);
CREATE INDEX idx_user_devices_platform ON user_devices (platform);
CREATE UNIQUE INDEX idx_user_devices_userid_token ON user_devices (user_id, token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_devices_userid_token;
DROP INDEX IF EXISTS idx_user_devices_platform;
DROP INDEX IF EXISTS idx_user_devices_token_hash;
DROP TABLE IF EXISTS user_devices;
DROP TYPE IF EXISTS platform;
-- +goose StatementEnd
