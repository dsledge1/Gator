-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;