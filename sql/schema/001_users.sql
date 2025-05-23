-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    hashed_password TEXT NOT NULL DEFAULT 'unset',
    is_chirpy_red BOOLEAN DEFAULT false,
    UNIQUE(email)
);

-- +goose Down
DROP TABLE users;
