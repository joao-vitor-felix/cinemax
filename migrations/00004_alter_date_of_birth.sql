-- +goose Up
ALTER TABLE users
ALTER COLUMN date_of_birth DROP NOT NULL;

-- +goose Down
ALTER TABLE users
ALTER COLUMN date_of_birth SET NOT NULL;

