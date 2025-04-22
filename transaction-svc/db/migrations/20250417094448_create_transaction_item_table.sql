-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_items (
    id CHAR(30) PRIMARY KEY,
    transaction_detail_id CHAR(30) NOT NULL,
    photo_id CHAR(30) NOT NULL,
    price INT NOT NULL,
    discount INT,
    final_price INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (transaction_detail_id) REFERENCES transaction_details(id)
) 

-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_items;
-- +goose StatementEnd