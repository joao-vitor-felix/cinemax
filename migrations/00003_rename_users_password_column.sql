-- +goose Up
ALTER TABLE users
RENAME COLUMN hashed_password TO password_hash;

-- +goose Down
ALTER TABLE users
RENAME COLUMN password_hash TO hashed_password;

