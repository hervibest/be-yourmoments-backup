-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS user_similar_photos (
    photo_id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    PRIMARY KEY (photo_id, user_id),
    similarity SMALLINT CHECK (similarity BETWEEN 1 AND 9) NOT NULL DEFAULT 5,
    is_wishlist boolean DEFAULT false,
    is_resend boolean DEFAULT false,
    is_cart boolean DEFAULT false,
    is_favorite boolean DEFAULT false,
    created_at TIMESTAMPTZ not null,
    updated_at TIMESTAMPTZ not null,
    FOREIGN KEY(photo_id) REFERENCES photos(id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_similar_photos;


-- +goose StatementEnd