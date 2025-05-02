-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    id CHAR(26) PRIMARY KEY NOT NULL,
    user_id CHAR(26) UNIQUE NOT NULL ,
    birth_date DATE,
    nickname VARCHAR(100) UNIQUE,
    biography TEXT,
    profile_url TEXT,
    profile_cover_url TEXT,
    similarity SMALLINT CHECK (similarity BETWEEN 1 AND 9) NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd