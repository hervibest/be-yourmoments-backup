-- +goose Up
-- +goose StatementBegin
CREATE TYPE discount_type AS ENUM ('FLAT', 'PERCENT');

CREATE TABLE IF NOT EXISTS creator_discounts (
    id CHAR(26) PRIMARY KEY NOT NULL,
    creator_id CHAR(26) NOT NULL,
    name VARCHAR(100) NOT NULL,
    min_quantity INT NOT NULL,
    discount_type discount_type NOT NULL,
    value INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (creator_id) REFERENCES creators(id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS creator_discounts;
DROP TYPE discount_type;
-- +goose StatementEnd