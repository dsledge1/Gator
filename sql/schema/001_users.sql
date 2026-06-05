-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    name TEXT UNIQUE NOT NULL
);
ps
-- +goose Down
DROP TABLE users;